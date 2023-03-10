package ops

import (
	mod "kusionstack.io/kpm/core"
	conf "kusionstack.io/kpm/utils"
)

func KpmInit(conf conf.Config) error {
	return mod.NewKclPkg(conf).InitEmptyModule()
}
