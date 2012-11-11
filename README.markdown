= go-dtrace

If you want to fire DTrace probes from within your Go program,
go-dtrace is for you. It wraps libusdt to provide DTrace probe
definition in go.

== compiling libusdt

In order to use go-dtrace, *libusdt* needs to be compiled as a shared
library. To achieve this in Mac OS X:

- `$ git clone https://github.com/chrisa/libusdt.git`
- `$ cd libusdt`
- `$ make`
- `$ rm usdt_tracepoints.o`
- `$ cc -dynamiclib -install_name /usr/local/lib/libusdt.dylib -flat_namespace -o libusdt.dylib *.o`
- `$ mkdir -p /usr/local/lib`
- `$ cp libusdt.dylib /usr/local/lib/libusdt.dylib`

Have instructions for another OS? Add 'em here!

== go install

Once *libusdt* is installed, you can install go-dtrace from the command line:

- `$ go get github.com/ecin/go-dtrace`

== example

```go
package main

import (
  "github.com/ecin/go-dtrace"
  "reflect"
  "time"
  "fmt"
)

func main() {
  provider := dtrace.NewProvider("go-dtrace", "example")

  probe := provider.AddProbe("fire", "1", reflect.Int, reflect.String)
  fmt.Printf("%s:%s:%s:%s is ready to fire!\n", provider.Name, provider.Module, probe.Function, probe.Name)

  // Enable the provider AFTER defining probes
  provider.Enable()

  for i := 0;; i += 1 {
    probe.Fire(i, "Boom!")
    time.Sleep(1 * time.Second)
  }
}
```

Run `sudo dtrace -n "go-dtrace*:::"` to see your DTrace probe fire at 1 glorious Hz.
