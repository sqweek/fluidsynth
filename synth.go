package fluidsynth

// #cgo pkg-config: fluidsynth
// #include <fluidsynth.h>
// #include <stdlib.h>
import "C"
import "reflect"

import (
	"os"
	"fmt"
	"unsafe"
)

type Synth struct {
	csettings *C.fluid_settings_t
	csynth *C.fluid_synth_t
	cdriver *C.fluid_audio_driver_t
}

func NewSynth(settings map[string]interface{}) *Synth {
	csettings, _ := C.new_fluid_settings()
	for key, value := range settings {
		ckey := C.CString(key)
		switch value := value.(type) {
		case string:
			cval := C.CString(value)
			C.fluid_settings_setstr(csettings, ckey, cval)
			C.free(unsafe.Pointer(cval))
		case int:
			C.fluid_settings_setint(csettings, ckey, C.int(value))
		case float64:
			C.fluid_settings_setnum(csettings, ckey, C.double(value))
		default:
			fmt.Fprintf(os.Stderr, "NewSynth: ignoring setting %s: unhandled type %T\n", key, value)
		}
		C.free(unsafe.Pointer(ckey))
	}
	csynth := C.new_fluid_synth(csettings)
	//cdriver := C.new_fluid_audio_driver(csettings, csynth)
	return &Synth{csettings, csynth, nil}
}

func (s *Synth) SFLoad(path string, resetPresets bool) int {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	creset := C.int(1)
	if !resetPresets {
		creset = 0
	}
	cfont_id, _ := C.fluid_synth_sfload(s.csynth, cpath, creset)
	return int(cfont_id)
}

/* XXX can this be run automatically on gc? */
func (s *Synth) Delete() {
	//C.delete_fluid_audio_driver(s.cdriver)
	C.delete_fluid_synth(s.csynth)
	C.delete_fluid_settings(s.csettings)
}

func (s *Synth) NoteOn(channel, note, velocity uint8) {
	C.fluid_synth_noteon(s.csynth, C.int(channel), C.int(note), C.int(velocity))
}

func (s *Synth) NoteOff(channel, note uint8) {
	C.fluid_synth_noteoff(s.csynth, C.int(channel), C.int(note))
}

func (s *Synth) ProgramChange(channel, program uint8) {
	C.fluid_synth_program_change(s.csynth, C.int(channel), C.int(program))
}

func (s *Synth) WriteFrames_int16(dst []int16) {
	if len(dst) % 2 != 0 {
		panic("dst not disivible by 2")
	}
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&dst))
	cbuf := unsafe.Pointer(hdr.Data)
	C.fluid_synth_write_s16(s.csynth, C.int(len(dst)/2), cbuf, 0, 2, cbuf, 1, 2)
}

