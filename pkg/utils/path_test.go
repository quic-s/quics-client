package utils_test

import (
	//"path"
	"testing"

	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestLocalAbsToRoot(t *testing.T) {
	abs := utils.LocalAbsToRoot("/home/bob/work/utils/src/github.com", "/home/bob")
	assert.Equal(t, "/bob/work/utils/src/github.com", abs)
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

func TestLocalRelToRoot(t *testing.T) {
	abs := utils.LocalRelToRoot("../qic", "/workspaces/quics-client")
	assert.Equal(t, "/quics-client/pkg/qic", abs)
}
