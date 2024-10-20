package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors" //nolint: depguard // import is necessary
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to os.ReadDir")
	}

	env := make(Environment)

	for _, entry := range dirEntries {
		entryInfo, err := entry.Info()
		if err != nil {
			return nil, errors.Wrap(err, "failed to entry.Info")
		}

		if entryInfo.IsDir() {
			continue
		}

		fileName := entryInfo.Name()

		// имя S не должно содержать =
		if strings.Contains(fileName, "=") {
			log.Printf("skipping %s, contains '=' character", fileName)
			continue
		}

		filePath := filepath.Join(dir, fileName)

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to os.ReadFile")
		}

		fileContentString := string(fileContent)

		// учитываем только первую строку
		fileContentString = strings.Split(fileContentString, "\n")[0]

		// пробелы и табуляция в конце T удаляются
		fileContentString = strings.TrimRight(fileContentString, " \t")

		// терминальные нули (0x00) заменяются на перевод строки (\n)
		fileContentString = strings.ReplaceAll(fileContentString, string([]byte{0}), "\n")

		env[fileName] = EnvValue{
			Value:      fileContentString,
			NeedRemove: entryInfo.Size() == 0,
		}
	}

	return env, nil
}
