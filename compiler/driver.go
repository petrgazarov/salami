package compiler

import (
	"fmt"
	"os"
	"path/filepath"
)

type driver struct {
	config compilerConfig
}

type target int

const (
	pulumiPython target = iota
)

var targetNames = map[target]string{
	pulumiPython: "PulumiPython",
}

type compilerConfig struct {
	target               target
	sourceDirectory      string
	destinationDirectory string
}

func newDriver(config compilerConfig) *driver {
	return &driver{
		config: config,
	}
}

func (driver *driver) compile() error {
	var files []string
	err := filepath.Walk(driver.config.sourceDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".sami" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(file)
	}

	return nil
}

func Compile() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	absModuleDir, err := filepath.Abs(cwd)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	sourceDirectory := filepath.Join(absModuleDir, "example")

	config := compilerConfig{
		target:               pulumiPython,
		sourceDirectory:      sourceDirectory,
		destinationDirectory: "./output",
	}

	driver := newDriver(config)
	if err := driver.compile(); err != nil {
		return fmt.Errorf("compilation error: %w", err)
	}

	return nil
}
