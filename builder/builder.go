package builder

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/linxiulei/sucker/analyzer"
	"github.com/linxiulei/sucker/elfpatch"
	"github.com/linxiulei/sucker/utils"
)

const DynDirPath string = "dynlib"
const InterpDirPath string = "interp"

type Builder struct {
}

type Packer struct {
	Stat      analyzer.Stat
	ElfObject *elfpatch.ElfObject
}

func (p *Packer) ExportDyn(destDir string) {
	destDynDir := path.Join(destDir, DynDirPath)
	if _, err := os.Stat(destDynDir); os.IsNotExist(err) {
		mode := os.FileMode(0755)
		os.Mkdir(destDynDir, mode)
	}

	for _, dynPath := range p.ElfObject.GetDynPathList() {
		dynBase := path.Base(dynPath)
		destDynPath := path.Join(destDynDir, dynBase)
		err := utils.CopyFile(dynPath, destDynPath, true)
		if err != nil {
			panic(err)
		}
	}
}

func (p *Packer) ExportExe(destDir string, execPrefix string) {
	executePath := path.Join(destDir, path.Base(p.ElfObject.FilePath))
	err := utils.CopyFile(p.ElfObject.FilePath, executePath, true)
	if err != nil {
		panic(err)
	}

	destDynDir := path.Join(execPrefix, DynDirPath)
	destDynDirAbs, _ := filepath.Abs(destDynDir)
	_, err = exec.Command(
		"sucker_patchelf",
		"--set-rpath",
		destDynDirAbs,
		executePath).Output()

	if err != nil {
		panic(err)
	}

	destInterpDir := path.Join(execPrefix, InterpDirPath)
	interpPath := path.Join(destInterpDir, path.Base(p.ElfObject.InterpPath))
	interpPathAbs, _ := filepath.Abs(interpPath)

	_, err = exec.Command(
		"sucker_patchelf",
		"--set-interpreter",
		interpPathAbs,
		executePath).Output()

	if err != nil {
		panic(err)
	}
}

func (p *Packer) ExportInterp(destDir string) {
	destInterpDir := path.Join(destDir, InterpDirPath)
	if _, err := os.Stat(destInterpDir); os.IsNotExist(err) {
		mode := os.FileMode(0755)
		os.Mkdir(destInterpDir, mode)
	}

	InterpPath := path.Join(destInterpDir, path.Base(p.ElfObject.InterpPath))
	err := utils.CopyFile(p.ElfObject.InterpPath, InterpPath, true)
	if err != nil {
		panic(err)
	}

}

func (p *Packer) Export(dirpath string, execPrefix string) {
	// export binary and its dependencies
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		mode := os.FileMode(0755)
		os.Mkdir(dirpath, mode)
	}

	p.ExportDyn(dirpath)
	p.ExportInterp(dirpath)
	p.ExportExe(dirpath, execPrefix)
}
