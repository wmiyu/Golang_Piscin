package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func symlinkTreatment(path string, name string, isDir bool, flag bool) {

	if isDir {
		return
	}

	fi, err := os.Lstat(filepath.Join(path, "/", name))
	if err != nil {
		fmt.Println("Error while get stat of file", err)
		os.Exit(0)
	}

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {

		filePath := filepath.Join(path, fi.Name())

		originFile, err := filepath.EvalSymlinks(filePath)
		if err != nil {
			originFile = "[broken]"
		}

		fmt.Println(filePath, "->", originFile)
	} else if flag {
		fmt.Println(path + name)
	}
}

func readDir(path string, flags map[string]bool, ext string) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error while read dir", err)
		os.Exit(0)
	}

	for _, files := range files {

		switch {
		case flags["D"]:
			if files.IsDir() {
				fmt.Println(path + "/" + files.Name())
			} else {
				break
			}
		case flags["F"] && ext != "":
			if filepath.Ext(files.Name()) == "."+ext {
				fmt.Println(path + "/" + files.Name())
			} else {
				break
			}
		case flags["F"]:
			if !files.IsDir() {
				fmt.Println(path + "/" + files.Name())
			} else {
				break
			}
		case flags["SL"]:
			symlinkTreatment(path+"/", files.Name(), files.IsDir(), false)
		default:
			symlinkTreatment(path+"/", files.Name(), files.IsDir(), true)
		}
		if files.IsDir() {
			readDir(path+"/"+files.Name(), flags, ext)
		}
	}
}

func main() {
	flagSl := flag.Bool("sl", false, "Find only sym-links")
	flagD := flag.Bool("d", false, "Find only directories")
	flagF := flag.Bool("f", false, "Find only files")
	flagExt := flag.String("ext", "", "Find files by extenstion")

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Error params: No path specified")
		os.Exit(0)
	}

	path := flag.Args()[0]

	flags := map[string]bool{"SL": *flagSl, "D": *flagD, "F": *flagF}

	if *flagExt != "" && !flags["F"] {
		fmt.Println("Error params: Required flag -f")
		os.Exit(0)
	}

	readDir(path, flags, *flagExt)
}
