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

package enums

import (
	"strconv"
	"strings"
)

type ResponseType int

const (
	ResponseNone        ResponseType = -1
	ResponseReject      ResponseType = -2
	ResponseAccept      ResponseType = -3
	ResponseDeleteEvent ResponseType = -4
	ResponseOk          ResponseType = -5
	ResponseCancel      ResponseType = -6
	ResponseClose       ResponseType = -7
	ResponseYes         ResponseType = -8
	ResponseNo          ResponseType = -9
	ResponseApply       ResponseType = -10
	ResponseHelp        ResponseType = -11
)

var (
	responseTypes = map[ResponseType]string{
		ResponseNone:        "none",
		ResponseReject:      "reject",
		ResponseAccept:      "accept",
		ResponseDeleteEvent: "delete-event",
		ResponseOk:          "ok",
		ResponseCancel:      "cancel",
		ResponseClose:       "close",
		ResponseYes:         "yes",
		ResponseNo:          "no",
		ResponseApply:       "apply",
		ResponseHelp:        "help",
	}
	responseNames = map[string]ResponseType{
		"none":         ResponseNone,
		"reject":       ResponseReject,
		"accept":       ResponseAccept,
		"delete-event": ResponseDeleteEvent,
		"ok":           ResponseOk,
		"cancel":       ResponseCancel,
		"close":        ResponseClose,
		"yes":          ResponseYes,
		"no":           ResponseNo,
		"apply":        ResponseApply,
		"help":         ResponseHelp,
	}
)

func (r ResponseType) String() string {
	if v, ok := responseTypes[r]; ok {
		return v
	}
	return strconv.Itoa(int(r))
}

func ResponseTypeFromName(name string) ResponseType {
	name = strings.ToLower(name)
	if v, ok := responseNames[name]; ok {
		return v
	}
	return ResponseNone
}