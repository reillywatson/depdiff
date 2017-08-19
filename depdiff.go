package depdiff

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

func DepDiff(pkgPath, oldCommit, newCommit string) ([]string, error) {
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
	for _, f := range files {
		if !strings.HasPrefix(f, "src/") {
			continue
		}
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		f = strings.TrimPrefix(f, "src/")
		if !strings.Contains(f, "/") {
			continue
		}
		f = f[:strings.LastIndex(f, "/")]
		if depMap[f] {
			changedPackages[f] = true
		}
	}
	var changedPkgList []string
	for pkg := range changedPackages {
		changedPkgList = append(changedPkgList, pkg)
	}
	sort.Strings(changedPkgList)
	return changedPkgList, nil
}
