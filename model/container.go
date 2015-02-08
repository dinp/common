package model

import (
	"fmt"
)

type Port struct {
	PublicPort int
}

func (this *Port) String() string {
	return fmt.Sprintf("<PublicPort:%d>", this.PublicPort)
}

type ContainerDto struct {
	Id      string
	Image   string
	AppName string
	Ports   []*Port
	Status  string
}

func (this *ContainerDto) String() string {
	return fmt.Sprintf(
		"<Id:%s, Image:%s, AppName:%s, Status:%s, Ports:%v>",
		this.Id,
		this.Image,
		this.AppName,
		this.Status,
		this.Ports,
	)
}

type Container struct {
	Id       string
	Ip       string
	Image    string
	AppName  string
	Ports    []*Port
	Status   string
	UpdateAt int64
}

func (this *Container) String() string {
	return fmt.Sprintf(
		"<Id:%s, Ip:%s, Image:%s, AppName:%s, Status:%s, Ports:%v>",
		this.Id,
		this.Ip,
		this.Image,
		this.AppName,
		this.Status,
		this.Ports,
	)
}
