package fluidsynth

// #cgo pkg-config: fluidsynth
// #include <fluidsynth.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

var settingNames map[string]*C.char
var nSettings = 0

type Settings struct {
	ptr *C.fluid_settings_t
}

func NewSettings() Settings {
	if settingNames == nil {
		settingNames = make(map[string]*C.char)
	}
	nSettings++
	return Settings{ptr: C.new_fluid_settings()}
}

func (s *Settings) Delete() {
	nSettings--
	if nSettings == 0 {
		settingNames = nil
	}
	C.delete_fluid_settings(s.ptr)
}

func cname(name string) *C.char {
	if cname, ok := settingNames[name]; ok {
		return cname
	}
	cname := C.CString(name)
	settingNames[name] = cname
	return cname
}

/* IsRealtime returns true if changing the specified setting immediately affects an associated Synth */
func (s *Settings) IsRealtime(name string) bool {
	return C.fluid_settings_is_realtime(s.ptr, cname(name)) == 1
}

func (s *Settings) SetInt(name string, val int) bool {
	return C.fluid_settings_setint(s.ptr, cname(name), C.int(val)) == 1
}

func (s *Settings) SetNum(name string, val float64) bool {
	return C.fluid_settings_setnum(s.ptr, cname(name), C.double(val)) == 1
}

func (s *Settings) SetString(name, val string) bool {
	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))
	return C.fluid_settings_setstr(s.ptr, cname(name), cval) == 1
	
}

func (s *Settings) GetInt(name string, val *int) bool {
	return C.fluid_settings_getint(s.ptr, cname(name), (*C.int)(unsafe.Pointer(val))) == 1
}

func (s *Settings) GetNum(name string, val *float64) bool {
	return C.fluid_settings_getnum(s.ptr, cname(name), (*C.double)(unsafe.Pointer(val))) == 1
}

func (s *Settings) GetString(name string, val *string) bool {
	var cstr *C.char
	ok := (C.fluid_settings_getstr(s.ptr, cname(name), &cstr) == 1)
	if ok {
		*val = C.GoString(cstr)
	}
	return ok
}
