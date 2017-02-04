package elfpatch

import (
	"debug/elf"
	"os"
	"os/exec"
	"path"
	"strings"
)

type ElfObject struct {
	FilePath   string
	originPath []string
	DynLibs    []*ElfObject
	InterpPath string
	DynPathMap map[string]string
	SoName     string
}

func (elfobject *ElfObject) Export(path string) {
}

func listDyn(loader string, execute string) map[string]string {
	dynMap := make(map[string]string)
	output, err := exec.Command(loader, "--list", execute).Output()
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(output), "\n") {
		trimLine := strings.Trim(line, "\t\n")
		segments := strings.Split(trimLine, " ")
		if len(segments) > 2 {
			dyn := segments[0]
			dynPath := segments[2]
			dynMap[dyn] = dynPath
		} else if path.IsAbs(segments[0]) {
			dynPath := segments[0]
			dyn := path.Base(dynPath)
			dynMap[dyn] = dynPath
		}
	}

	return dynMap
}

func (elfObject *ElfObject) Scan() {
	f, err := os.Open(elfObject.FilePath)
	if err != nil {
		panic(err)
	}

	_elf, err := elf.NewFile(f)
	if err != nil {
		panic(err)
	}

	importLibraries, err := _elf.ImportedLibraries()
	if err != nil {
		panic(err)
	}

	if elfObject.InterpPath == "" {
		interp := ""

		for _, p := range _elf.Progs {
			if p.Type == elf.PT_INTERP {
				// note: sub tailing zero
				buf := make([]byte, p.Filesz-1)
				p.ReadAt(buf, 0)
				interp = string(buf)
			}
		}

		if interp == "" {
			panic("no interp")
		}

		elfObject.InterpPath = interp
	}

	if elfObject.DynPathMap == nil {
		elfObject.DynPathMap = listDyn(elfObject.InterpPath, elfObject.FilePath)
	}

	for _, lib := range importLibraries {
		libPath := elfObject.DynPathMap[lib]
		depElfObject := new(ElfObject)
		depElfObject.FilePath = libPath
		depElfObject.InterpPath = elfObject.InterpPath
		depElfObject.Scan()
		elfObject.DynLibs = append(elfObject.DynLibs, depElfObject)
	}
}

func AppendString(slice []string, data string) []string {
	m := len(slice)
	if cap(slice) < (m + 1) {
		newSlice := make([]string, m*2+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	// set len as used
	slice = slice[0 : m+1]

	slice[m] = data
	return slice
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getDynPathList(pathList *[]string, elfObject *ElfObject) {
	for _, dynlib := range elfObject.DynLibs {
		if !contains(*pathList, dynlib.FilePath) {
			*pathList = AppendString(*pathList, dynlib.FilePath)
		}
	}

	for _, dynlib := range elfObject.DynLibs {
		if dynlib.DynLibs != nil {
			getDynPathList(pathList, dynlib)
		}
	}
}

func (elfObject *ElfObject) GetDynPathList() []string {
	pathList := new([]string)
	getDynPathList(pathList, elfObject)
	return *pathList
}
