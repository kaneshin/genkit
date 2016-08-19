package templates

import "text/template"

var templates = map[string]string{"field.tmpl": `{{if .Struct}}struct {{initialCap .Name}}{{.Type}}
let {{.Name}}: {{initialCap .Name}}?
{{else}}let {{.Name}}: {{.Type}}{{end}}
`,
	"imports.tmpl": `{{range .}}import {{.}}
{{end}}
`,
	"init.tmpl": `init?(dictionary: [String: AnyObject]) {
{{range $name, $type := .Elements}}{{if eq "" $type}}    guard let {{$name}} = dictionary["{{$name}}"] as? [String: AnyObject] else {
        return nil
    }
    self.{{$name}} = {{initialCap $name}}(dictionary: {{$name}})
{{else}}    guard let {{$name}} = dictionary["{{$name}}"] as? {{$type}} else {
        return nil
    }
    self.{{$name}} = {{$name}}
{{end}}{{end}}
}
`,
	"protocol.tmpl": `protocol {{initialCap .TypeName}}: RequestType {
}

extension {{initialCap .TypeName}} {
    var baseURL: NSURL {
        return NSURL(string: "{{.URL}}")!
    }
}

`,
	"struct.tmpl": `{{$TypeName := .TypeName}}
{{$Name := .Name}}
struct {{initialCap .Name}}{{goType .Definition}}
{{range .Definition.Links}}
struct {{.Method}}{{initialCap $Name}}Request: {{initialCap $TypeName}} {
    typealias Response = {{initialCap $Name}}

    var method: HTTPMethod {
        return {{printf ".%s" .Method}}
    }

    var path: String {
        return "{{.HRef}}"
    }

    func responseFromObject(object: AnyObject, URLResponse: NSHTTPURLResponse) throws -> Response {
        guard let dictionary = object as? [String: AnyObject], let rateLimit = RateLimit(dictionary: dictionary) else {
            throw ResponseError.UnexpectedObject(object)
        }
        return rateLimit
    }
}{{end}}
`,
}

// Parse parses declared templates.
func Parse(t *template.Template) (*template.Template, error) {
	for name, s := range templates {
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		if _, err := tmpl.Parse(s); err != nil {
			return nil, err
		}
	}
	return t, nil
}

