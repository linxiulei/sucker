package elfpatch

import (
	"fmt"
	"testing"
)

func TestElfObjectScan(t *testing.T) {
	elfObject := new(ElfObject)
	elfObject.FilePath = "/bin/ls"
	elfObject.Scan()

	fmt.Println(elfObject.DynPathMap)
	v := elfObject.DynPathMap["ld-linux-x86-64.so.2"]

	if v != "/lib64/ld-linux-x86-64.so.2" {
		t.Error(
			"expected", "/lib64/ld-linux-x86-64.so.2",
			"got", v,
		)
	}

	if elfObject.InterpPath != "/lib64/ld-linux-x86-64.so.2" {
		t.Error(
			"expected", "/lib64/ld-linux-x86-64.so.2",
			"got", v,
		)
	}

	dynlist := []string{
		"libselinux.so.1",
		"librt.so.1",
		"libcap.so.2",
		"libacl.so.1",
		"libc.so.6",
	}

	for _, i := range dynlist {
		fmt.Println(i)
	}
}
