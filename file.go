package blimpy

import (
	"os"
)

type File struct {
	path        string   `db:"-"`
	file        *os.File `db:"-"`
	Id          string   `db:"id"`
	Name        string   `db:"name"`
	Description string   `db:"description"`
	MimeType    string   `db:"mime_type"`
}

func (self *File) Open() error {
	var err error

	if self.file, err = os.OpenFile(self.path, os.O_RDWR|os.O_CREATE, os.FileMode(0644)); err != nil {
		return err
	} else {
		return nil
	}
}

func (self *File) Close() {
	self.file.Close()
	self.file = nil
}

func (self *File) IsOpen() bool {
	return self.file != nil
}

func (self *File) File() *os.File {
	return self.file
}

func (self *File) Read(p []byte) (int, error) {
	return self.file.Read(p)
}

func (self *File) Write(p []byte) (int, error) {
	return self.file.Write(p)
}

func (self *File) Seek(offset int64, whence int) (int64, error) {
	return self.file.Seek(offset, whence)
}
