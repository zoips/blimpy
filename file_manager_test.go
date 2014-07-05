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
