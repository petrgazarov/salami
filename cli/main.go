package main

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
)

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
	fmt.Println("Usage: \n \t salami <command>\n ")
	fmt.Println("The commands are:")
	for id, command_ := range command_list {
		fmt.Println(id, ":", command_["description"])
	}
	fmt.Println()
}

func main() {
	flag.Parse()
	command := flag.Arg(0)

	if command == "" {
		fmt.Println("===== Welcome to Salami-CLI =====\n ")
		fmt.Println("This tool provides a very simple command to generate your LLMs!\n ")
		showCommands()
		return
	}

	switch cmd := command; cmd {
	case command_list["version"]["cmd"]:
		fmt.Println("Salami version", "1.0")
	case command_list["compile"]["cmd"]:
		fmt.Println("compile")
	case "help":
		showCommands()
	default:
		color.Red("Invalid command passed. Type 'salami help'")
	}
}
