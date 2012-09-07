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
import "errors"
import "fmt"

type Provider struct {
  Name string
  Module string
  Probes []Probe
  provider_t *C.usdt_provider_t
}

func (provider Provider) Error() (errMsg string) {
  errMsg = C.GoString(provider.provider_t.error)
  if (errMsg != "") {
    errMsg = fmt.Sprintf("dtrace: [%s] %s", provider.String(), errMsg)
  }

  return
}

func (provider Provider) String() string {
  return fmt.Sprintf("%s:%s", provider.Name, provider.Module)
}

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

  cTypes := make([]*C.char, len(signature) + 1)

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

func (provider Provider) Enable() (err error) {
  errCode := C.usdt_provider_enable(provider.provider_t)
  if errCode != 0 {
    err = errors.New(provider.Error())
  }
  return
}

func (provider Provider) IsEnabled() bool {
  return C.int(provider.provider_t.enabled) != 0
}
