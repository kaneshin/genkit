package main

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"

	"github.com/kaneshin/genkit"
)

var (
	templates = genkit.Templates()
	newlines  = regexp.MustCompile(`(?m:\s*$)`)
)

// Generate generates code according to the schema.
func Generate(s *genkit.Schema) ([]byte, error) {
	var buf bytes.Buffer

	for i := 0; i < 2; i++ {
		s.Resolve(nil)
	}

	templates.ExecuteTemplate(&buf, "imports.tmpl", []string{
		"UIKit", "APIKit",
	})

	typeName := fmt.Sprintf("%sRequestType", s.ID)
	templates.ExecuteTemplate(&buf, "protocol.tmpl", map[string]string{
		"TypeName": typeName,
		"URL":      s.URL(),
	})

	for _, name := range sortedKeys(s.Properties) {
		schema := s.Properties[name]
		// Skipping definitions because there is no links, nor properties.
		if schema.Links == nil && schema.Properties == nil {
			continue
		}

		context := struct {
			TypeName   string
			Name       string
			Definition *genkit.Schema
		}{
			TypeName:   typeName,
			Name:       name,
			Definition: schema,
		}

		templates.ExecuteTemplate(&buf, "struct.tmpl", context)
	}

	return buf.Bytes(), nil

	// Remove blank lines added by text/template
	bytes := newlines.ReplaceAll(buf.Bytes(), []byte(""))

	return bytes, nil
	// Format sources
	// clean, err := format.Source(bytes)
	// if err != nil {
	// 	return buf.Bytes(), err
	// }
	// return clean, nil
}

func sortedKeys(m map[string]*genkit.Schema) (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}
