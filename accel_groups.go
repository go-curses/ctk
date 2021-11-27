package ctk

import (
	"sync"

	"github.com/go-curses/cdk"
	"github.com/gofrs/uuid"
)

var (
	ctkAccelGroups     = make(map[uuid.UUID][]AccelGroup)
	ctkAccelGroupsLock = &sync.RWMutex{}
)

type AccelGroups interface {
	Object

	AddAccelGroup(object Object, accelGroup AccelGroup)
	RemoveAccelGroup(object Object, accelGroup AccelGroup)
	Activate(object Object, accelKey cdk.Key, accelMods cdk.ModMask) (value bool)
	FromObject(object Object) (groups []AccelGroup)
}

type CAccelGroups struct {
	CObject
}

func NewAccelGroups() (groups AccelGroups) {
	groups = new(CAccelGroups)
	groups.Init()
	return
}

func (a *CAccelGroups) AddAccelGroup(object Object, accelGroup AccelGroup) {
	oid := object.ObjectID()
	a.Lock()
	ctkAccelGroupsLock.Lock()
	if _, ok := ctkAccelGroups[oid]; !ok {
		ctkAccelGroups[oid] = make([]AccelGroup, 0)
	}
	ctkAccelGroups[oid] = append(ctkAccelGroups[oid], accelGroup)
	ctkAccelGroupsLock.Unlock()
	a.Unlock()
}

func (a *CAccelGroups) RemoveAccelGroup(object Object, accelGroup AccelGroup) {
	oid := object.ObjectID()
	agid := accelGroup.ObjectID()
	a.Lock()
	ctkAccelGroupsLock.Lock()
	index := -1
	if _, ok := ctkAccelGroups[oid]; ok {
		for idx, ag := range ctkAccelGroups[oid] {
			if ag.ObjectID() == agid {
				index = idx
				break
			}
		}
		if index > -1 {
			ctkAccelGroups[oid] = append(
				ctkAccelGroups[oid][:index],
				ctkAccelGroups[oid][index+1:]...,
			)
		}
	}
	ctkAccelGroupsLock.Unlock()
	a.Unlock()
}

func (a *CAccelGroups) Activate(object Object, accelKey cdk.Key, accelMods cdk.ModMask) (value bool) {
	for _, ag := range a.FromObject(object) {
		if ag.AccelGroupActivate(accelKey, accelMods) {
			break
		}
	}
	return
}

func (a *CAccelGroups) FromObject(object Object) (groups []AccelGroup) {
	oid := object.ObjectID()
	a.RLock()
	ctkAccelGroupsLock.RLock()
	var ok bool
	if groups, ok = ctkAccelGroups[oid]; !ok {
		ctkAccelGroupsLock.RUnlock()
		a.RUnlock()
		return []AccelGroup{}
	}
	ctkAccelGroupsLock.RUnlock()
	a.RUnlock()
	return
}
