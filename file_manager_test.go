package blimpy

import (
	"testing"
	"os"
	"io/ioutil"
	"path/filepath"
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
