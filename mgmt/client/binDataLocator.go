package client

type BinDataLocator struct {
}

func (l *BinDataLocator) Get(path string) ([]byte, error) {
	// Trim leading /
	if len(path) > 0 && path[0] == '/' {
		path = path[1:len(path)]
	}

	data, err := Asset(path)
	if err == nil {
		return data, nil
	}

	return l.serveIndex()
}

func (l *BinDataLocator) serveIndex() ([]byte, error) {
	index, err := Asset("index.html")
	if err != nil {
		return nil, err
	}

	return index, nil
}
