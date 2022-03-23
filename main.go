package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
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

// Add to slice, inclusive range
func AppendRange(indxs *[]int, start int, end int) {
	for i := start; i <= end; i++ {
		*indxs = append(*indxs, i)
	}
}

func ParseRange(indxs *[]int, s string) error {
	values := strings.Split(s, "-")
	switch len(values) {
	case 1:
		str := strings.TrimSpace(values[0])
		if val, err := strconv.Atoi(str); err != nil {
			return err
		} else {
			*indxs = append(*indxs, val)
			return nil
		}
	case 2:
		s1 := strings.TrimSpace(values[0])
		s2 := strings.TrimSpace(values[1])
		v1, e1 := strconv.Atoi(s1)
		v2, e2 := strconv.Atoi(s2)
		if e1 != nil && e2 != nil && v1 <= v2 {
			return errors.New("input: unable to parse range")
		} else {
			AppendRange(indxs, v1, v2)
			return nil
		}
	}
	return nil
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

	cwd, _ := os.Getwd()
	ignoredPath := path.Join(cwd, "ignored.txt")

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
		if _, err := os.Stat(ignoredPath); os.IsExist(err) {
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

func SyncPacmanPackages() {
	home, _ := os.UserHomeDir()
	pkgListPath := path.Join(home, "dev", "ansible-setup", "arch", "packages.yml")

	pkgs := GeneratePacmanList(pkgListPath)
	PrintPackageList(pkgs)
	fmt.Println("Select packages to add to package list. (eg. 1,2,5-7)")
	indxs := ReadIndexes()
	pkgs = SelectPackages(pkgs, indxs)
	fmt.Println("================================================")
	fmt.Println("The following packages will be added to ansible:")
	fmt.Println("================================================")
	PrintSlice(pkgs)
	fmt.Println("================================================")
	fmt.Println("Continue? (y,N)")
	input, err := GetInput()
	if err != nil {
		fmt.Println("Error reading input.")
		panic(err)
	}
	if strings.TrimSpace(strings.ToLower(input)) == "y" {
		fmt.Println("yes")
		AddPackagesYml(pkgs, pkgListPath)
	}
}

// Add packages to the input yml file
func AddPackagesYml(pkgs []string, filePath string) {

	// open output file
	fo, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
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

func main() {
	SyncPacmanPackages()
}
