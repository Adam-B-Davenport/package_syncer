package main

import (
	"fmt"
	"os"
	"strings"
)

func PrintSlice(slice []string) {
	for _, str := range slice {
		fmt.Println(str)
	}
}

func PrintPackageList(pkgs []string) {
	for i, pkg := range pkgs {
		fmt.Printf("[%d] %s\n", i, pkg)
	}
}

func Contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func ReadTextList(filePath string) ([]string, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	text := strings.TrimSpace(string(raw))
	return strings.Split(text, "\n"), nil
}