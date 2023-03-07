package pool

import (
	"encoding/json"
	"os"
	"time"
)

type AddrModel struct {
	ID       int64
	Addr     string
	Header   json.RawMessage
	Speed    time.Duration
	Health   int64
	CreateAt time.Time
	CheckAt  time.Time
}

var DefaultFile = &FileStorager{
	FileName: "proxies.json",
}

type FileStorager struct {
	FileName string
}

func (s *FileStorager) Save(rows []*AddrModel) error {
	raw, err := json.Marshal(rows)
	if err != nil {
		return err
	}
	return os.WriteFile(s.FileName, raw, os.ModePerm)
}

func (s *FileStorager) Read() ([]*AddrModel, error) {
	raw, err := os.ReadFile(s.FileName)
	if err != nil {
		return nil, err
	}
	rows := []*AddrModel{}
	err = json.Unmarshal(raw, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
