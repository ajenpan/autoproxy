package pool

import (
	"encoding/json"
	"os"

	"autoproxy/util"
)

type FileStorager struct {
	filename string
	*MemStorager
}

func NewFileStorager(filename string) *FileStorager {
	fs := &FileStorager{
		filename:    filename,
		MemStorager: NewMemStorager(),
	}

	exist, _ := util.FileExist(filename)

	if !exist {
		return fs
	}

	all, err := fs.readfile()
	if err != nil {
		return nil
	}

	for _, v := range all {
		fs.Save(v)
	}
	return fs
}

func (s *FileStorager) Flash() error {
	s.MemStorager.Flash()
	all, _ := s.MemStorager.GetAll()
	return s.writefile(all)
}

func (s *FileStorager) readfile() ([]*ProxyAddr, error) {
	raw, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}
	rows := []*ProxyAddr{}
	err = json.Unmarshal(raw, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *FileStorager) writefile(rows []*ProxyAddr) error {
	// raw, err := json.Marshal(rows)
	raw, err := json.MarshalIndent(rows, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, raw, 0644)
}
