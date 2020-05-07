package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	middleBlock = "├───"
	endBlock    = "└───"
	wall        = "│"
	tab         = "\t"
)

func printDir(out io.Writer, path string, printFiles bool, level int, openNode map[int]struct{}) {
	var new_path string
	listDir, err := ioutil.ReadDir(path)
	var listDirActual []os.FileInfo

	if err != nil {
		panic("Can't open dir")
	}

	if printFiles {
		listDirActual = listDir
	} else {
		for _, file := range listDir {
			if file.IsDir() {
				listDirActual = append(listDirActual, file)
			}
		}
	}

	for num, f := range listDirActual {
		var walls string
		isNotLastBlock := num != (len(listDirActual) - 1)

		if level > 0 {
			for i := 0; i < level; i++ {

				_, found := openNode[i]
				if found {
					walls += wall + tab
				} else {
					walls += tab
				}
			}
		}

		if isNotLastBlock {
			walls += middleBlock
		} else {
			walls += endBlock
		}

		if f.IsDir() {
			if isNotLastBlock {
				openNode[level] = struct{}{}
			} else {
				delete(openNode, level)
			}

			new_path = path + "/" + f.Name()
			fmt.Fprintln(out, walls+f.Name())
			printDir(out, new_path, printFiles, level+1, openNode)

		} else if printFiles {
			fmt.Fprint(out, walls+f.Name())
			if f.Size() != 0 {
				fmt.Fprintf(out, " (%db)\n", f.Size())
			} else {
				fmt.Fprint(out, " (empty)\n")
			}
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	level := 0
	var openNode = make(map[int]struct{})
	printDir(out, path, printFiles, level, openNode)
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
