// The core abstraction in kpm for working with a package or module of kcl.
package mod

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"kusionstack.io/kpm/gen/pkg"
	conf "kusionstack.io/kpm/utils"
)

const kclMod = "kcl.mod"

type KclPkg struct {
	HomePath string       `toml:"-"`
	Pkg      *pkg.Package `toml:"package"`
}

func NewKclPkg(conf conf.Config) KclPkg {
	return KclPkg{
		HomePath: conf.ExecPath,
		Pkg: &pkg.Package{
			Name:    conf.Name,
			Version: conf.Version,
			Edition: conf.Edition,
		},
	}
}

func (kclPkg KclPkg) InitEmptyModule() error {
	kclModPath := filepath.Join(kclPkg.HomePath, kclMod)
	_, err := os.Stat(kclModPath)
	if os.IsNotExist(err) {
		var buf bytes.Buffer
		err := toml.NewEncoder(&buf).Encode(kclPkg)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(kclModPath, buf.Bytes(), 0644)
		fmt.Printf("buf.String(): %v, %s\n", buf.String(), kclPkg.HomePath)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("kpm:%s already exists", kclModPath)
	}
}
