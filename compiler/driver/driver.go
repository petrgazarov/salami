package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"salami/compiler/config"
	"salami/compiler/lexer"
	"salami/compiler/parser"
	"salami/compiler/semantic_analyzer"
	"salami/compiler/symbol_table"
	"sync"
)

func Run() {
	compilerConfig := config.GetConfig()
	var files []string

	filepath.Walk(compilerConfig.SourceDir, func(path string, info os.FileInfo, err error) error {
		// TODO: handle error from filepath.Walk

		if !info.IsDir() && filepath.Ext(path) == ".sami" {
			files = append(files, path)
		}
		return nil
	})

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			processFile(file)
			wg.Done()
		}(file)
	}
	wg.Wait()
	symbolTable := symbol_table.NewSymbolTable(resources, variables)
	semanticAnalyzer := semantic_analyzer.NewSemanticAnalyzer(symbolTable)
	semanticAnalyzer.Analyze()
}

func processFile(filePath string) {
	fmt.Println("\nProcessing file:", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	lexerInstance := lexer.NewLexer(filePath, string(content))
	tokens := lexerInstance.Run()
	parserInstance := parser.NewParser(tokens)
	resources, variables, parsingError := parserInstance.Parse()
	if parsingError != nil {
		fmt.Println(parsingError)
		return
	}
}
