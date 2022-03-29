package main

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

// The ansible package list should contain only 1 variable, the list of packages
func ReadAnsibleList(filePath string) ([]string, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	list := strings.Split(string(raw), "\n")
	res := make([]string, 0)
	for _, line := range list {
		// Proper yaml list item contains '- '
		if strings.Contains(line, "- ") && !strings.Contains(line, ":") {
			pkg := strings.Split(line, "- ")
			if len(pkg) >= 2 {
				res = append(res, strings.TrimSpace(pkg[1]))
			}
		}
	}
	return res, nil
}

func ComparePackages(installed []string, ansible []string, ignored []string) []string {
	sort.StringSlice.Sort(installed)
	sort.StringSlice.Sort(ansible)
	sort.StringSlice.Sort(ignored)
	res := make([]string, 0)
	for _, pkg := range installed {
		if !Contains(ansible, pkg) && !Contains(ignored, pkg) {
			res = append(res, pkg)
		}
	}

	return res
}

func ReadIndexes() []int {
	input, err := GetInput()
	if err != nil {
		fmt.Println("Error parsing input range.")
		panic(err)
	}
	input = strings.TrimSpace(input)
	indexes := make([]int, 0)
	for _, s := range strings.Split(input, ",") {
		err := ParseRange(&indexes, s)
		if err != nil {
			fmt.Println("Error parsing input range.")
			panic(err)
		}
	}
	return indexes
}

func GeneratePacmanList(pkgListPath string) []string {
	ansChan := make(chan []string)
	insChan := make(chan []string)
	ignChan := make(chan []string)

	go func() {
		pkgs, err := ReadAnsibleList(pkgListPath)
		if err != nil {
			fmt.Println("Error reading packages from yml.")
			panic(err)
		}
		ansChan <- pkgs
	}()
	go func() {
		pkgs, err := BasePackages()
		if err != nil {
			fmt.Println("Error reading packages from pacman.")
			panic(err)
		}
		insChan <- pkgs
	}()
	go func() {
		ignoredPath := IgnoredPath()
		if _, err := os.Stat(ignoredPath); !os.IsNotExist(err) {
			pkgs, err := ReadTextList(ignoredPath)
			if err != nil {
				fmt.Println("Error reading ignored package list.")
				panic(err)
			}
			ignChan <- pkgs
		} else {
			ignChan <- make([]string, 0)
		}
	}()
	ansiblePkgs := <-ansChan
	installedPkgs := <-insChan
	ignoredPkgs := <-ignChan
	return ComparePackages(installedPkgs, ansiblePkgs, ignoredPkgs)
}

func SelectPackages(pkgs []string, indexes []int) []string {
	res := make([]string, 0)
	for _, i := range indexes {
		res = append(res, pkgs[i])
	}
	return res
}

func PromptAdd(pkgs []string, target string) []string {
	PrintPackageList(pkgs)
	fmt.Printf("Select packages to add to %s. (eg. 1,2,5-7)\n", target)
	indxs := ReadIndexes()
	toAdd := SelectPackages(pkgs, indxs)
	fmt.Println("===================================================")
	fmt.Printf("The following packages will be added to %s:\n", target)
	fmt.Println("===================================================")
	PrintStringSlice(toAdd)
	fmt.Println("===================================================")
	fmt.Println("Continue? (y,N)")
	input, err := GetInput()
	if err != nil {
		fmt.Println("Error reading input.")
		panic(err)
	}
	if strings.TrimSpace(strings.ToLower(input)) == "y" {
		return toAdd
	} else {
		return nil
	}
}

// Add packages to the input yml file
func AddPackagesYml(pkgs []string, ymlPath string) {

	// open output file
	fo, err := os.OpenFile(ymlPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	for _, pkg := range pkgs {
		pkg := fmt.Sprintf("  - %s\n", pkg)
		fo.WriteString(pkg)
	}

}

func SyncPacmanPackages(pkgListPath string) {
	pkgs := GeneratePacmanList(pkgListPath)
	pkgsToAdd := PromptAdd(pkgs, "Ansible List")
	if pkgsToAdd != nil {
		AddPackagesYml(pkgsToAdd, pkgListPath)
	}
}

func main() {
	home, _ := os.UserHomeDir()
	pkgListPath := path.Join(home, "dev", "ansible-setup", "arch", "packages.yml")
	SyncPacmanPackages(pkgListPath)
	CheckIgnored(pkgListPath)
}
