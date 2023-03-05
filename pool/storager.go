package pool

import (
	"encoding/json"
	"os"
	"time"
)

type Model struct {
	ID int64
	*ProxyAddrItem
	Speed      int64
	Health     int64
	CreateTime time.Time
	UpdateTime time.Time
}

type FileStorager struct {
	FileName string
}

func (s *FileStorager) Save(rows []*Model) error {
	raw, err := json.Marshal(rows)
	if err != nil {
		return err
	}
	return os.WriteFile(s.FileName, raw, os.ModePerm)
}

func (s *FileStorager) Read() ([]*Model, error) {
	raw, err := os.ReadFile(s.FileName)
	if err != nil {
		return nil, err
	}
	rows := []*Model{}
	err = json.Unmarshal(raw, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
