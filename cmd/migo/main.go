package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/walkline/migo"
)

var currentPath string

func handleNewCommand() {
	if len(flag.Args()) < 3 {
		fmt.Println("type and name required")
		os.Exit(1)
	}

	tmpl := migo.Templater{}
	err := tmpl.LoadTemplates(currentPath)
	if err != nil {
		panic(err)
	}

	if ver == "-1" {
		v, err := tmpl.TampleteWithType(migo.TemplateTypeVersion).BuildVersion()
		if err != nil {
			panic(err)
		}
		ver = v
	}

	name := strings.Replace(flag.Args()[2], " ", "-", -1)
	v, err := migo.VersionFromString(ver + "-" + name)
	if err != nil {
		panic(err)
	}

	switch flag.Args()[1] {
	case "sql":
		upData, err := tmpl.ContentForTemplateType(migo.TemplateTypeSQLUp, v)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.up.sql", currentPath, v), upData, 0755)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}

		downData, err := tmpl.ContentForTemplateType(migo.TemplateTypeSQLDown, v)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.down.sql", currentPath, v), downData, 0755)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}
	case "go":
		goData, err := tmpl.ContentForTemplateType(migo.TemplateTypeGo, v)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.go", currentPath, v), goData, 0755)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}

	default:
		fmt.Printf("unsupported file: %v", err)
		os.Exit(1)
	}
}

func handleInitCommand() {
	t := migo.Templater{}
	err := t.LoadTemplates(currentPath)
	if err != nil {
		panic(err)
	}

}

var ver string

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	currentPath = dir

	var verPtr = flag.String("version", "-1", "set version manualy")
	flag.Parse()
	ver = *verPtr

	if len(flag.Args()) < 1 {
		fmt.Println(`new subcommand is required

Usage sample:
	$ migo new go "[MK-2014] Create users table"
	$ migo new sql "[MK-2015] Clean users table"`)
		os.Exit(1)
	}

	switch flag.Args()[0] {
	case "new":
		handleNewCommand()
	case "init":
		handleInitCommand()
	default:
		os.Exit(1)
	}
}
