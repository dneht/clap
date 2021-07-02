package refer

const (
	DeployAppPackageType = "package"
	DeployAppDeployType  = "deploy"
	RenderAppJsonnetType = "jsonnet"

	//1打包中、2打包完成、3打包失败、6已有发布
	DeployStatusBuilding  = 1
	DeployStatusBuildEnd  = 2
	DeployStatusBuildFail = 3
	DeployStatusPackHear  = 6
)

var ConnectErrorDeployStatus = DeployStatus{
	Status: "ConnectError",
}

var FailedDeployStatus = DeployStatus{
	Status: "Failed",
}

var DefaultDeployStatus = DeployStatus{
	Status: "Creating",
}
