package main

import (
	"fmt"
	"os"
	"path"
	"sort"
)

func IgnoredPath() string {
	home, _ := os.UserHomeDir()
	return path.Join(home, ".config", "package_syncer", "ignored.txt")
}

func CheckIgnored(pkgListPath string) {
	pkgs := GeneratePacmanList(pkgListPath)
	toAdd := PromptAdd(pkgs, "Ignored List")
	AddIgnored(IgnoredPath(), toAdd)
}

func AddIgnored(filePath string, pkgs []string) {
	if _, err := os.Stat(filePath); os.IsExist(err) {
		current, err := ReadTextList(filePath)
		if err != nil {
			fmt.Println("Error reading ignored files.")
			panic(err)
		}
		pkgs = append(pkgs, current...)
	}
	sort.StringSlice.Sort(pkgs)
	WriteTextList(filePath, pkgs)
}
