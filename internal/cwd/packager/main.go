package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Println("Argument should be the path to packager config json file")
		return
	}
	path := os.Args[1]
	if err := load(path); err != nil {
		log.Println(err.Error())
		panic(err)
	}
}

func load(path string) error {
	//log.Println("Ela Packager")
	pkg := NewPackage()
	if err := pkg.LoadConfig(path); err != nil {
		return err
	}
	if err := pkg.Compile("."); err != nil {
		return err
	}
	return nil
}
