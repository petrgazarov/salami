package main

import (
	"flag"
	"fmt"
	"salami/common/driver"

	"github.com/fatih/color"
)

const version = "0.0.1"

var command_list = map[string]map[string]string{
	"version": {
		"cmd":         "version",
		"description": "Print the version of the installed Salami-CLI",
	},
	"compile": {
		"cmd":         "compile",
		"description": "Run the compilation end-to-end",
	},
}

func showCommands() {
	fmt.Println("Usage: \n \t salami-cli <command>\n ")
	fmt.Println("The commands are:")
	for id, command_ := range command_list {
		fmt.Println(color.HiBlueString(id), ":", command_["description"])
	}
	fmt.Println()
}

func runSystem() {
	errors := driver.Run()

	for _, err := range errors {
		color.Red(err.Error())
	}
}

func main() {
	flag.Parse()
	command := flag.Arg(0)

	if command == "" {
		color.Green("====================================")
		color.Green("======= Welcome to Salami-CLI ======")
		color.Green("====================================\n ")
		fmt.Println("Salami is a declarative domain-specific language for cloud infrastructure, \ncentered around natural language descriptions. You can think of Salami as \nwriting documentation for each cloud resource object, and letting the compiler \ntake care of converting that to structured code. The compiler uses LLM to convert \nnatural language into structured code, with Terraform currently as the supported target.\n ")
		showCommands()
		return
	}

	switch cmd := command; cmd {
	case command_list["version"]["cmd"]:
		fmt.Println("Salami v", version)
	case command_list["compile"]["cmd"]:
		runSystem()
	case "help":
		showCommands()
	default:
		color.Red("Invalid command passed. Type 'salami help'")
	}
}
