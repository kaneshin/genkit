package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/kaneshin/genkit"
)

const (
	NAME    = "genkit"
	OUTPUT  = "genkit_gen.go"
	PACKAGE = "genkit"
)

var (
	pkg     = flag.String("pkg", "", fmt.Sprintf("package name to use in the generated code. (default \"%s\")", PACKAGE))
	genPath = flag.String("path", "", "generate in a specific directory.")
	output  = flag.String("output", "", fmt.Sprintf("generate in a specific file name. (default \"%s\")", OUTPUT))
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", NAME)
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", NAME))
	flag.Usage = Usage
	flag.Parse()

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(3)
	}

	// validate `pkg'
	if len(*pkg) == 0 {
		*pkg = PACKAGE
	}

	switch {
	case len(*genPath) == 0:
		*genPath = filepath.Dir(os.Args[0])
	case len(*genPath) > 0:
		if d, err := filepath.Abs(*genPath); err == nil {
			*genPath = d
		}
	}

	os.Exit(run())
}

func run() int {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	log.SetFlags(0)
	log.SetPrefix("genkit: ")

	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("missing schema file")
	}

	var i io.Reader
	var err error
	if flag.Arg(0) == "-" {
		i = os.Stdin
	} else {
		if i, err = os.Open(flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
	}

	var s genkit.Schema
	d := json.NewDecoder(i)
	if err := d.Decode(&s); err != nil {
		log.Fatal(err)
	}

	code, err := Generate(&s)
	if err != nil {
		log.Fatal(err)
	}

	var o io.Writer
	if *output == "" {
		o = os.Stdout
	} else {
		dest := path.Join(*genPath, *output)
		if o, err = os.Create(dest); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Fprintln(o, string(code))

	return 0
}
