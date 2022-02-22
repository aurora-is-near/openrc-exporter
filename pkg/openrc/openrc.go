package openrc

// #cgo LDFLAGS: -lrc
/*
#include <rc.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

const (
	Prefix     = C.RC_PREFIX
	SysconfDir = C.RC_SYSCONFDIR
	LibDir     = C.RC_LIBDIR
	LibexecDir = C.RC_LIBEXECDIR

	SvcDir      = C.RC_SVCDIR
	RunlevelDir = C.RC_RUNLEVELDIR
	InitDir     = C.RC_INITDIR
	ConfDir     = C.RC_CONFDIR
	PluginDir   = C.RC_PLUGINDIR

	InitFifo     = C.RC_INIT_FIFO
	ProfileEnv   = C.RC_PROFILE_ENV
	SysWhitelist = C.RC_SYS_WHITELIST
	UsrWhitelist = C.RC_USR_WHITELIST
	Conf         = C.RC_CONF
	ConfD        = C.RC_CONF_D

	// Service states
	ServiceStopped  = C.RC_SERVICE_STOPPED
	ServiceStarted  = C.RC_SERVICE_STARTED
	ServiceStopping = C.RC_SERVICE_STOPPING
	ServiceStarting = C.RC_SERVICE_STARTING
	ServiceInactive = C.RC_SERVICE_INACTIVE

	ServiceHotplugged = C.RC_SERVICE_HOTPLUGGED

	ServiceFailed      = C.RC_SERVICE_FAILED
	ServiceScheduled   = C.RC_SERVICE_SCHEDULED
	ServiceWasInactive = C.RC_SERVICE_WASINACTIVE
	ServiceCrashed     = C.RC_SERVICE_CRASHED
)

var (
	ServiceStateNames = map[int]string{
		ServiceStarted:     "started",
		ServiceStopped:     "stopped",
		ServiceStarting:    "starting",
		ServiceStopping:    "stopping",
		ServiceInactive:    "inactive",
		ServiceWasInactive: "wasinactive",
		ServiceHotplugged:  "hotplugged",
		ServiceFailed:      "failed",
		ServiceScheduled:   "scheduled",
		ServiceCrashed:     "crashed",
	}
)

func RunlevelGet() string {
	runlevelCharPtr := C.rc_runlevel_get()
	defer C.free(unsafe.Pointer(runlevelCharPtr))

	return C.GoString(runlevelCharPtr)
}

func RunlevelExists(runlevel string) bool {
	runlevelCharPtr := C.CString(runlevel)
	defer C.free(unsafe.Pointer(runlevelCharPtr))

	result := C.rc_runlevel_exists(runlevelCharPtr)
	return bool(result)
}

func RunlevelStacks(runlevel string) []string {
	runlevelCharPtr := C.CString(runlevel)
	defer C.free(unsafe.Pointer(runlevelCharPtr))

	stackStringlist := C.rc_runlevel_stacks(runlevelCharPtr)
	defer C.rc_stringlist_free(stackStringlist)

	return goStringList(stackStringlist)
}

func RunlevelList() []string {
	stringList := C.rc_runlevel_list()
	defer C.rc_stringlist_free(stringList)

	return goStringList(stringList)
}

func RunlevelSet(runlevel string) bool {
	runlevelCharPtr := C.CString(runlevel)
	defer C.free(unsafe.Pointer(runlevelCharPtr))

	result := C.rc_runlevel_set(runlevelCharPtr)
	return bool(result)
}

func RunlevelStarting() bool {
	return bool(C.rc_runlevel_starting())
}

func RunlevelStopping() bool {
	return bool(C.rc_runlevel_stopping())
}

func ServiceAdd(runlevel string, service string) bool {
	runlevelCharPtr := C.CString(runlevel)
	defer C.free(unsafe.Pointer(runlevelCharPtr))

	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	result := C.rc_service_add(runlevelCharPtr, serviceCharPtr)
	return bool(result)
}

func ServiceDelete(runlevel string, service string) bool {
	runlevelCharPtr := C.CString(runlevel)
	defer C.free(unsafe.Pointer(runlevelCharPtr))

	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	result := C.rc_service_delete(runlevelCharPtr, serviceCharPtr)
	return bool(result)
}

func ServiceDescription(service string, option *string) string {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	optionCharPtr := (*C.char)(nil)
	if option != nil {
		optionCharPtr = C.CString(*option)
		defer C.free(unsafe.Pointer(optionCharPtr))
	}

	result := C.rc_service_description(serviceCharPtr, optionCharPtr)
	return C.GoString(result)
}

func ServiceExists(service string) bool {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	return bool(C.rc_service_exists(serviceCharPtr))
}

func ServiceInRunLevel(service string, runlevel string) bool {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	runlevelCharPtr := C.CString(runlevel)
	defer C.free(unsafe.Pointer(runlevelCharPtr))

	result := C.rc_service_in_runlevel(serviceCharPtr, runlevelCharPtr)
	return bool(result)
}

func ServiceMark(service string, state int) bool {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	result := C.rc_service_mark(serviceCharPtr, C.RC_SERVICE(state))
	return bool(result)
}

func ServiceExtraCommands(service string) []string {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	stringList := C.rc_service_extra_commands(serviceCharPtr)
	defer C.rc_stringlist_free(stringList)

	return goStringList(stringList)
}

func ServiceResolve(service string) string {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	fullPath := C.rc_service_resolve(serviceCharPtr)
	defer C.free(unsafe.Pointer(fullPath))

	return C.GoString(fullPath)
}

func ServiceState(service string) int {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	return int(C.rc_service_state(serviceCharPtr))
}

func ServiceValueGet(service string, option string) string {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	optionCharPtr := C.CString(option)
	defer C.free(unsafe.Pointer(optionCharPtr))

	value := C.rc_service_value_get(serviceCharPtr, optionCharPtr)
	defer C.free(unsafe.Pointer(value))

	return C.GoString(value)
}

func ServiceValueSet(service, option, value string) bool {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	optionCharPtr := C.CString(option)
	defer C.free(unsafe.Pointer(optionCharPtr))

	valueCharPtr := C.CString(value)
	defer C.free(unsafe.Pointer(valueCharPtr))

	result := C.rc_service_value_set(serviceCharPtr, optionCharPtr, valueCharPtr)
	return bool(result)
}

func ServicesInRunlevel(runlevel *string) []string {
	runlevelCharPtr := (*C.char)(nil)
	if runlevel != nil {
		runlevelCharPtr = C.CString(*runlevel)
		defer C.free(unsafe.Pointer(runlevelCharPtr))
	}

	stringList := C.rc_services_in_runlevel(runlevelCharPtr)
	defer C.rc_stringlist_free(stringList)

	return goStringList(stringList)
}

func ServicesInRunlevelStacked(runlevel string) []string {
	runlevelCharPtr := C.CString(runlevel)
	defer C.free(unsafe.Pointer(runlevelCharPtr))

	stringList := C.rc_services_in_runlevel_stacked(runlevelCharPtr)
	defer C.rc_stringlist_free(stringList)

	return goStringList(stringList)
}

func ServicesInState(state int) []string {
	stringList := C.rc_services_in_state(C.RC_SERVICE(state))
	defer C.rc_stringlist_free(stringList)

	return goStringList(stringList)
}

func ServicesScheduled(service string) []string {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	stringList := C.rc_services_scheduled(serviceCharPtr)
	defer C.rc_stringlist_free(stringList)

	return goStringList(stringList)
}

func ServiceDaemonsCrashed(service string) bool {
	serviceCharPtr := C.CString(service)
	defer C.free(unsafe.Pointer(serviceCharPtr))

	result := C.rc_service_daemons_crashed(serviceCharPtr)
	return bool(result)
}
