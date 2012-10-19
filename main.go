package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type FileInfo struct {
	Hash, Root, Path string
}

func main() {
	in := bufio.NewReader(os.Stdin)

	// [path][root]info
	files := make(map[string]map[string]FileInfo)

	for {
		line, err := in.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from stdin:", err)
			return
		}

		hash, filename := line[:40], line[42:]
		if !strings.Contains(filename, "/") {
			continue // Ignoring root-level files.
		}

		firstSlash := strings.Index(filename, "/")
		root, path := filename[:firstSlash], filename[firstSlash+1:]

		path = path[:len(path)-1] // Remove trailing newline

		switch path {
		case "addoninfo.txt", "addonimage.jpg":
			continue // Ignored filenames
		}

		if _, ok := files[path]; !ok {
			files[path] = make(map[string]FileInfo)
		}
		files[path][root] = FileInfo{
			Root: root,
			Path: path,
			Hash: hash,
		}
	}

	collisions := make(map[string][]FileInfo)
	duplicates := make(map[string][]FileInfo)

	for path, pathFiles := range files {
		if len(pathFiles) == 1 {
			continue
		}

	pathLoop:
		for root1, file1 := range pathFiles {
			for root2, file2 := range pathFiles {
				if root1 == root2 {
					continue
				}

				if file1.Hash != file2.Hash {
					current, _ := collisions[path]
					collisions[path] = append(current, file1)
					continue pathLoop
				}
			}
			current, _ := duplicates[path]
			duplicates[path] = append(current, file1)
		}
	}

	for path, instances := range duplicates {
		fmt.Printf("Duplicate %q:\nHash %s\n", path, instances[0].Hash)
		for _, instance := range instances {
			fmt.Printf("\t%s\n", instance.Root)
		}
		fmt.Println()
	}

	for path, instances := range collisions {
		fmt.Printf("Collision %q:\n", path)
		for _, instance := range instances {
			fmt.Printf("\t%s\n\t\tHash: %s\n", instance.Root, instance.Hash)
		}
		fmt.Println()
	}

}
