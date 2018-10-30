package model

type Queue struct {
	Go []*ContainerInfo
	C []*ContainerInfo
}

type ContainerInfo struct {
	Id        string
	IsRunning bool
}
