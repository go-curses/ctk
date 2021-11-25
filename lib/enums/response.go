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
