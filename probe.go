package dtrace

/*
#cgo LDFLAGS: -lusdt -L/usr/local/lib
#include <stdlib.h>
#include "usdt.h"
*/
import "C"
import "unsafe"
import "fmt"

type Probe struct {
  Function string
  Name string
  probedef_t *C.usdt_probedef_t
}

func (probe Probe) String() string {
  return fmt.Sprintf("%s:%s", probe.Function, probe.Name)
}

func (probe Probe) IsEnabled() (enabled bool) {
  cProbe := (*C.usdt_probe_t)(probe.probedef_t.probe)

  if C.int(C.usdt_is_enabled(cProbe)) == 1 {
    enabled = true
  } else {
    enabled = false
  }

  return
}

// Could use reflect package to throw ArgumentError
// if args don't match the probe's probedef
func (probe Probe) Fire(args ...interface{}) {
  nargv := make([]unsafe.Pointer, len(args) + 1)
  argc := probe.probedef_t.argc

  for i, arg := range args {
    if i > int(argc) {
      break
    }

    switch arg.(type) {
      case int:
        x := C.int(arg.(int))
        nargv[i] = unsafe.Pointer(uintptr(x))
      case string:
        nargv[i] = unsafe.Pointer(C.CString(arg.(string)))
      default:
        nargv[i] = nil
    }
  }

  cProbe := (*C.usdt_probe_t)(probe.probedef_t.probe)
  C.usdt_fire_probe(cProbe, argc, &nargv[0])
}
