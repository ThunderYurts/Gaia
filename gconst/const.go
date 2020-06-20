package gconst

const (
	// YurtImage is the base url to yurt image
	YurtImage = "registry.cn-hangzhou.aliyuncs.com/se-devgo/yurt:v2.0"

	// ActionPort is port action server use in container
	ActionPort = "8080/tcp"
	// SyncPort is port action server use in container
	SyncPort = "8000/tcp"
	// YurtFilter is for container list
	YurtFilter = "label=app=thunderyurt"

	GaiaRoot = "/gaia"
)
