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
		return
	}
}

func load(path string) error {
	log.Println("Ela Packager")
	pkger := &Config{}
	if err := pkger.LoadFrom(path); err != nil {
		return err
	}
	if err := pkger.Compile("."); err != nil {
		return err
	}
	return nil
}
