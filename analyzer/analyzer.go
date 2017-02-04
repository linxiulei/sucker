package analyzer

import (
	"github.com/linxiulei/sucker/elfpatch"
)

type Stat struct {
	DynFiles []string
	Execute  *elfpatch.ElfObject
}

type Analyzer struct {
	Loader string
}

func (a *Analyzer) AnaFromFile(filepath string) *elfpatch.ElfObject {
	elfObject := new(elfpatch.ElfObject)
	elfObject.FilePath = filepath
	elfObject.Scan()

	return elfObject
}
