package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCfg(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "devices.json")
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestDevicesInheritDefaults(t *testing.T) {
	t.Setenv("EOS_CONFIG", writeCfg(t, `{
	  "defaults": {"username": "u", "password": "p"},
	  "devices": [
	    {"name": "a", "host": "192.0.2.1"},
	    {"name": "b", "host": "192.0.2.2", "username": "other", "password": "x"}
	  ]}`))
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if a := c.Devices[0]; a.Username != "u" || a.Password != "p" {
		t.Errorf("device a did not inherit defaults: %+v", a)
	}
	if b := c.Devices[1]; b.Username != "other" || b.Password != "x" {
		t.Errorf("device-level credentials must win: %+v", b)
	}
}

func TestDeviceWithoutCredsAndNoDefaultsFails(t *testing.T) {
	t.Setenv("EOS_CONFIG", writeCfg(t, `{"devices":[{"name":"a","host":"192.0.2.1"}]}`))
	if _, err := Load(); err == nil {
		t.Fatal("expected error for device without credentials")
	}
}
