package client

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

type FsLocator struct {
	Root string
}

func (l *FsLocator) Get(filePath string) ([]byte, error) {
	// Trim leading /
	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:len(filePath)]
	}

	if data, err := l.readFile(filePath); err == nil {
		return data, nil
	}

	return l.serveIndex()
}

func (l *FsLocator) serveIndex() ([]byte, error) {
	data, err := l.readFile("index.html")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (l *FsLocator) readFile(filePath string) ([]byte, error) {
	fullPath := path.Join(l.Root, filePath)
	if _, err := os.Stat(fullPath); err != nil {
		log.Println("ERROR", err, fullPath)
		return nil, err
	}

	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
