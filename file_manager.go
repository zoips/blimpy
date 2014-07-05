package blimpy

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
	"path/filepath"
)

const (
	dirNames = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type FileManager interface {
	GetRoot() string
	SetRoot(root string) error
	GetFile(id string) *File
	DeleteFile(id string) error
	InsertFile(file *File) error
	UpdateFile(file *File) error
}

type FSFileManager struct {
	root     string
	dbHandle *sql.DB
	dbMap    *gorp.DbMap
}

func NewFSFileManager(root string) (*FSFileManager, error) {
	self := FSFileManager{}

	err := self.SetRoot(root)
	if err != nil {
		return nil, err
	}

	return &self, nil
}

func (self *FSFileManager) GetRoot() string {
	return self.root
}

func (self *FSFileManager) SetRoot(root string) error {
	if self.dbHandle != nil {
		if err := self.dbHandle.Close(); err != nil {
			return err
		}
	}

	db, err := sql.Open("sqlite3", filepath.Join(root, "blimpy_files.sqlite3"))
	if err != nil {
		return err
	}
	self.root = root
	self.dbHandle = db
	self.dbMap = &gorp.DbMap{
		Db:      self.dbHandle,
		Dialect: gorp.SqliteDialect{},
	}

	self.dbMap.AddTableWithName(File{}, "files").SetKeys(false, "Id")

	err = self.dbMap.CreateTablesIfNotExists()
	if err != nil {
		return err
	}

	err = self.ensureStorageDirectories()
	if err != nil {
		return err
	}

	return nil
}

func (self *FSFileManager) ensureStorageDirectories() error {
	for _, s := range dirNames {
		err := os.Mkdir(filepath.Join(self.root, string(s)), os.FileMode(0755))
		if err != nil && !os.IsExist(err) {
			return err
		}

		for _, s2 := range dirNames {
			err := os.Mkdir(filepath.Join(self.root, string(s), string(s2)), os.FileMode(0755))
			if err != nil && !os.IsExist(err) {
				return err
			}
		}
	}

	return nil
}

func (self *FSFileManager) InsertFile(file *File, fd *os.File) error {
	hash := sha512.New()

	_, err := io.Copy(hash, fd)
	if err != nil {
		return err
	}

	_, err = fd.Seek(0, 0)
	if err != nil {
		return err
	}

	file.Id = hex.EncodeToString(hash.Sum(nil))
	file.path = filepath.Join(self.root, file.Id[0:1], file.Id[1:2])

	err = file.Open()
	if err != nil {
		return err
	}

	io.Copy(file, fd)

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	return self.dbMap.Insert(file)
}

func (self *FSFileManager) GetFile(id string) *File {
	return nil
}

func (self *FSFileManager) UpdateFile(file *File) error {
	return nil
}

func (self *FSFileManager) DeleteFile(id string) error {
	return nil
}
