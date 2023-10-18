package main

import (
	"flag"
	"fmt"
	"salami/common/constants"
	"salami/common/driver"

	"github.com/fatih/color"
)

const GENERAL_COMMAND = "salami"

var command_list = map[string]map[string]string{
	"version": {
		"cmd":         "version",
		"description": "Print the version of the installed Salami CLI",
	},
	"compile": {
		"cmd":         "compile",
		"description": "Run the compilation end-to-end",
	},
}

func showCommands() {
	fmt.Println("Usage: \n \t ", GENERAL_COMMAND, "<command>\n ")
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
		color.Green("======= Welcome to Salami CLI ======")
		color.Green("====================================\n ")
		fmt.Println("Salami is a declarative domain-specific language for cloud infrastructure, " +
			"centered around natural language descriptions. You can think of Salami as writing " +
			"documentation for each cloud resource object, and letting the compiler take care of " +
			"converting that to IaC. The compiler uses LLM to convert natural language to IaC, " +
			"with Terraform currently as the supported target.\n ")
		showCommands()
		return
	}

	switch cmd := command; cmd {
	case command_list["version"]["cmd"]:
		fmt.Println("Salami version " + constants.SalamiVersion)
	case command_list["compile"]["cmd"]:
		runSystem()
	case "help":
		showCommands()
	default:
		msg := "Invalid command passed. Type '" + GENERAL_COMMAND + " help'"
		color.Red(msg)
	}
}
