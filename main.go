package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
)

// Retrieve all non aur packages
func BasePackages() ([]string, error) {

	pacman := exec.Command("pacman", "-Qent")
	awk := exec.Command("awk", "{print $1}")

	// Retrieve foreign packages and pass as args to grep -v
	aur, err := exec.Command("pacman", "-Qm").Output()
	if err != nil {
		return nil, err
	}
	args := fmt.Sprintf(" '%s'", string(aur))
	vgrep := exec.Command("grep", "-v", args)

	awk.Stdin, err = pacman.StdoutPipe()
	if err != nil {
		return nil, err
	}
	vgrep.Stdin, err = awk.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = pacman.Start()
	if err != nil {
		return nil, err
	}
	err = awk.Start()
	if err != nil {
		return nil, err
	}

	output, err := vgrep.Output()
	if err != nil {
		return nil, err
	}

	text := strings.TrimSpace(string(output))

	res := strings.Split(text, "\n")
	return res, nil
}

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

func PrintSlice(slice []string) {
	for _, str := range slice {
		fmt.Println(str)
	}
}

func PrintPackagePrompt(pkgs []string) {
	for i, pkg := range pkgs {
		fmt.Printf("[%d] %s\n", i, pkg)
	}
}

func ReadPackagsYml() ([]string, error) {
	home, _ := os.UserHomeDir()
	filePath := path.Join(home, "dev", "ansible-setup", "arch", "packages.yml")
	return ReadAnsibleList(filePath)
}

func Contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
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

func main() {
	ansChan := make(chan []string)
	insChan := make(chan []string)
	ignChan := make(chan []string)

	go func() {
		pkgs, err := ReadPackagsYml()
		if err != nil {
			fmt.Println("Error reading packages from yml.")
			panic(err)
		}
		ansChan <- pkgs
	}()
	go func() {
		pkgs, err := BasePackages()
		if err != nil {
			fmt.Println("Error reading packages from yml.")
			panic(err)
		}
		insChan <- pkgs

	}()
	go func() {
		if _, err := os.Stat("path/to/ignored"); os.IsExist(err) {

		} else {
			ignChan <- make([]string, 0)
		}
	}()
	ansiblePkgs := <-ansChan
	installedPkgs := <-insChan
	ignoredPkgs := <-ignChan
	// fmt.Println("Ansible")
	// PrintPackagePrompt(ansiblePkgs)
	// fmt.Println("installed")
	// PrintPackagePrompt(installedPkgs)
	// fmt.Println("ignored")
	// PrintPackagePrompt(ignoredPkgs)
	PrintPackagePrompt(ComparePackages(installedPkgs, ansiblePkgs, ignoredPkgs))
}
