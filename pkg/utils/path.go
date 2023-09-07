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

// ex) SplitBeforeAfterRoot("/a/rootDir/c") -> ("/a", "/rootDir/c")
func SplitBeforeAfterRoot(path string) (string, string) {

	dirlist := GetRootDirs()

	for _, v := range dirlist { // v == /rootDir/c
		if path[:len(v)] == v {
			dir, _ := filepath.Split(v)
			log.Println("before, after :", dir[len(dir)-1:], " , ", path[len(dir)-1:])
			return dir[:len(dir)-1], path[len(dir)-1:]
		}

	}
	return "", ""
}
