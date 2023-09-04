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

// ex) SplitPathWithRoot("/a/b/c") -> ("/a/b", "/c")
func SplitBeforeAfterRoot(path string) (string, string) {

	dirlist := GetRootDirs()

	for _, v := range dirlist {
		if path[:len(v)] == v {
			return path[:len(v)], path[len(v):]
		}

	}
	return "", ""
}
