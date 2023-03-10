package conf

import "testing"

func TestSetName(t *testing.T) {
	conf := NewEmptyConf()
	conf = conf.SetName("test_name")
	if conf.Name != "test_name" {
		t.Errorf("conf.SetName failed %s", conf.Name)
	}
}

func TestSetEdition(t *testing.T) {
	conf := NewEmptyConf()
	conf = conf.SetEdition("test_edition")
	if conf.Edition != "test_edition" {
		t.Errorf("conf.SetEdition failed %s", conf.Edition)
	}
}

func TestSetVersion(t *testing.T) {
	conf := NewEmptyConf()
	conf = conf.SetVersion("test_version")
	if conf.Version != "test_version" {
		t.Errorf("conf.SetVersion failed %s", conf.Version)
	}
}

func TestSetExecPath(t *testing.T) {
	conf := NewEmptyConf()
	conf = conf.SetExecPath("test_exec_path")
	if conf.ExecPath != "test_exec_path" {
		t.Errorf("conf.SetExecPath failed %s", conf.ExecPath)
	}
}
