<%=licenseText%>
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"reflect"
)

func contains(needle string, haystack []string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}

	return false
}

var fns = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
	"isLastReturnError": func(returns []*Arg) bool {
		l := len(returns)
		if l == 0 {
			return false
		}
		return returns[l-1].Type == "error"
	},
	"lastReturnName": func(returns []*Arg) string {
		l := len(returns)
		if l == 0 {
			return ""
		}
		return returns[l-1].Name
	},
}

func main() {
	file := flag.String("input", "", "path to the file containing the interface")
	target := flag.String("target", "", "name of the interface to use")
	ignore := flag.String("ignore", "", "ignore the following methods (separate with commas)")
	templateFlag := flag.String("template", "", "path to the template")
	formatCode := flag.Bool("format", true, "format output using gofmt")
	output := flag.String("output", "-", "path to the output file (use - for stdout)")
	imports := flag.String("imports", "-", "custom imports (separate with commas)")

	flag.Parse()

	ignoredMethods := []string{}
	if *ignore != "" {
		ignoredMethods = strings.Split(*ignore, ",")
	}

	customImports := []string{}
	if *imports != "" {
		customImports = strings.Split(*imports, ",")
	}

	tmpl, err := loadTemplate(*templateFlag)
	if err != nil {
		panic(err)
	}

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, *file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	methods := []*ast.Field{}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name != *target {
				break
			}

			i, ok := x.Type.(*ast.InterfaceType)
			if !ok {
				break
			}

			for _, m := range i.Methods.List {
				methods = append(methods, m)
			}
		}

		return true
	})

	m := []*Method{}
	for _, met := range methods {
		if t, ok := met.Type.(*ast.FuncType); ok {
			name := met.Names[0].Name

			if contains(name, ignoredMethods) {
				continue
			}

			counter := 0
			params := []*Arg{}
			for _, par := range t.Params.List {
				name := fmt.Sprintf("param%d", counter)
				if len(par.Names) > 0 {
					name = par.Names[0].Name
				}

				params = append(params, &Arg{Name: name, Type: getType(par.Type)})
				counter++
			}

			counter = 0
			returns := []*Arg{}
			for _, ret := range t.Results.List {
				name := fmt.Sprintf("result%d", counter)
				if len(ret.Names) > 0 {
					name = ret.Names[0].Name
				}

				returns = append(returns, &Arg{Name: name, Type: getType(ret.Type)})
				counter++
			}

			doclines := []string{}
			if met.Doc != nil && len(met.Doc.List) > 0 {
				for _, line := range met.Doc.List {
					doclines = append(doclines, line.Text)
				}
			}

			m = append(m, &Method{Name: name, Params: params, Returns: returns, Doc: doclines})
		}
	}

	t, err := template.New("t").Funcs(fns).Parse(tmpl)
	if err != nil {
		panic(err)
	}

	vm := &ViewModel{
		File:          *file,
		Target:        *target,
		Template:      *templateFlag,
		Output:        *output,
		Methods:       m,
		CustomImports: customImports,
	}

	var sink bytes.Buffer
	err = t.Execute(&sink, vm)
	if err != nil {
		panic(err)
	}

	b := sink.Bytes()
	if *formatCode {
		b, err = format.Source(b)
		if err != nil {
			fmt.Printf("unable to format generated code")
			panic(err)
		}
	}

	var out io.Writer = os.Stdout
	if *output != "-" {
		f, err := os.OpenFile(*output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Printf("unable to open output")
			panic(err)
		}
		defer f.Close()
		out = f
	}

	_, err = out.Write(b)
	if err != nil {
		fmt.Printf("unable to write output")
		panic(err)
	}
}

// ViewModel abstracts file
type ViewModel struct {
	File          string
	Target        string
	Template      string
	Output        string
	Methods       []*Method
	CustomImports []string
}

// Method logic
type Method struct {
	Name    string
	Doc     []string
	Params  []*Arg
	Returns []*Arg
}

// Arg to Method
type Arg struct {
	Name string
	Type string
}

func getType(n ast.Expr) string {
	switch x := n.(type) {
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", x.X, x.Sel)
	case *ast.Ident:
		return fmt.Sprintf("%s", x.Name)
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", getType(x.X))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", getType(x.Elt))
	}

	return fmt.Sprintf("unable to process type: %T", n)
}

func loadTemplate(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
