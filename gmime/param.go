package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"

type GMimeParamsCallback func(name string, value string)

type Parametrized interface {
	SetParameter(name, value string)
	Parameter(name string) string
	ForEachParam(callback GMimeParamsCallback)
}

func forEachParam(params *C.GMimeParam, callback GMimeParamsCallback) {
	for params != nil {
		cName := C.g_mime_param_get_name(params)
		name := C.GoString(cName)
		cValue := C.g_mime_param_get_value(params)
		value := C.GoString(cValue)
		callback(name, value)
		params = C.g_mime_param_next(params)
	}
}
