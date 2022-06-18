package ctk

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/sync"
	"github.com/go-curses/cdk/log"
)

const TypeAccelMap cdk.CTypeTag = "ctk-accel-map"

var (
	ctkAccelMaps     = make(map[uuid.UUID]*CAccelMap)
	ctkAccelMapsLock = &sync.RWMutex{}
)

func init() {
	_ = cdk.TypesManager.AddType(TypeAccelMap, nil)
}

// AccelMap Hierarchy:
//	Object
//    +- AccelMap
//
// The AccelMap is the global accelerator mapping object. There is only one
// AccelMap instance for any given Display.
type AccelMap interface {
	Object

	Init() (already bool)
	AddEntry(accelPath string, accelKey cdk.Key, accelMods cdk.ModMask)
	LookupEntry(accelPath string) (accelerator Accelerator, ok bool)
	ChangeEntry(accelPath string, accelKey cdk.Key, accelMods cdk.ModMask, replace bool) (ok bool)
	Load(fileName string)
	LoadFromString(accelMap string)
	Save(fileName string)
	SaveToString() (accelMap string)
	LockPath(accelPath string)
	UnlockPath(accelPath string)
}

// The CAccelMap structure implements the AccelMap interface and is exported to
// facilitate type embedding with custom implementations. No member variables
// are exported as the interface methods are the only intended means of
// interacting with AccelMap objects.
type CAccelMap struct {
	CObject

	accelerators map[string]Accelerator
}

// GetAccelMap is the getter for the current Application AccelMap singleton.
// Returns nil if there is no Application present for the current thread.
func GetAccelMap() AccelMap {
	if acd, err := cdk.GetLocalContext(); err != nil {
		log.Error(err)
	} else {
		if app, ok := acd.Data.(Application); ok {
			return app.AccelMap()
		}
	}
	return nil
}

// Init initializes an AccelMap object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the AccelMap instance. Init is used in the
// NewAccelMap constructor and only necessary when implementing a derivative
// AccelMap type.
func (a *CAccelMap) Init() (already bool) {
	if a.InitTypeItem(TypeAccelMap, a) {
		return true
	}
	a.CObject.Init()
	a.accelerators = make(map[string]Accelerator)
	return false
}

// AddEntry registers a new accelerator with the global accelerator map. This
// method should only be called once per accel_path with the canonical accel_key
// and accel_mods for this path. To change the accelerator during runtime
// programmatically, use AccelMap.ChangeEntry(). The accelerator path must
// consist of "<WINDOWTYPE>/Category1/Category2/.../Action", where <WINDOWTYPE>
// should be a unique application-specific identifier, that corresponds to the
// kind of window the accelerator is being used in, e.g. "Gimp-Image",
// "Abiword-Document" or "Gnumeric-Settings". The Category1/.../Action portion
// is most appropriately chosen by the action the accelerator triggers, i.e. for
// accelerators on menu items, choose the item's menu path, e.g. "File/Save As",
// "Image/View/Zoom" or "Edit/Select All". So a full valid accelerator path may
// look like: "<Gimp-Toolbox>/File/Dialogs/Tool Options...".
//
// Parameters
// 	accel_path	valid accelerator path
// 	accel_key	the accelerator key
// 	accel_mods	the accelerator modifiers
func (a *CAccelMap) AddEntry(accelPath string, accelKey cdk.Key, accelMods cdk.ModMask) {
	a.RLock()
	if _, ok := a.accelerators[accelPath]; ok {
		a.RUnlock()
		log.ErrorDF(1, "accelerator exists for path: %v", accelPath)
		return
	}
	a.RUnlock()
	a.Lock()
	a.accelerators[accelPath] = NewAccelerator(accelPath, accelKey, accelMods)
	a.Unlock()
}

// LookupEntry returns the accelerator entry for accel_path.
//
// Parameters
// 	accel_path	a valid accelerator path
func (a *CAccelMap) LookupEntry(accelPath string) (accelerator Accelerator, ok bool) {
	a.RLock()
	if accelerator, ok = a.accelerators[accelPath]; !ok {
		accelerator = nil
	}
	a.RUnlock()
	return
}

// ChangeEntry updates the accel_key and accel_mods currently associated with
// accel_path. Due to conflicts with other accelerators, a change may not always
// be possible, replace indicates whether other accelerators may be deleted to
// resolve such conflicts. A change will only occur if all conflicts could be
// resolved (which might not be the case if conflicting accelerators are
// locked). Successful changes are indicated by a TRUE return value.
//
// Parameters
// 	accel_path	a valid accelerator path
// 	accel_key	the new accelerator key
// 	accel_mods	the new accelerator modifiers
// 	replace	TRUE if other accelerators may be deleted upon conflicts
func (a *CAccelMap) ChangeEntry(accelPath string, accelKey cdk.Key, accelMods cdk.ModMask, replace bool) (ok bool) {
	// get the Accelerator for accelPath
	// find any other association of accelKey + accelMods
	// if the found is locked and not replace, nok
	// else if the accelPath is locked, nok
	var foundOtherAccel Accelerator
	for _, accel := range a.accelerators {
		if accel.Match(accelKey, accelMods) {
			foundOtherAccel = accel
			break
		}
	}
	if foundOtherAccel != nil {
		if foundOtherAccel.IsLocked() && !replace {
			// existing key+mods assignment is treated immutable
			return false
		}
		// unset the existing key+mods assignment
		foundOtherAccel.UnsetKeyMods()
	}
	if foundPathAccel, pok := a.accelerators[accelPath]; pok {
		foundPathAccel.Configure(accelKey, accelMods)
		return true
	}
	a.LogError("accelPath not found: %v", accelPath)
	return false
}

// Load parses a file previously saved with AccelMap.Save() for accelerator
// specifications, and propagates them accordingly.
//
// Parameters
// 	file_name	a file containing accelerator specifications
func (a *CAccelMap) Load(fileName string) {

}

var rxAccelMapLineParser = regexp.MustCompile(`^\s*<([^>]+?)>/((?:[A-Z][- a-zA-Z\d]+?[a-zA-Z-\d]|/)+)\s*=\s*(.+?)\s*$`)

func (a *CAccelMap) LoadFromString(accelMap string) {
	parsed := make(map[string]string)
	for _, line := range strings.Split(accelMap, "\n") {
		line = strings.TrimSpace(line)
		if rxAccelMapLineParser.MatchString(line) {
			m := rxAccelMapLineParser.FindAllStringSubmatch(line, -1)
			if len(m) == 1 && len(m[0]) == 4 {
				p := fmt.Sprintf("<%s>/%s", m[0][1], m[0][2])
				parsed[p] = m[0][3]
			}
		}
	}
	for path, keyMods := range parsed {
		if key, mods, err := cdk.ParseKeyMods(keyMods); err != nil {
			a.LogErr(err)
		} else {
			a.Lock()
			if accelerator, ok := a.accelerators[path]; ok {
				accelerator.Configure(key, mods)
			} else {
				a.accelerators[path] = NewAccelerator(path, key, mods)
			}
			a.Unlock()
		}
	}
}

// Save stores the current accelerator specifications (accelerator path, key and
// modifiers) to file_name. The file is written in a format suitable to be read
// back in by AccelMap.Load().
//
// Parameters
// 	file_name	the name of the file to contain accelerator specifications
func (a *CAccelMap) Save(fileName string) {

}

func (a *CAccelMap) SaveToString() (accelMap string) {
	a.RLock()
	for path, accelerator := range a.accelerators {
		_, key, mods := accelerator.Settings()
		accelMap += fmt.Sprintf("%v = %v%v\n", path, mods.String(), key)
	}
	a.RUnlock()
	return
}

// LockPath locks the given accelerator path. If the accelerator map doesn't yet
// contain an entry for accel_path, a new one is created.
//
// Locking an accelerator path prevents its accelerator from being changed
// during runtime. A locked accelerator path can be unlocked by
// AccelMap.UnlockPath(). Refer to AccelMap.ChangeEntry() for information about
// runtime accelerator changes.
//
// If called more than once, accel_path remains locked until
// AccelMap.UnlockPath() has been called an equivalent number of times.
//
// Note that locking of individual accelerator paths is independent of locking
// the AccelGroup containing them. For runtime accelerator changes to be
// possible both the accelerator path and its AccelGroup have to be unlocked.
//
// Parameters
// 	accel_path	a valid accelerator path
func (a *CAccelMap) LockPath(accelPath string) {
	a.Lock()
	if accelerator, ok := a.accelerators[accelPath]; !ok {
		a.accelerators[accelPath] = NewDefaultAccelerator(accelPath)
	} else {
		accelerator.Lock()
	}
	a.Unlock()
}

// UnlockPath undoes the last call to AccelMap.LockPath() on this accel_path.
// Refer to AccelMap.LockPath() for information about accelerator path locking.
//
// Parameters
// 	accel_path	a valid accelerator path
func (a *CAccelMap) UnlockPath(accelPath string) {

}