package model

import (
	"fmt"
)

type Node struct {
	Ip       string
	MemFree  uint64
	UpdateAt int64
}

func (this *Node) String() string {
	return fmt.Sprintf(
		"<Ip:%s, MemFree:%d, UpdateAt:%d>",
		this.Ip,
		this.MemFree,
		this.UpdateAt,
	)
}

type NodeRequest struct {
	Node
	Containers []*ContainerDto
}

func (this NodeRequest) String() string {
	return fmt.Sprintf(
		"<Node:%v, Containers:%v>",
		this.Node,
		this.Containers,
	)
}

type NodeResponse struct {
	Code int
}

type NodeSlice []*Node

// Len is part of sort.Interface.
func (s NodeSlice) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s NodeSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface.
func (s NodeSlice) Less(i, j int) bool {
	return s[i].MemFree > s[j].MemFree
}
