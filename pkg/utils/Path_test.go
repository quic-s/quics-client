package utils_test

import (
	//"path"
	"testing"

	"github.com/quic-s/quics-client/pkg/utils"
)

func TestRtoA(t *testing.T) {
	abs, err := utils.LocalAbsToRoot("/home/bob/work/utils/src/github.com", "/home/bob")
	t.Error(abs)
	t.Error(err)

	// rel := "/home/bob/work/utils/src/github.com"
	// root := "bob"

	// dir, file := path.Split(rel)
	// result := file
	// t.Error(result, "")
	// for file != root {
	// 	dir, file = path.Split(dir[:len(dir)-1])
	// 	result = file + "/" + result
	// 	t.Log(result)
	// }

}
