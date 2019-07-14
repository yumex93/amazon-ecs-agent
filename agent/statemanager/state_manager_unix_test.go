// +build !windows,unit

// Copyright 2014-2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.
package statemanager_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	apitask "github.com/aws/amazon-ecs-agent/agent/api/task"
	"github.com/aws/amazon-ecs-agent/agent/config"
	"github.com/aws/amazon-ecs-agent/agent/engine"
	"github.com/aws/amazon-ecs-agent/agent/engine/dockerstate"
	engine_testutils "github.com/aws/amazon-ecs-agent/agent/engine/testutils"
	"github.com/aws/amazon-ecs-agent/agent/statemanager"
	"github.com/aws/amazon-ecs-agent/agent/taskresource/efs"
	resourcestatus "github.com/aws/amazon-ecs-agent/agent/taskresource/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateManager(t *testing.T) {
	tmpDir, err := ioutil.TempDir("/tmp", "ecs_statemanager_test")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpDir)
	cfg := &config.Config{DataDir: tmpDir}
	manager, err := statemanager.NewStateManager(cfg)
	assert.Nil(t, err, "Error loading manager")

	err = manager.Load()
	assert.Nil(t, err, "Expected loading a non-existent file to not be an error")

	// Now let's make some state to save
	containerInstanceArn := ""
	taskEngine := engine.NewTaskEngine(&config.Config{}, nil, nil, nil, nil, dockerstate.NewTaskEngineState(),
		nil, nil)

	manager, err = statemanager.NewStateManager(cfg, statemanager.AddSaveable("TaskEngine", taskEngine),
		statemanager.AddSaveable("ContainerInstanceArn", &containerInstanceArn))
	require.Nil(t, err)

	containerInstanceArn = "containerInstanceArn"

	testTask := &apitask.Task{Arn: "test-arn"}
	taskEngine.(*engine.DockerTaskEngine).State().AddTask(testTask)

	err = manager.Save()
	require.Nil(t, err, "Error saving state")

	assertFileMode(t, filepath.Join(tmpDir, "ecs_agent_data.json"))

	// Now make sure we can load that state sanely
	loadedTaskEngine := engine.NewTaskEngine(&config.Config{}, nil, nil, nil, nil, dockerstate.NewTaskEngineState(),
		nil, nil)
	var loadedContainerInstanceArn string

	manager, err = statemanager.NewStateManager(cfg, statemanager.AddSaveable("TaskEngine", &loadedTaskEngine),
		statemanager.AddSaveable("ContainerInstanceArn", &loadedContainerInstanceArn))
	require.Nil(t, err)

	err = manager.Load()
	require.Nil(t, err, "Error loading state")

	assert.Equal(t, containerInstanceArn, loadedContainerInstanceArn, "Did not load containerInstanceArn correctly")

	if !engine_testutils.DockerTaskEnginesEqual(loadedTaskEngine.(*engine.DockerTaskEngine), (taskEngine.(*engine.DockerTaskEngine))) {
		t.Error("Did not load taskEngine correctly")
	}

	// I'd rather double check .Equal there; let's make sure ListTasks agrees.
	tasks, err := loadedTaskEngine.ListTasks()
	assert.Nil(t, err, "Error listing tasks")
	require.Equal(t, 1, len(tasks), "Should have a task!")
	assert.Equal(t, "test-arn", tasks[0].Arn, "Wrong arn")
}

func assertFileMode(t *testing.T, path string) {
	info, err := os.Stat(path)
	assert.Nil(t, err)

	mode := info.Mode()
	assert.Equal(t, os.FileMode(0600), mode, "Wrong file mode")
}

// verify that the state manager correctly loads the existing task networking related fields in state file.
// if we change those fields in the future, we should modify this test to test the new fields
func TestLoadsDataForAWSVPCTask(t *testing.T) {
	cfg := &config.Config{DataDir: filepath.Join(".", "testdata", "v11", "task-networking")}

	taskEngine := engine.NewTaskEngine(&config.Config{}, nil, nil, nil, nil, dockerstate.NewTaskEngineState(), nil, nil)
	var containerInstanceArn, cluster, savedInstanceID string

	stateManager, err := statemanager.NewStateManager(cfg,
		statemanager.AddSaveable("TaskEngine", taskEngine),
		statemanager.AddSaveable("ContainerInstanceArn", &containerInstanceArn),
		statemanager.AddSaveable("Cluster", &cluster),
		statemanager.AddSaveable("EC2InstanceID", &savedInstanceID),
	)
	assert.NoError(t, err)
	err = stateManager.Load()
	assert.NoError(t, err)

	assert.Equal(t, "state-file", cluster)

	tasks, err := taskEngine.ListTasks()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))

	task := tasks[0]
	assert.Equal(t, "arn:aws:ecs:us-west-2:1234567890:task/fad405be-8705-4175-877b-db50109a15f2", task.Arn)
	assert.Equal(t, "task-networking-state", task.Family)
	assert.NotNil(t, task.ENI)

	eni := task.ENI
	assert.Equal(t, "eni-089ba8329b8e3f6ec", eni.ID)
	assert.Equal(t, "ip-172-31-10-246.us-west-2.compute.internal", eni.GetHostname())

	ipv4Addresses := eni.GetIPV4Addresses()
	assert.Equal(t, 1, len(ipv4Addresses))
	assert.Equal(t, "172.31.10.246", ipv4Addresses[0])
}

// verify that the state manager correctly loads gpu related fields in state file
func TestLoadsDataForGPU(t *testing.T) {
	cfg := &config.Config{DataDir: filepath.Join(".", "testdata", "v18", "gpu")}
	taskEngine := engine.NewTaskEngine(&config.Config{}, nil, nil, nil, nil, dockerstate.NewTaskEngineState(), nil, nil)
	var containerInstanceArn, cluster, savedInstanceID string
	var sequenceNumber int64
	stateManager, err := statemanager.NewStateManager(cfg,
		statemanager.AddSaveable("TaskEngine", taskEngine),
		statemanager.AddSaveable("ContainerInstanceArn", &containerInstanceArn),
		statemanager.AddSaveable("Cluster", &cluster),
		statemanager.AddSaveable("EC2InstanceID", &savedInstanceID),
		statemanager.AddSaveable("SeqNum", &sequenceNumber),
	)
	assert.NoError(t, err)

	err = stateManager.Load()
	assert.NoError(t, err)
	assert.Equal(t, "state-file", cluster)
	assert.EqualValues(t, 0, sequenceNumber)

	tasks, err := taskEngine.ListTasks()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))

	task := tasks[0]
	assert.Equal(t, "arn:aws:ecs:us-west-2:1234567890:task/33425c99-5db7-45fb-8244-bc94d00661e4", task.Arn)
	assert.Equal(t, "gpu-state", task.Family)
	assert.Equal(t, 1, len(task.Containers))
	assert.Equal(t, 2, len(task.Associations))

	association1 := task.Associations[0]
	assert.Equal(t, "container_1", association1.Containers[0])
	assert.Equal(t, "gpu", association1.Type)
	assert.Equal(t, "0", association1.Name)

	association2 := task.Associations[1]
	assert.Equal(t, "container_1", association2.Containers[0])
	assert.Equal(t, "gpu", association2.Type)
	assert.Equal(t, "1", association2.Name)

	container := task.Containers[0]
	assert.Equal(t, "container_1", container.Name)
	assert.Equal(t, []string{"0", "1"}, container.GPUIDs)
	assert.Equal(t, "0,1", container.Environment["NVIDIA_VISIBLE_DEVICES"])
}

// verify that the state manager correctly loads efs related fields in state file
func TestLoadsDataForEFS(t *testing.T) {
	cfg := &config.Config{DataDir: filepath.Join(".", "testdata", "v23", "efs")}
	taskEngine := engine.NewTaskEngine(&config.Config{}, nil, nil, nil, nil, dockerstate.NewTaskEngineState(), nil, nil)
	var containerInstanceArn, cluster, savedInstanceID string
	var sequenceNumber int64
	stateManager, err := statemanager.NewStateManager(cfg,
		statemanager.AddSaveable("TaskEngine", taskEngine),
		statemanager.AddSaveable("ContainerInstanceArn", &containerInstanceArn),
		statemanager.AddSaveable("Cluster", &cluster),
		statemanager.AddSaveable("EC2InstanceID", &savedInstanceID),
		statemanager.AddSaveable("SeqNum", &sequenceNumber),
	)
	assert.NoError(t, err)

	err = stateManager.Load()
	assert.NoError(t, err)
	assert.Equal(t, "state-file", cluster)
	assert.EqualValues(t, 0, sequenceNumber)

	tasks, err := taskEngine.ListTasks()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))

	task := tasks[0]
	assert.Equal(t, "arn:aws:ecs:us-west-2:1234567890:task/fad405be-8705-4175-877b-db50109a15f2", task.Arn)
	assert.Equal(t, "efs-state", task.Family)
	assert.Equal(t, 2, len(task.Containers))

	volume := task.Volumes[0]
	assert.Equal(t, "efsVolume", volume.Name)
	assert.Equal(t, "efs", volume.Type)
	vol := volume.Volume.(*efs.EFSConfig)
	assert.Equal(t, "fs-324562", vol.FileSystem)
	assert.Equal(t, "/", vol.RootDir)
	assert.Equal(t, "/data/efs/ecs-efs-state-1-efsVolume-a68ef4b6e0fba38d3500", vol.TargetDir)
	assert.Equal(t, "fs-324562.us-west-2.amazon.com", vol.DNSEndpoints[0])
	assert.Equal(t, "rsize=1048576,wsize=1048576,timeo=10,hard,retrans=2,noresvport,vers=4", vol.MountOptions)
	assert.True(t, vol.ReadOnly)
	assert.Equal(t, "/var/lib/data", vol.DataDirOnHost)

	efsRes, ok := task.ResourcesMapUnsafe["efs"]
	assert.True(t, ok)
	createDep := efsRes[0].GetContainerDependencies(resourcestatus.ResourceCreated)
	assert.Equal(t, "~internal~ecs~pause", createDep[0].ContainerName)
	assert.Equal(t, "RESOURCES_PROVISIONED", createDep[0].SatisfiedStatus.String())
	removeDep := efsRes[0].GetContainerDependencies(resourcestatus.ResourceRemoved)
	assert.Equal(t, "container_1", removeDep[0].ContainerName)
	assert.Equal(t, "STOPPED", removeDep[0].SatisfiedStatus.String())

	assert.Equal(t, resourcestatus.ResourceCreated, efsRes[0].GetDesiredStatus())
	assert.Equal(t, resourcestatus.ResourceCreated, efsRes[0].GetKnownStatus())
	assert.Equal(t, resourcestatus.ResourceStatusNone, efsRes[0].GetAppliedStatus())
}
