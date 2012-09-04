package dtrace

/*
Compile libusdt with:
cc -dynamiclib -install_name /usr/local/lib/libusdt.dylib -flat_namespace -o libusdt.dylib *.o
*/

/*
#cgo LDFLAGS: -lusdt -L/usr/local/lib
#include <stdlib.h>
#include "usdt.h"
*/
import "C"
import "unsafe"
import "reflect"

type Probe struct {
  Function string
  Name string
  probedef_t *C.usdt_probedef_t
}

type Provider struct {
  Name string
  Module string
  Probes []Probe
  provider_t *C.usdt_provider_t
}
/*
// Probably look into defining new errors in Go
type Error struct {

}
*/



func NewProvider(name string, module string) (provider *Provider) {
  cName := C.CString(name)
  cModule := C.CString(module)
  defer C.free(unsafe.Pointer(cName))
  defer C.free(unsafe.Pointer(cModule))

  var probes []Probe

  provider = &Provider{
    name,
    module,
    probes,
    C.usdt_create_provider(cName, cModule),
  }

  return
}

func (provider *Provider) AddProbe(function string, name string, signature ...reflect.Kind) Probe {
  cFunction := C.CString(function)
  cName := C.CString(name)
  defer C.free(unsafe.Pointer(cFunction))
  defer C.free(unsafe.Pointer(cName))

  cTypes := make([]*C.char, len(signature))

  for i, kind := range signature {
    switch kind {
      case reflect.Int:
        cTypes[i] = C.CString("int")
        defer C.free(unsafe.Pointer(cTypes[i]))
      case reflect.String:
        cTypes[i] = C.CString("char *")
        defer C.free(unsafe.Pointer(cTypes[i]))
      default:
        cTypes[i] = nil
    }
  }

  probedef := C.usdt_create_probe(cFunction, cName, C.size_t(len(signature)), &cTypes[0])
  C.usdt_provider_add_probe(provider.provider_t, probedef)

  newProbe := Probe{
    function,
    name,
    probedef,
  }

  provider.Probes = append(provider.Probes, newProbe)

  return newProbe
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

// Missing error handling
func (provider Provider) Enable() {
  C.usdt_provider_enable(provider.provider_t)
}

func (provider Provider) IsEnabled() bool {
  return C.int(provider.provider_t.enabled) == 1
}

// Could use reflect package to throw ArgumentError
// if args don't match the probe's probedef
func (probe Probe) fire(args ...interface{}) {
  nargv := make([]unsafe.Pointer, len(args))
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
