package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dragmz/cpre"
	"github.com/pkg/errors"
)

type includes []string

type args struct {
	Path    string
	Include includes
}

func (i *includes) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *includes) String() string {
	return strings.Join(*i, ",")
}

func run(a args) error {
	if a.Path == "" {
		return errors.New("path is required")
	}

	p := cpre.NewPreprocessor(cpre.PreprocessorConfig{
		Include: cpre.NewIncluder(a.Include),
	})

	bs, err := os.ReadFile(a.Path)
	if err != nil {
		return errors.Wrapf(err, "failed to read file: '%s'", a.Path)
	}

	processed := p.Process(string(bs))

	fmt.Println(string(processed))

	return nil
}

func main() {
	var a args

	flag.StringVar(&a.Path, "path", "", "path to the file to be processed; required")
	flag.Var(&a.Include, "include", "include path; can be specified multiple times")

	flag.Parse()

	err := run(a)
	if err != nil {
		panic(err)
	}
}
