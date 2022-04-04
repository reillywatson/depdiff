package main

import (
	"fmt"
	"os"

	"github.com/reillywatson/depdiff"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: depdiff moduleName pkgPath oldCommit newCommit")
		os.Exit(1)
	}
	moduleName := os.Args[1]
	pkgPath := os.Args[2]
	oldCommit := os.Args[3]
	newCommit := os.Args[4]
	pkgs, err := depdiff.DepDiff(moduleName, pkgPath, oldCommit, newCommit)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, pkg := range pkgs {
		fmt.Println(pkg)
	}
}
