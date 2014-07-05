package blimpy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFSFileManagerEnsureDirectory(t *testing.T) {
	fm := FSFileManager{}
	root, _ := ioutil.TempDir("", "fm-test")

	defer os.RemoveAll(root)

	fm.root = root

	err := fm.ensureStorageDirectories()
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range dirNames {
		if _, err := os.Stat(filepath.Join(root, string(s))); err != nil {
			t.Fatal(err)
		}

		for _, s2 := range dirNames {
			if _, err = os.Stat(filepath.Join(root, string(s), string(s2))); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestFSFileManagerSetRoot(t *testing.T) {
	fm := FSFileManager{}
	root, _ := ioutil.TempDir("", "fm-test")

	defer os.RemoveAll(root)

	err := fm.SetRoot(root)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(root, fileDatabaseName)); err != nil {
		t.Fatal(err)
	}

	if fm.dbHandle == nil || fm.dbMap == nil {
		t.Fatal("Expected database handles")
	}

	tables := []string{
		"files",
	}

	stmt, err := fm.dbHandle.Prepare("select name from sqlite_master where type = 'table' and name = ?")
	if err != nil {
		t.Fatal(err)
	}

	for _, table := range tables {
		rows, err := stmt.Query(table)
		if err != nil {
			t.Fatal(err)
		}

		defer rows.Close()

		if !rows.Next() {
			t.Fatalf("Expected table %s", table)
		}
	}
}

func TestFSFileManagerInsertFile(t *testing.T) {
	root, _ := ioutil.TempDir("", "fm-test")
	fm, _ := NewFSFileManager(root)

	defer os.RemoveAll(root)

	for i := 0; i < 10; i++ {
		fd, _ := ioutil.TempFile("", "fm-test-file")
		data := fmt.Sprintf("test file %d", i)

		fd.Write([]byte(data))
		fd.Seek(0, 0)

		file := File{
			Name:        "foo",
			Description: "bar",
			MimeType:    "test/foo",
		}

		err := fm.InsertFile(&file, fd)
		if err != nil {
			t.Fatal(err)
		}

		if file.Id == "" {
			t.Fatal("Expected an id")
		}

		_, err = os.Stat(filepath.Join(root, file.Id[0:1], file.Id[1:2], file.Id))
		if err != nil {
			t.Fatal(err)
		}

		err = file.Open()
		if err != nil {
			t.Fatal(err)
		}

		data2, err := ioutil.ReadAll(&file)
		if err != nil {
			t.Fatal(err)
		}

		if string(data2) != data {
			t.Fatalf("Expected %s == %s", data, data2)
		}
	}
}

func TestFSFileManagerGetFile(t *testing.T) {
	root, _ := ioutil.TempDir("", "fm-test")
	fm, _ := NewFSFileManager(root)

	defer os.RemoveAll(root)

	for i := 0; i < 10; i++ {
		fd, _ := ioutil.TempFile("", "fm-test-file")
		data := fmt.Sprintf("test file %d", i)

		fd.Write([]byte(data))
		fd.Seek(0, 0)

		file := File{
			Name:        "foo",
			Description: "bar",
			MimeType:    "test/foo",
		}

		err := fm.InsertFile(&file, fd)
		if err != nil {
			t.Fatal(err)
		}

		file2, err := fm.GetFile(file.Id)
		if err != nil {
			t.Fatal(err)
		}

		if file2 == nil {
			t.Fatal("Expected file")
		}

		err = file2.Open()
		if err != nil {
			t.Fatal(err)
		}

		data2, err := ioutil.ReadAll(file2)
		if err != nil {
			t.Fatal(err)
		}

		if string(data2) != data {
			t.Fatalf("Expected %s == %s", data, data2)
		}
	}
}
