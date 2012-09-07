package dtrace

import (
  "testing"
  "reflect"
)

func TestNewProvider(t *testing.T) {
  provider := NewProvider("golang", "dtrace")

  if provider.Name != "golang" {
    t.Error("Provider constructor didn't assign name")
  }

  if provider.Module != "dtrace" {
    t.Error("Provider constructor didn't assign module")
  }
}

func TestAddProbe(t *testing.T) {
  provider := NewProvider("golang", "dtrace")

  provider.AddProbe("Hello", "world", reflect.Int, reflect.String)

  if len(provider.Probes) != 1 {
    t.Error("AddProbe didn't add a probe to the provider")
  }
}

func TestEnable(t *testing.T) {
  provider := NewProvider("golang", "dtrace")

  if err := provider.Enable(); err == nil {
    t.Error("Enabling a provider without probes should fail")
  }

  provider.AddProbe("Hello", "World", reflect.Int, reflect.Int)

  if provider.IsEnabled() {
    t.Error("Provider isn't disabled by default")
  }

  if err := provider.Enable(); err != nil {
    t.Error("Error enabling valid provider")
  }

  if !provider.IsEnabled() {
    t.Error("Couldn't enable Provider")
  }

  if err := provider.Enable(); err != nil {
    t.Error("Got fatal error when enabling an already enabled provider")
  }

  if provider.Error() == "" {
    t.Error("Didn't get a non-fatal error when enabling an already enabled provider")
  }
}

func TestFire(t *testing.T) {
  provider := NewProvider("golang", "dtrace")

  // Testing that we don't panic while firing probes
  defer func() {
    if x := recover(); x != nil {
      t.Errorf("Panic firing probes: %s", x)
    }
  }()

  probe1 := provider.AddProbe("Probe", "1")
  probe2 := provider.AddProbe("Probe", "2", reflect.Int)
  probe3 := provider.AddProbe("Probe", "3", reflect.Int, reflect.String)

  provider.Enable()

  probe1.fire()
  probe2.fire(1)
  probe3.fire(1, "lasers!")
}

func TestString(t *testing.T) {
  provider := NewProvider("golang", "dtrace")

  if provider.String() != "golang:dtrace" {
    t.Error("Provider String() didn't return name:module")
  }

  if probe := provider.AddProbe("Hello", "World"); probe.String() != "Hello:World" {
    t.Error("Probe String() didn't return function:name")
  }
}
