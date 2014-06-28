package blimpy

import (
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DataDir     string `json:"DataDir"`
	ApiPort     int    `json:"apiPort"`
}

type Blimpy struct {
	metadataDb *gorp.DbMap
}

func New(config *Config) (*Blimpy, error) {
	return nil, nil
}
