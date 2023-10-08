package utils

import (
	"log"
	"path"
	"path/filepath"
)

// LocalAbsToRoot converts an absolute path to a relative path to the root
// ex) LocalAbsToRoot("/a/rootDir/c", "/a/rootDir") -> "/rootDir/c"
func LocalAbsToRoot(abs string, root string) string {

	rootdir, _ := path.Split(root) //rootdir == /a
	result := abs[len(rootdir):]   //result == /roorDir/c
	return result
}

// LocalRelToRoot converts a relative path to a relative path to the root
// ex) LocalRelToRoot("./b/c", "a/b") -> "/c"
func LocalRelToRoot(rel string, root string) string {

	abs, err := filepath.Abs(rel)
	if err != nil {
		log.Panic(err)

	}
	return LocalAbsToRoot(abs, root)

}
