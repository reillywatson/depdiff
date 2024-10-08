package depdiff

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func DepDiff(moduleName, pkgPath, oldCommit, newCommit string) ([]string, error) {
	filesBytes, err := exec.Command("git", "diff", "--name-only", oldCommit, newCommit).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Exec error: %s. Output: %s", err, string(filesBytes))
	}
	depsBytes, err := exec.Command("go", "list", "-f", `{{ .Deps }}`, pkgPath).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Exec error: %s. Output: %s", err, string(filesBytes))
	}
	depsStr := string(depsBytes)
	deps := strings.Split(depsStr[1:len(depsStr)-1], " ")
	deps = append(deps, pkgPath)
	depMap := map[string]bool{}
	for _, dep := range deps {
		depMap[dep] = true
	}
	files := strings.Split(string(filesBytes), "\n")
	changedPackages := map[string]bool{}
	rootPath, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Exec error: %s. Output: %s", err, string(rootPath))
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Getwd error: %s", wd)
	}
	wd = strings.TrimSpace(wd)
	prefix := strings.TrimPrefix(wd, strings.TrimSpace(string(rootPath)))
	prefix = strings.TrimPrefix(prefix, "/")
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		f = moduleName + strings.TrimPrefix(f, prefix)
		if !strings.Contains(f, "/") {
			continue
		}
		if strings.HasSuffix(f, ".go") {
			f = f[:strings.LastIndex(f, "/")]
			if depMap[f] {
				changedPackages[f] = true
			}
		} else {
			// these may be embedded files, look for any parent directories that are dependencies
			for k := range depMap {
				if strings.HasPrefix(f, k) {
					changedPackages[k] = true
				}
			}
		}
	}
	var changedPkgList []string
	for pkg := range changedPackages {
		changedPkgList = append(changedPkgList, pkg)
	}
	sort.Strings(changedPkgList)
	return changedPkgList, nil
}
