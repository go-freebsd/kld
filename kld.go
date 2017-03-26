package kld

import (
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
	"unsafe"
)

// #include <stdlib.h>
// #include <sys/param.h>
// #include <sys/linker.h>
import "C"

// Extension kernel file / module extension string
const Extension = ".ko"

// LoadFile load the passed kernel module by .ko file path
func LoadFile(path string) error {
	cpath := C.CString(path)
	v, err := C.kldload(cpath)
	C.free(unsafe.Pointer(cpath))
	if v != 0 {
		return err
	}
	return nil
}

// Load load the passed kernel module by name
func Load(name string) error {
	path, err := Find(name)
	if err != nil {
		return err
	}
	return LoadFile(path)
}

// UnloadFile Unload the passed kernel module by .ko file path
func UnloadFile(id int, force bool) error {
	var v C.int
	var err error
	if force {
		v, err = C.kldunloadf(C.int(id), C.LINKER_UNLOAD_FORCE)
	} else {
		v, err = C.kldunload(C.int(id))
	}

	if v != 0 {
		return err
	}

	return nil
}

// Unload Unload the passed kernel module by name
func Unload(name string, force bool) error {
	files, err := LoadedFiles()
	if err != nil {
		return err
	}
	var file *File
	for _, f := range files {
		if f.Name() == name || f.Name() == name+Extension {
			file = f
			break
		}
	}
	if file == nil {
		return fmt.Errorf("Unkown module: %s(%s)", name, Extension)
	}
	return UnloadFile(file.ID(), force)
}

// Path will return the module path for the passed module name
// or an error if there is no module
func Find(name string) (string, error) {
	pathstr, err := syscall.Sysctl("kern.module_path")
	if err != nil {
		return "", nil
	}

	paths := strings.Split(pathstr, ";")
	for _, p := range paths {
		for _, n := range []string{name, name + Extension} {
			filepath := path.Join(p, n)
			info, err := os.Stat(filepath)
			if err == nil && info.IsDir() == false {
				return filepath, nil
			}
		}
	}

	return "", fmt.Errorf("Module %s(%s) not found in %s",
		name, Extension, pathstr)
}

// Loaded returns true if the kernel module is already loaded
func Loaded(name string) (bool, error) {
	path, err := Find(name)
	if err != nil {
		return false, err
	}

	files, err := LoadedFiles()
	if err != nil {
		return false, err
	}

	for _, f := range files {
		if f.Pathname() == path {
			return true, nil
		}
	}
	return false, nil
}

// File a kld file
type File struct {
	wrap C.struct_kld_file_stat
}

// The name of the file
func (f File) Name() string {
	return C.GoString(&f.wrap.name[0])
}

// The full name of the file file including the path
func (f File) Pathname() string {
	return C.GoString(&f.wrap.pathname[0])
}

// Size the amount of memory bytes allocated by the file
func (f File) Size() int {
	return int(f.wrap.size)
}

// Address the load address of the kld file
func (f File) Address() unsafe.Pointer {
	return unsafe.Pointer(f.wrap.address)
}

// Refs the number of modules referenced by this file
func (f File) Refs() int {
	return int(f.wrap.refs)
}

// ID the id of the file
func (f File) ID() int {
	return int(f.wrap.id)
}

func (f File) String() string {
	return fmt.Sprintf("%d: %s (Size: %d Address: %d Refs: %d Path: %s)",
		f.ID(), f.Name(), f.Size(), f.Address(), f.Refs(), f.Pathname())
}

// LoadedFiles returns currently loaded files
func LoadedFiles() ([]*File, error) {
	var files []*File
	for fileid := C.kldnext(C.int(0)); fileid != C.int(0); fileid = C.kldnext(fileid) {
		var f File
		f.wrap.version = C.int(unsafe.Sizeof(f.wrap))
		if v, err := C.kldstat(fileid, &f.wrap); v == -1 {
			return nil, err
		}
		files = append(files, &f)
	}
	return files, nil
}
