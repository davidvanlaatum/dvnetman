package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <file> <dir>\n", os.Args[0])
		os.Exit(1)
	}
	file, dir := os.Args[1], filepath.Clean(os.Args[2])+string(filepath.Separator)

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer f.Close()

	fileList := map[string]bool{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		fileList[s.Text()] = true
	}

	var extra []string
	err = filepath.Walk(
		dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Base(path) == ".openapi-generator" {
				return filepath.SkipDir
			}
			if info.IsDir() || slices.Contains(
				[]string{".openapi-generator-ignore", "openapi-generator.yaml", ".editorconfig"}, filepath.Base(path),
			) {
				return nil
			}
			if _, ok := fileList[strings.TrimPrefix(path, dir)]; !ok {
				extra = append(extra, path)
			}
			return nil
		},
	)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if len(extra) > 0 {
		fmt.Println("Extra files:")
		for _, f := range extra {
			fmt.Println(f)
		}
		os.Exit(1)
	}
}
