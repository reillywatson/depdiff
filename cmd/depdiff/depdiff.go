package main

import (
	"fmt"
	"github.com/reillywatson/depdiff"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: shoulddeploy pkgPath oldCommit newCommit")
		os.Exit(1)
	}
	pkgPath := os.Args[1]
	oldCommit := os.Args[2]
	newCommit := os.Args[3]
	pkgs, err := depdiff.DepDiff(pkgPath, oldCommit, newCommit)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, pkg := range pkgs {
		fmt.Println(pkg)
	}
}
