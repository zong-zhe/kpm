// The core abstraction in kpm for working with a package or module of kcl.
package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	conf "kusionstack.io/kpm/utils"
)

func TestInitEmptyModule(t *testing.T) {
	pwd, _ := os.Getwd()
	execPath := filepath.Join(pwd, "testdata")
	expectKclModPath := filepath.Join(pwd, "testdata", "kcl.mod")

	if _, err := os.Stat(expectKclModPath); err == nil {
		err := os.Remove(expectKclModPath)
		if err != nil {
			fmt.Println("Error deleting file:", err)
		} else {
			fmt.Println("File deleted successfully")
		}
	}

	config := conf.NewEmptyConf().SetName("test name").SetVersion("test version").SetEdition("test edition").SetExecPath(execPath)
	kclPkg := NewKclPkg(config)

	err := kclPkg.InitEmptyModule()
	if err != nil {
		t.Errorf("gen kcl.mod failed")
	}

	_, err = os.Stat(expectKclModPath)
	if os.IsNotExist(err) {
		t.Errorf("gen kcl.mod failed")
	}

	err = kclPkg.InitEmptyModule()
	if err == nil {
		t.Errorf("gen kcl.mod failed")
	} else {
		if err.Error() != fmt.Sprintf("kpm:%s already exists", expectKclModPath) {
			t.Errorf("The kcl.mod already exists: '%s'", err)
		}
	}

	os.Remove(expectKclModPath)
}
