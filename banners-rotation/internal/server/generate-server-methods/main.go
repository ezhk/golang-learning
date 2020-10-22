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

var gRPCTemplate = `// generated by generate-grpc-methods; DO NOT EDIT
package server

import (
	"context"

	"github.com/ezhk/golang-learning/banners-rotation/internal/structs"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Convert{{.Name | Title}}ToSimpleResponse(b structs.{{.Name | Title}}) *SimpleResponse {
	return &SimpleResponse{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		CreatedAt:   timestamppb.New(b.CreatedAt),
		UpdatedAt:   timestamppb.New(b.UpdatedAt),
	}
}

func ConvertSimpleUpdateRequestTo{{.Name | Title}}(r *SimpleUpdateRequest) structs.{{.Name | Title}} {
	return structs.{{.Name | Title}}{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
}

func (s Server) Create{{.Name | Title}}(ctx context.Context, r *SimpleCreateRequest) (*SimpleResponse, error) {
	{{.Name | ToLower}}, err := s.storage.Create{{.Name | Title}}(r.Name, r.Description)
	if err != nil {
		return nil, err
	}

	return Convert{{.Name | Title}}ToSimpleResponse({{.Name | ToLower}}), nil
}

func (s Server) Read{{.Name | Title}}s(ctx context.Context, empty *empty.Empty) (*MultipleSimpleResponse, error) {
	{{.Name | ToLower}}s, err := s.storage.Read{{.Name | Title}}s()
	if err != nil {
		return nil, err
	}

	simpleResponses := make([]*SimpleResponse, 0)
	for _, {{.Name | ToLower}} := range {{.Name | ToLower}}s {
		simpleResponses = append(simpleResponses, Convert{{.Name | Title}}ToSimpleResponse(*{{.Name | ToLower}}))
	}

	return &MultipleSimpleResponse{Objects: simpleResponses}, nil
}

func (s Server) Update{{.Name | Title}}(ctx context.Context, r *SimpleUpdateRequest) (*SimpleResponse, error) {
	{{.Name | ToLower}} := ConvertSimpleUpdateRequestTo{{.Name | Title}}(r)
	b, err := s.storage.Update{{.Name | Title}}({{.Name | ToLower}})
	if err != nil {
		return nil, err
	}

	return Convert{{.Name | Title}}ToSimpleResponse(b), nil
}

func (s Server) Delete{{.Name | Title}}(ctx context.Context, r *SimpleRequestID) (*SimpleResponseID, error) {
	if err := s.storage.Delete{{.Name | Title}}(r.ID); err != nil {
		return nil, err
	}

	return &SimpleResponseID{ID: r.ID}, nil
}
`

func generate(objectName string) ([]byte, error) {
	var buf bytes.Buffer

	tmpl, err := template.New("generate-grpc-methods").Funcs(FunctionsMap).Parse(gRPCTemplate)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(&buf, struct{ Name string }{Name: objectName})
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
		objectName string
		fileName   string
	)

	flag.StringVar(&objectName, "object", "banner", "object name")
	flag.StringVar(&fileName, "file", "banner.go", "output filepath")
	flag.Parse()

	generatedContent, err := generate(objectName)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(fileName, generatedContent, 0600); err != nil {
		log.Fatal(err)
	}
}
