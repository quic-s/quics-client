package utils

import (
	"path"
)

func LocalAbsToRoot(rel string, root string) (string, error) {

	rootdir, _ := path.Split(root)

	result := rel[len(rootdir):]

	// dir, file := path.Split(rel)
	// result := file
	// for file != root {
	// 	dir, file = path.Split(dir[:len(dir)-1])
	// 	result = file + "/" + result
	// }
	// if file == "" {
	// 	return "", fmt.Errorf("given root folder does not exist")
	// }

	return "/" + result, nil
}
