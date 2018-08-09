# Copyright 2014-2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

# Build

#Param(
#    [string]$S_REPOSITORY,
#    [string]$S_REGION,
#    [string]$R_REPOSITORY,
#    [string]$R_REGION
#)

$S_REPOSITORY="115215150660.dkr.ecr.us-west-2.amazonaws.com"
$S_REGION="us-west-2"
$R_REPOSITORY="632651270272.dkr.ecr.cn-northwest-1.amazonaws.com.cn"
$R_REGION="cn-northwest-1"


Invoke-Expression "${PSScriptRoot}\..\misc\windows-iam\Setup_Iam.ps1"
Invoke-Expression "${PSScriptRoot}\..\misc\windows-listen80\Setup_Listen80.ps1"
Invoke-Expression "${PSScriptRoot}\..\misc\windows-telemetry\build.ps1"
Invoke-Expression "${PSScriptRoot}\..\misc\windows-python\build.ps1"
Invoke-Expression "${PSScriptRoot}\..\misc\container-health-windows\build.ps1"

# Login ECR
$sAccessKey=(Get-SSMParameter -Name SAccessKey -WithDecryption $TRUE).Value
$sSecretKey=(Get-SSMParameter -Name SSecretKey -WithDecryption $TRUE).Value
Invoke-Expression -Command (Get-ECRLoginCommand -Region "${S_REGION}").Command

$array = @("amazon-ecs-windows-telemetry-test", "amazon-ecs-windows-python", "amazon-ecs-containerhealthcheck")
foreach ($element in $array){
    $result=(GET-ECRRepository -Region "${S_REGION}"|Select-Object RepositoryName|Select-String -Pattern $element)
    if ($result -eq $null) { New-ECRRepository -Region "${S_REGION}" -RepositoryName $element }
    docker tag "amazon/${element}:make" "${S_REPOSITORY}/${element}:make"
    docker push "${S_REPOSITORY}/${element}:make"
    docker rmi "${S_REPOSITORY}/${element}:make"
}

$array = @("amazon-ecs-iamrolecontainer", "amazon-ecs-listen80")
foreach ($element in $array){
    $result=(GET-ECRRepository -Region "${S_REGION}"|Select-Object RepositoryName|Select-String -Pattern $element)
    if ($result -eq $null) { New-ECRRepository -Region "${S_REGION}" -RepositoryName $element }
    docker tag "amazon/${element}:latest" "${S_REPOSITORY}/${element}:latest"
    docker push "${S_REPOSITORY}/${element}:latest"
    docker rmi "${S_REPOSITORY}/${element}:latest"
}

#upload to replication ecr
#$accessKey=(Get-SSMParameter -Name AccessKey -WithDecryption $TRUE).Value
#$secretKey=(Get-SSMParameter -Name SecretKey -WithDecryption $TRUE).Value
#Set-AWSCredential -AccessKey $accessKey -SecretKey $secretKey -StoreAs default
#Invoke-Expression –Command (Get-ECRLoginCommand –Region "${R_REGION}").Command

#$array = @("amazon-ecs-windows-telemetry-test", "amazon-ecs-windows-python", "amazon-ecs-containerhealthcheck")
#foreach ($element in $array){
#    $result=(GET-ECRRepository -Region "${R_REGION}"|Select-Object RepositoryName|Select-String -Pattern $element)
#    if ($result -eq $null) { New-ECRRepository -Region "${R_REGION}" -RepositoryName $element }
#    docker tag "amazon/${element}:make" "${R_REPOSITORY}/${element}:make"
#    docker push "${R_REPOSITORY}/${element}:make"
#    docker rmi "${R_REPOSITORY}/${element}:make"
#}

#$array = @("amazon-ecs-iamrolecontainer", "amazon-ecs-listen80")
#foreach ($element in $array){
#    $result=(GET-ECRRepository -Region "${R_REGION}"|Select-Object RepositoryName|Select-String -Pattern $element)
#    if ($result -eq $null) { New-ECRRepository -Region "${R_REGION}" -RepositoryName $element }
#    docker tag "amazon/${element}:latest" "${R_REPOSITORY}/${element}:latest"
#    docker push "${R_REPOSITORY}/${element}:latest"
#    docker rmi "${R_REPOSITORY}/${element}:latest"
#}

Set-AWSCredential -AccessKey $sAccessKey -SecretKey $sSecretKey -StoreAs default