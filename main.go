package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
)

func main() {
	count := 0

	mgmtPath := "D:/Projects/Go/src/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager"

	havaLiveTestMgmts := make([]Mgmt, 0)
	noneLiveTestMgmts := make([]Mgmt, 0)

	if err := filepath.WalkDir(mgmtPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		if !strings.Contains(d.Name(), "arm") {
			return nil
		}
		count++

		// fmt.Println(path)
		// TODO

		// get version number
		// readme autorest.md module-version: 1.2.0
		moduleVersion := ""
		autorest, err := os.ReadFile(filepath.Join(path, "autorest.md"))
		if err != nil {
			return err
		}
		for _, v := range strings.Split(string(autorest), "\n") {
			if strings.Contains(v, "module-version: ") {
				moduleVersion, _ = strings.CutPrefix(v, "module-version: ")
				moduleVersion = strings.TrimSpace(moduleVersion)
				break
			}
		}
		// fmt.Println(d.Name(), moduleVersion)
		version, err := semver.NewVersion(moduleVersion)
		if err != nil {
			return err
		}

		if version.Major() > 0 {
			// 判断当前目录下是否存在live test
			havaLiveTest := false
			err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				if strings.HasSuffix(info.Name(), "_live_test.go") {
					// fmt.Println("Found live test file:", p)
					havaLiveTest = true
					return filepath.SkipDir
				}

				return nil
			})
			if err != nil {
				return err
			}

			newMgmt := Mgmt{
				armName:       d.Name(),
				armPath:       path,
				moduleVersion: moduleVersion,
				havaLiveTest:  havaLiveTest,
			}
			if havaLiveTest {
				havaLiveTestMgmts = append(havaLiveTestMgmts, newMgmt)
			} else {
				noneLiveTestMgmts = append(noneLiveTestMgmts, newMgmt)
			}

		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("mgmt count:", count)

	fmt.Println("mgmt hava live test:", len(havaLiveTestMgmts))
	for _, v := range havaLiveTestMgmts {
		if v.havaLiveTest {
			fmt.Printf("%s/v%s\n", v.armName, v.moduleVersion)
		}
	}

	fmt.Println("mgmt none live test:", len(noneLiveTestMgmts))
	for _, v := range noneLiveTestMgmts {
		if !v.havaLiveTest {
			fmt.Printf("%s/v%s\n", v.armName, v.moduleVersion)
		}
	}
}

type Mgmt struct {
	armName string

	armPath string

	moduleVersion string

	havaLiveTest bool
}
