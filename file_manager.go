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
	dirNames         = "abcdefghijklmnopqrstuvwxyz0123456789"
	fileDatabaseName = "blimpy_files.sqlite3"
)

func makeFilePath(root, id string) string {
	return filepath.Join(root, id[0:1], id[1:2], id)
}

type FileManager interface {
	GetRoot() string
	SetRoot(root string) error
	GetFile(id string) (*File, error)
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

	db, err := sql.Open("sqlite3", filepath.Join(root, fileDatabaseName))
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
	file.path = makeFilePath(self.root, file.Id)

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

func (self *FSFileManager) GetFile(id string) (*File, error) {
	var file File

	err := self.dbMap.SelectOne(&file, "select * from files where id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	file.path = makeFilePath(self.root, id)

	return &file, nil
}

func (self *FSFileManager) UpdateFile(file *File) error {
	return nil
}

func (self *FSFileManager) DeleteFile(id string) error {
	return nil
}
