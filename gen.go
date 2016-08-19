//go:generate templates -s templates -o templates/templates.go
package genkit

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	bundle "github.com/kaneshin/genkit/templates"
)

var templates *template.Template

func init() {
	templates = template.New("field.tmpl").Funcs(helpers)
	templates = template.Must(bundle.Parse(templates))
}

// Templates retuns a reference of templates.
func Templates() *template.Template {
	return templates
}

// Resolve resolves reference inside the schema.
func (s *Schema) Resolve(r *Schema) *Schema {
	if r == nil {
		r = s
	}
	for n, d := range s.Definitions {
		s.Definitions[n] = d.Resolve(r)
	}
	for n, p := range s.Properties {
		s.Properties[n] = p.Resolve(r)
	}
	for n, p := range s.PatternProperties {
		s.PatternProperties[n] = p.Resolve(r)
	}
	if s.Items != nil {
		s.Items = s.Items.Resolve(r)
	}
	if s.Ref != nil {
		s = s.Ref.Resolve(r)
	}
	if len(s.OneOf) > 0 {
		s = s.OneOf[0].Ref.Resolve(r)
	}
	if len(s.AnyOf) > 0 {
		s = s.AnyOf[0].Ref.Resolve(r)
	}
	for _, l := range s.Links {
		l.Resolve(r)
	}
	return s
}

// Types returns the array of types described by this schema.
func (s *Schema) Types() (types []string, err error) {
	if arr, ok := s.Type.([]interface{}); ok {
		for _, v := range arr {
			types = append(types, v.(string))
		}
	} else if str, ok := s.Type.(string); ok {
		types = append(types, str)
	} else {
		err = fmt.Errorf("unknown type %v", s.Type)
	}
	return types, err
}

// GoType returns the Go type for the given schema as string.
func (s *Schema) GoType() string {
	return s.goType(true, true)
}

// IsCustomType returns true if the schema declares a custom type.
func (s *Schema) IsCustomType() bool {
	return len(s.Properties) > 0
}

func (s *Schema) goType(required bool, force bool) (goType string) {
	// Resolve JSON reference/pointer
	types, err := s.Types()
	if err != nil {
		panic(err)
	}
	for _, kind := range types {
		switch kind {
		case "boolean":
			goType = "Bool"
		case "string":
			switch s.Format {
			case "date-time":
				goType = "NSDate"
			default:
				goType = "String"
			}
		case "number":
			goType = "Double"
		case "integer":
			goType = "Int"
		case "any":
			goType = "Any"
		case "array":
			if s.Items != nil {
				goType = "[" + s.Items.goType(required, force) + "]"
			} else {
				goType = "[Any]"
			}
		case "object":
			// Check if patternProperties exists.
			if s.PatternProperties != nil {
				for _, prop := range s.PatternProperties {
					goType = fmt.Sprintf("map[string]%s", prop.GoType())
					break // We don't support more than one pattern for now.
				}
				continue
			}

			buf := bytes.NewBufferString(" {\n")
			elms := map[string]string{}
			for _, name := range sortedKeys(s.Properties) {
				prop := s.Properties[name]
				req := contains(name, s.Required) || force
				templates.ExecuteTemplate(buf, "field.tmpl", struct {
					Definition *Schema
					Struct     bool
					Name       string
					Type       string
				}{
					Definition: prop,
					Struct:     prop.Type.([]interface{})[0].(string) == "object",
					Name:       name,
					Type:       prop.goType(req, force),
				})
				if prop.Type.([]interface{})[0].(string) == "object" {
					elms[name] = ""
				} else {
					elms[name] = prop.goType(req, force)
				}
			}

			templates.ExecuteTemplate(buf, "init.tmpl", struct {
				Elements map[string]string
			}{
				Elements: elms,
			})

			buf.WriteString("}")
			goType = buf.String()
		case "null":
			continue
		default:
			panic(fmt.Sprintf("unknown type %s", kind))
		}
	}
	if goType == "" {
		panic(fmt.Sprintf("type not found : %s", types))
	}
	return goType
}

// Values returns function return values types.
func (s *Schema) Values(name string, l *Link) []string {
	var values []string
	name = returnType(name, s, l)
	if s.EmptyResult(l) {
		values = append(values, "error")
	} else if s.ReturnsCustomType(l) {
		values = append(values, fmt.Sprintf("*%s", name), "error")
	} else {
		values = append(values, s.ReturnedGoType(l), "error")
	}
	return values
}

// URL returns schema base URL.
func (s *Schema) URL() string {
	for _, l := range s.Links {
		if l.Rel == "self" {
			return l.HRef.String()
		}
	}
	return ""
}

// ReturnsCustomType returns true if the link returns a custom type.
func (s *Schema) ReturnsCustomType(l *Link) bool {
	if l.TargetSchema != nil {
		return len(l.TargetSchema.Properties) > 0
	}
	return len(s.Properties) > 0
}

// ReturnedGoType returns Go type returned by the given link as a string.
func (s *Schema) ReturnedGoType(l *Link) string {
	if l.TargetSchema != nil {
		return l.TargetSchema.goType(true, false)
	}
	return s.goType(true, false)
}

// EmptyResult retursn true if the link result should be empty.
func (s *Schema) EmptyResult(l *Link) bool {
	var (
		types []string
		err   error
	)
	if l.TargetSchema != nil {
		types, err = l.TargetSchema.Types()
	} else {
		types, err = s.Types()
	}
	if err != nil {
		return true
	}
	return len(types) == 1 && types[0] == "null"
}

// Parameters returns function parameters names and types.
func (l *Link) Parameters(name string) ([]string, map[string]string) {
	if l.HRef == nil {
		// No HRef property
		panic(fmt.Errorf("no href property declared for %s", l.Title))
	}
	var order []string
	params := make(map[string]string)
	for _, name := range l.HRef.Order {
		def := l.HRef.Schemas[name]
		order = append(order, name)
		params[name] = def.GoType()
	}
	if l.Schema != nil {
		order = append(order, "o")
		t, required := l.GoType()
		if l.AcceptsCustomType() {
			params["o"] = paramType(name, l)
		} else {
			params["o"] = t
		}
		if !required {
			params["o"] = "*" + params["o"]
		}
	}
	if l.Rel == "instances" && strings.ToUpper(l.Method) == "GET" {
		order = append(order, "lr")
		params["lr"] = "*ListRange"
	}
	return order, params
}

// AcceptsCustomType returns true if the link schema is not a primitive type
func (l *Link) AcceptsCustomType() bool {
	if l.Schema != nil && l.Schema.IsCustomType() {
		return true
	}
	return false
}

// Resolve resolve link schema and href.
func (l *Link) Resolve(r *Schema) {
	if l.Schema != nil {
		l.Schema = l.Schema.Resolve(r)
	}
	if l.TargetSchema != nil {
		l.TargetSchema = l.TargetSchema.Resolve(r)
	}
	l.HRef.Resolve(r)
}

// GoType returns Go type for the given schema as string and a bool specifying whether it is required
func (l *Link) GoType() (string, bool) {
	t := l.Schema.goType(true, false)
	if t[0] == '*' {
		return t[1:], false
	}
	return t, true
}
