package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: shoulddeploy pkgPath oldCommit newCommit")
		os.Exit(1)
	}
	pkgPath := os.Args[1]
	oldCommit := os.Args[2]
	newCommit := os.Args[3]
	filesBytes, _ := exec.Command("git", "diff", "--name-only", oldCommit, newCommit).CombinedOutput()
	depsBytes, _ := exec.Command("go", "list", "-f", `{{ .Deps }}`, pkgPath).CombinedOutput()
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
	if len(changedPackages) == 0 {
		return
	}
	changedPkgList := []string{}
	for pkg := range changedPackages {
		changedPkgList = append(changedPkgList, pkg)
	}
	sort.Strings(changedPkgList)
	fmt.Println(changedPkgList)
}
