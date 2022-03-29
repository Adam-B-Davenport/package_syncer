package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)

func Check(err error, message string) {
	if err != nil {
		fmt.Println(message)
		panic(err)
	}
}

func GetInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}

func PrintStringSlice(slice []string) {
	for _, str := range slice {
		fmt.Println(str)
	}
}
func PrintIntSlice(slice []int) {
	for _, i := range slice {
		fmt.Println(i)
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

func WriteTextList(filePath string, list []string) error {
	dir := path.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, OS_USER_RWX)
		Check(err, "Could not create parent directory.")
	}
	fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0600)
	Check(err, "Error writing to file.")

	defer func() {
		if err := fp.Close(); err != nil {
			panic(err)
		}
	}()
	for _, str := range list {
		line := fmt.Sprintf("%s\n", str)
		_, err := fp.WriteString(line)
		Check(err, "Failed to write to file.")
	}
	return nil

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
