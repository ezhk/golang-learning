package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Environment map[string]string

func ReadDir(dir string) (Environment, error) {
	// ReadDir reads a specified directory and returns map of env variables.
	// Variables represented as files where filename is name of variable, file first line is a value.
	env := make(Environment)
	dataDir, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, d := range dataDir {
		fileName := d.Name()

		// skip recursive directory search and files with "="
		if d.IsDir() || strings.Contains(fileName, "=") {
			continue
		}

		filePath := filepath.Join(dir, fileName)
		val, err := readFileValue(filePath)
		// skip file with error
		if err != nil {
			continue
		}
		env[fileName] = val
	}

	return env, nil
}

func readFileValue(filePath string) (string, error) {
	// function read file and return first "trail trimmed" string
	r, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	st, err := r.Stat()
	if err != nil {
		return "", err
	}

	// don't read empty files
	if st.Size() <= 0 {
		return "", nil
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	// read only first line if it doesn't not empty
	for _, val := range strings.Split(string(bytes), "\n") {
		if val == "" {
			continue
		}

		// convering 0x00 to \n magic
		val = strings.ReplaceAll(val, string('\x00'), "\n")

		return strings.TrimRightFunc(val, unicode.IsSpace), nil
	}

	return "", errors.New("empty result")
}
