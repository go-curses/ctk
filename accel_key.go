// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctk

import (
	"fmt"

	"github.com/go-curses/cdk"
	"github.com/go-curses/ctk/lib/enums"
)

type AccelKey interface {
	GetKey() cdk.Key
	GetMods() cdk.ModMask
	GetFlags() enums.AccelFlags
	Match(key cdk.Key, mods cdk.ModMask) (match bool)
	String() (key string)
}

type CAccelKey struct {
	Key   cdk.Key
	Mods  cdk.ModMask
	Flags enums.AccelFlags
}

func NewAccelKey(key cdk.Key, mods cdk.ModMask, flags enums.AccelFlags) (accelKey AccelKey) {
	accelKey = &CAccelKey{
		Key:   key,
		Mods:  mods,
		Flags: flags,
	}
	return
}

func (a *CAccelKey) GetKey() cdk.Key {
	return a.Key
}

func (a *CAccelKey) GetMods() cdk.ModMask {
	return a.Mods
}

func (a *CAccelKey) GetFlags() enums.AccelFlags {
	return a.Flags
}

func (a *CAccelKey) Match(key cdk.Key, mods cdk.ModMask) (match bool) {
	k, m := a.GetKey(), a.GetMods()
	return key == k && mods == m
}

func (a *CAccelKey) String() (key string) {
	mods := a.Mods.String()
	name := cdk.LookupKeyName(a.Key)
	if len(mods) > 0 {
		return fmt.Sprintf("%v%v", mods, name)
	}
	return fmt.Sprintf("%v", name)
}