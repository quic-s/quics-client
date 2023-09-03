package utils

import (
	"log"
	"path"
	"path/filepath"
)

// LocalAbsToRoot converts an absolute path to a relative path to the root
// ex) LocalAbsToRoot("/a/b/c", "a/b") -> "/c"
func LocalAbsToRoot(abs string, root string) string {

	rootdir, _ := path.Split(root)
	result := abs[len(rootdir):]
	return "/" + result
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

// ex) SplitPathWithRoot("/a/b/c", "a/b") -> "/a/b", "/c"
func SplitPathWithRoot(path string, root string) (string, string) {

	rootdir, _ := filepath.Split(root)
	afterPath := "/" + path[len(rootdir):]
	beforePath := path[:len(rootdir)]
	return beforePath, afterPath
}
