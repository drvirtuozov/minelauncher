package main

import (
	"io/ioutil"
	"path"
)

func getLibsPaths(dir string) (paths []string) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filepath := path.Join(dir, file.Name())

		if file.IsDir() {
			paths = append(paths, getLibsPaths(filepath)...)
		} else {
			paths = append(paths, filepath)
		}
	}

	return paths
}
