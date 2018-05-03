package util

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/drvirtuozov/minelauncher/config"
	"github.com/drvirtuozov/minelauncher/events"
)

type passThruReader struct {
	io.Reader
	total  int64
	length int64
	text   string
}

func (pt *passThruReader) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)

	if n > 0 {
		prevFraction := float64(pt.total) / float64(pt.length)
		prevPercentage := int(prevFraction * 100)
		pt.total += int64(n)
		fraction := float64(pt.total) / float64(pt.length)
		percentage := int(fraction * 100)

		if percentage > prevPercentage {
			events.TaskProgress <- events.ProgressBarFraction{
				Fraction: fraction,
				Text:     pt.text + " " + strconv.Itoa(percentage) + "%",
			}
		}
	}

	return n, err
}

func GetLibsPaths(dir string) (paths []string, err error) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filepath := path.Join(dir, file.Name())

		if file.IsDir() {
			filepaths, err := GetLibsPaths(filepath)

			if err != nil {
				return nil, err
			}

			paths = append(paths, filepaths...)
		} else {
			paths = append(paths, filepath)
		}
	}

	return paths, nil
}

func DownloadZip(url string) (filePath string, err error) {
	res, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	cfg, err := config.Get()

	if err != nil {
		return "", err
	}

	file, err := ioutil.TempFile("", cfg.Launcher+"-update")

	if err != nil {
		return "", err
	}

	defer file.Close()
	reader := &passThruReader{
		Reader: res.Body,
		length: res.ContentLength,
		text:   "Downloading client update...",
	}
	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return "", err
	}

	_, err = file.Write(data)

	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func Unzip(zipPath, destPath string) error {
	zipReader, err := zip.OpenReader(zipPath)

	if err != nil {
		return err
	}

	if err := os.MkdirAll(destPath, 0777); err != nil {
		return err
	}

	for i, file := range zipReader.File {
		filePath := path.Join(destPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}

		events.TaskProgress <- events.ProgressBarFraction{
			Fraction: float64(i+1) / float64(len(zipReader.File)),
			Text:     fmt.Sprintf("Extracting files... %d of %d", i+1, len(zipReader.File)),
		}
		fileReader, err := file.Open()

		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())

		if err != nil {
			return err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}

		fileReader.Close()
		targetFile.Close()
	}

	return nil
}

func CopyDir(dirPath, destDir string) error {
	files, err := ioutil.ReadDir(dirPath)

	if err != nil {
		return err
	}

	for _, file := range files {
		fromPath := path.Join(dirPath, file.Name())
		toPath := path.Join(destDir, file.Name())

		if file.IsDir() {
			os.MkdirAll(toPath, file.Mode())
			CopyDir(fromPath, toPath)
			continue
		}

		bytes, err := ioutil.ReadFile(fromPath)

		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(toPath, bytes, file.Mode()); err != nil {
			return err
		}
	}

	return nil
}

func GetCommitFromFilename(filename string) string {
	return strings.TrimSuffix(filename[strings.LastIndex(filename, "-")+1:], path.Ext(filename))
}
