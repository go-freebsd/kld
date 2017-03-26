package kld

import (
	"testing"
)

func TestLoadFiles(t *testing.T) {
	files, err := LoadedFiles()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(files) < 1 {
		t.Error("LoadedFiles should contain at least the kernel")
	}

	foundKernel := false
	for _, f := range files {
		t.Log(f.String())
		if f.Name() == "kernel" {
			foundKernel = true
		}
	}

	if foundKernel != true {
		t.Error("LoadedFiles should contain the kernel file")
	}
}

func TestFindModule(t *testing.T) {
	path, err := Find("pf")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("path for pf: %s", path)
	}

	path, err = Find("asd")
	if path != "" || err == nil {
		t.Error(err)
	}
}

func TestLoadAndUnload(t *testing.T) {
	err := Load("aesni")
	if err != nil {
		t.Error(err)
	}

	loaded, err := Loaded("aesni")
	if !loaded {
		t.Error("aesni should have been loaded")
	}
	if err != nil {
		t.Error(err)
	}

	err = Unload("aesni", true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestLoadUnloadError(t *testing.T) {
	err := LoadFile("/tmp/foobar.ko")
	if err == nil {
		t.Error("loading this should fail")
	}

	err = UnloadFile(999999, false)
	if err == nil {
		t.Error("unloading this should fail")
	}
}
