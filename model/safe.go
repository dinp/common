package model

import (
	"fmt"
	"sync"
)

type SafeApp struct {
	sync.RWMutex
	// key: Ip-ContainerId
	M map[string]*Container
	// if ToUpdate: sync routes to redis
	ToUpdate bool
}

type SafeRealState struct {
	sync.RWMutex
	// key: AppName
	M map[string]*SafeApp
}

func NewSafeApp() *SafeApp {
	// why set ToUpdate = true?
	// it can clear dirty routes in redis
	return &SafeApp{M: make(map[string]*Container), ToUpdate: true}
}

func NewSafeRealState() *SafeRealState {
	return &SafeRealState{M: make(map[string]*SafeApp)}
}

func MakeContainerKey(ip, containerId string) string {
	return fmt.Sprintf("%s-%s", ip, containerId)
}

func (this *SafeApp) AddContainer(c *Container) {
	if c.Ports == nil || len(c.Ports) == 0 {
		return
	}

	key := MakeContainerKey(c.Ip, c.Id)

	this.Lock()
	defer this.Unlock()

	old, exists := this.M[key]
	if !(exists && old.Ports[0].PublicPort == c.Ports[0].PublicPort) {
		this.ToUpdate = true
	}
	this.M[key] = c
}

func (this *SafeApp) ContainerCount() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.M)
}

func (this *SafeApp) IsOldVersion(newImage string) (isOld bool, olds []*Container) {
	this.RLock()
	defer this.RUnlock()
	if len(this.M) == 0 {
		isOld = true
		return
	}

	for _, c := range this.M {
		if c.Image != newImage {
			isOld = true
			break
		}
	}

	if isOld {
		for _, c := range this.M {
			olds = append(olds, c)
		}
	}

	return
}

func (this *SafeApp) Containers() (cs []*Container) {
	this.RLock()
	defer this.RUnlock()
	for _, c := range this.M {
		cs = append(cs, c)
	}
	return
}

func (this *SafeApp) IsNeedUpdateRouter() bool {
	this.RLock()
	defer this.RUnlock()
	return this.ToUpdate
}

func (this *SafeApp) NeedUpdateRouter(needUpdate bool) {
	this.Lock()
	defer this.Unlock()
	this.ToUpdate = needUpdate
}

func (this *SafeRealState) UpdateContainer(c *Container) {
	sa, exists := this.GetSafeApp(c.AppName)
	if exists {
		sa.AddContainer(c)
	} else {
		sa = NewSafeApp()
		sa.AddContainer(c)
		this.AddSafeApp(c.AppName, sa)
	}
}

func (this *SafeRealState) AddSafeApp(name string, a *SafeApp) {
	this.Lock()
	defer this.Unlock()
	this.M[name] = a
}

func (this *SafeApp) ContainerExists(c *Container) bool {
	key := MakeContainerKey(c.Ip, c.Id)
	this.RLock()
	defer this.RUnlock()
	_, ok := this.M[key]
	return ok
}

func (this *SafeRealState) RealAppExists(name string) bool {
	this.RLock()
	defer this.RUnlock()
	_, ok := this.M[name]
	return ok
}

func (this *SafeRealState) Keys() []string {
	this.RLock()
	defer this.RUnlock()
	size := len(this.M)
	L := make([]string, size)
	i := 0
	for k, _ := range this.M {
		L[i] = k
		i++
	}
	return L
}

func (this *SafeRealState) DeleteByIp(ip string) {
	appNames := this.Keys()
	for _, name := range appNames {
		sa, ok := this.GetSafeApp(name)
		if !ok {
			continue
		}

		sa.DeleteByIp(ip)
	}
}

func (this *SafeApp) DeleteByIp(ip string) {
	needDelete := make([]string, 0)
	this.RLock()
	for _, c := range this.M {
		if c.Ip == ip {
			needDelete = append(needDelete, MakeContainerKey(c.Ip, c.Id))
		}
	}
	this.RUnlock()

	if len(needDelete) == 0 {
		return
	}

	this.Lock()
	for _, containerKey := range needDelete {
		delete(this.M, containerKey)
		this.ToUpdate = true
	}
	this.Unlock()
}

func (this *SafeRealState) GetSafeApp(name string) (*SafeApp, bool) {
	this.RLock()
	defer this.RUnlock()
	sa, ok := this.M[name]
	return sa, ok
}

func (this *SafeRealState) DeleteSafeApp(name string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, name)
}

func (this *SafeRealState) DeleteStale(before int64) {
	this.RLock()
	defer this.RUnlock()
	for _, sa := range this.M {
		sa.DeleteStale(before)
	}
}

func (this *SafeApp) DeleteStale(before int64) {
	needDelete := make([]string, 0)
	this.RLock()
	for _, c := range this.M {
		if c.UpdateAt < before {
			needDelete = append(needDelete, MakeContainerKey(c.Ip, c.Id))
		}
	}
	this.RUnlock()

	if len(needDelete) == 0 {
		return
	}

	this.Lock()
	for _, containerKey := range needDelete {
		delete(this.M, containerKey)
		this.ToUpdate = true
	}
	this.Unlock()
}

func (this *SafeApp) DeleteContainer(c *Container) {
	key := MakeContainerKey(c.Ip, c.Id)
	this.Lock()
	defer this.Unlock()
	delete(this.M, key)
	this.ToUpdate = true
}

func (this *SafeApp) HasRelation(ip string) bool {
	this.RLock()
	defer this.RUnlock()
	if len(this.M) == 0 {
		return false
	}

	for _, c := range this.M {
		if c.Ip == ip {
			return true
		}
	}

	return false
}

func (this *SafeRealState) HasRelation(appName, ip string) bool {
	sa, exists := this.GetSafeApp(appName)
	if !exists {
		return false
	}

	return sa.HasRelation(ip)
}
