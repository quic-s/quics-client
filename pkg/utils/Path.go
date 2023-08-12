package utils

import (
	"path"
	"path/filepath"
)

func LocalAbsToRoot(abs string, root string) string {

	rootdir, _ := path.Split(root)

	result := abs[len(rootdir):]

	// dir, file := path.Split(rel)
	// result := file
	// for file != root {
	// 	dir, file = path.Split(dir[:len(dir)-1])
	// 	result = file + "/" + result
	// }
	// if file == "" {
	// 	return "", fmt.Errorf("given root folder does not exist")
	// }

	return "/" + result
}

func LocalRelToRoot(rel string, root string) (string, error) {

	abs, err := filepath.Abs(rel)
	if err != nil {
		return "", err

	}
	return LocalAbsToRoot(abs, root), nil

}
