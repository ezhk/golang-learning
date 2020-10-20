package main

import (
	"bytes"
	"flag"
	"go/format"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

var FunctionsMap = template.FuncMap{
	"ToLower": strings.ToLower,
	"Title":   strings.Title,
}

var TestsTemplate = `// generated by generate-table-tests; DO NOT EDIT
// +build integration

package storage

import (
	"testing"

	"github.com/ezhk/golang-learning/banners-rotation/internal/config"
	"github.com/stretchr/testify/suite"
)

type {{.Name | Title}}TestSuite struct {
	suite.Suite
	db *Storage
}

func Test{{.Name | Title}}Suite(t *testing.T) {
	suite.Run(t, new({{.Name | Title}}TestSuite))
}

func (s *{{.Name | Title}}TestSuite) SetupTest() {
	cfg := config.NewConfig("testdata/config.yaml")

	db, err := NewStorage(cfg)
	s.NoError(err)

	// Define storage.
	s.db = db

	// Clean previous values.
	s.TearDownTest()
}

func (s *{{.Name | Title}}TestSuite) TearDownTest() {
	// Clean exists {{.Name | ToLower}}s.
	{{.Name | ToLower}}s, err := s.db.Read{{.Name | Title}}s()
	s.NoError(err)
	for _, {{.Name | ToLower}} := range {{.Name | ToLower}}s {
		err = s.db.Delete{{.Name | Title}}({{.Name | ToLower}}.ID)
		s.NoError(err)
	}
}

func (s *{{.Name | Title}}TestSuite) Test{{.Name | Title}}Operations() {
	// Create new {{.Name | ToLower}}.
	{{.Name | ToLower}}, err := s.db.Create{{.Name | Title}}("test {{.Name | ToLower}}")
	s.NoError(err)
	s.Equal("test {{.Name | ToLower}}", {{.Name | ToLower}}.Name)

	{{.Name | ToLower}}.Name = "updated test {{.Name | ToLower}}"
	updated{{.Name | Title}}, err := s.db.Update{{.Name | Title}}({{.Name | ToLower}})
	s.NoError(err)
	s.Equal("updated test {{.Name | ToLower}}", updated{{.Name | Title}}.Name)

	{{.Name | ToLower}}s, err := s.db.Read{{.Name | Title}}s()
	s.NoError(err)
	s.Greater(len({{.Name | ToLower}}s), 0)
	s.Equal("updated test {{.Name | ToLower}}", {{.Name | ToLower}}s[0].Name)

	// Call "duplicate key value violates unique constraint".
	_, err = s.db.Create{{.Name | Title}}("updated test {{.Name | ToLower}}")
	s.Error(err)

	err = s.db.Delete{{.Name | Title}}({{.Name | ToLower}}s[0].ID)
	s.NoError(err)
}
`

func generate(tableName string) ([]byte, error) {
	var buf bytes.Buffer

	tmpl, err := template.New("go-generate-table-tests").Funcs(FunctionsMap).Parse(TestsTemplate)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(&buf, struct{ Name string }{Name: tableName})
	if err != nil {
		return nil, err
	}

	source, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return source, nil
}

func main() {
	var (
		tableName string
		fileName  string
	)

	flag.StringVar(&tableName, "table", "banner", "table name")
	flag.StringVar(&fileName, "file", "banner.go", "output filepath")
	flag.Parse()

	generatedContent, err := generate(tableName)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(fileName, generatedContent, 0600); err != nil {
		log.Fatal(err)
	}
}
