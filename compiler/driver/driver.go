package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"salami/compiler/config"
	"salami/compiler/lexer"
	"salami/compiler/parser"
	"salami/compiler/semantic_analyzer"
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

	for _, file := range files {
		processFile(file)
	}
}

func processFile(filePath string) {
	fmt.Println("Processing file:", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tokenChannel := make(chan lexer.Token)

	go func() {
		lexerInstance := lexer.NewLexer(filePath, string(content))
		for {
			token := lexerInstance.NextToken()
			tokenChannel <- token
			if token.Type == lexer.EOF {
				close(tokenChannel)
				break
			}
		}
	}()

	irChannel := make(chan *parser.IR)

	parser := parser.NewParser(tokenChannel, irChannel)
	go parser.Parse()

	semanticAnalyzer := semantic_analyzer.NewSemanticAnalyzer(irChannel)
	semanticAnalyzer.Analyze()
}
