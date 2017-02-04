package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/linxiulei/sucker/analyzer"
	"github.com/linxiulei/sucker/builder"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [OPTIONS] [cmd1, cmd2...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	var exportDir = flag.String("export-dir", "export_dir", "dir for export")
	var execPrefix = flag.String("exec-prefix", "", "dir of running cmd")
	flag.Parse()

	if *execPrefix == "" {
		execPrefix = exportDir
	}

	a := new(analyzer.Analyzer)
	packer := new(builder.Packer)

	for _, i := range flag.Args() {
		elfObject := a.AnaFromFile(i)
		packer.ElfObject = elfObject
		packer.Export(*exportDir, *execPrefix)
	}
}
