package client

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/mattn/go-zglob"
	"github.com/otiai10/copy"
	"gotest.tools/v3/assert"
	"kcl-lang.io/kpm/pkg/downloader"
	"kcl-lang.io/kpm/pkg/features"
	"kcl-lang.io/kpm/pkg/utils"
)

func testRunWithModSpecVersion(t *testing.T, kpmcli *KpmClient) {
	pkgPath := getTestDir("test_run_with_modspec_version")
	modbkPath := filepath.Join(pkgPath, "kcl.mod.bk")
	modPath := filepath.Join(pkgPath, "kcl.mod")
	modExpect := filepath.Join(pkgPath, "kcl.mod.expect")
	lockbkPath := filepath.Join(pkgPath, "kcl.mod.lock.bk")
	lockPath := filepath.Join(pkgPath, "kcl.mod.lock")
	lockExpect := filepath.Join(pkgPath, "kcl.mod.lock.expect")
	err := copy.Copy(modbkPath, modPath)
	if err != nil {
		t.Fatal(err)
	}

	err = copy.Copy(lockbkPath, lockPath)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		// remove the copied files
		err := os.RemoveAll(modPath)
		if err != nil {
			t.Fatal(err)
		}
		err = os.RemoveAll(lockPath)
		if err != nil {
			t.Fatal(err)
		}
	}()

	res, err := kpmcli.Run(
		WithRunSource(
			&downloader.Source{
				Local: &downloader.Local{
					Path: pkgPath,
				},
			},
		),
	)

	if err != nil {
		t.Errorf("Failed to run package: %v", err)
	}

	assert.Equal(t, res.GetRawYamlResult(), "res: Hello World!")
	expectedMod, err := os.ReadFile(modExpect)
	if err != nil {
		t.Fatal(err)
	}
	gotMod, err := os.ReadFile(modPath)
	if err != nil {
		t.Fatal(err)
	}

	expectedLock, err := os.ReadFile(lockExpect)
	if err != nil {
		t.Fatal(err)
	}

	gotLock, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, utils.RmNewline(string(expectedMod)), utils.RmNewline(string(gotMod)))
	assert.Equal(t, utils.RmNewline(string(expectedLock)), utils.RmNewline(string(gotLock)))
}

func TestRun(t *testing.T) {
	features.Enable(features.SupportNewStorage)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunWithOciDownloader", testRunWithOciDownloader)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunDefaultRegistryDep", testRunDefaultRegistryDep)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunInVendor", testRunInVendor)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunRemoteWithArgsInvalid", testRunRemoteWithArgsInvalid)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunRemoteWithArgs", testRunRemoteWithArgs)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunWithNoSumCheck", testRunWithNoSumCheck)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunWithGitPackage", testRunWithGitPackage)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunGit", testRunGit)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunOciWithSettingsFile", testRunOciWithSettingsFile)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunWithModSpecVersion", testRunWithModSpecVersion)

	features.Disable(features.SupportNewStorage)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunWithOciDownloader", testRunWithOciDownloader)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunDefaultRegistryDep", testRunDefaultRegistryDep)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunInVendor", testRunInVendor)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunRemoteWithArgsInvalid", testRunRemoteWithArgsInvalid)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunRemoteWithArgs", testRunRemoteWithArgs)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunWithNoSumCheck", testRunWithNoSumCheck)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunWithGitPackage", testRunWithGitPackage)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunGit", testRunGit)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunOciWithSettingsFile", testRunOciWithSettingsFile)
	RunTestWithGlobalLockAndKpmCli(t, "TestRunWithModSpecVersion", testRunWithModSpecVersion)
}
func TestRunWithHyphenEntries(t *testing.T) {
	testFunc := func(t *testing.T, kpmcli *KpmClient) {
		pkgPath := getTestDir("test_run_hyphen_entries")

		res, err := kpmcli.Run(
			WithRunSource(
				&downloader.Source{
					Local: &downloader.Local{
						Path: pkgPath,
					},
				},
			),
		)

		if err != nil {
			t.Fatal(err)
		}

		expect, err := os.ReadFile(filepath.Join(pkgPath, "stdout"))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, utils.RmNewline(res.GetRawYamlResult()), utils.RmNewline(string(expect)))
	}

	RunTestWithGlobalLockAndKpmCli(t, "testRunWithHyphenEntries", testFunc)
}

func TestFuckingWindowsPath(t *testing.T) {
	pwd, _ := os.Getwd()
	fmt.Printf("pwd: %v\n", pwd)

	absPath, _ := filepath.Abs(".")
	fmt.Printf("absPath: %v\n", absPath)
	fmt.Printf("%v == %v: %v\n", absPath, pwd, absPath == pwd)
	path1 := "D:\\a\\kpm\\kpm\\pkg\\client\\test_data\\test_run_hyphen_entries"
	pathUrl, _ := url.Parse(path1)
	fmt.Printf("pathUrl: %v\n", pathUrl.String())
	path2 := "d:\\a\\kpm\\kpm\\pkg\\client\\test_data\\test_run_hyphen_entries"

	match, _ := zglob.Match(path1, path2)

	fmt.Println("Paths are equal:", match)
}
