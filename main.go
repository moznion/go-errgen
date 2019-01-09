package errgen

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

// Run generates code for errors from a struct that defines errors.
func Run(typ string, prefix string, outputFilePath string) {
	dir := "."
	p, err := build.Default.ImportDir(dir, 0)
	if err != nil {
		log.Fatalf("[ERROR] cannot process directory %s: %s", dir, err)
	}

	msgs := make([]string, 0)
	for _, goFile := range p.GoFiles {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, goFile, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("[ERROR] failed parsing a go file: filename=%s, err=%s", goFile, err)
		}

		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structName := typeSpec.Name.Name
				if typ != structName {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				pkgName := f.Name.Name
				header := fmt.Sprintf(`// This package was auto generated.
// DO NOT EDIT BY YOUR HAND!

package %s
`,
					pkgName,
				)
				body := ""

				i := 1
				isFmtImported := false
				isErrorsImported := false
				for _, field := range structType.Fields.List {
					func() {
						defer func() {
							i++
						}()

						tagValue := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

						if tagValue.Get("obsoleted") != "" {
							return
						}

						name := field.Names[0].Name
						msg := tagValue.Get("errmsg")
						if msg == "" {
							log.Printf("[WARN] `errmsg` tag is missing at `%s` field", name)
							return
						}

						vars := tagValue.Get("vars")
						if vars != "" && !isFmtImported {
							header += "\nimport \"fmt\""
							isFmtImported = true
						} else if !isErrorsImported {
							header += "\nimport \"errors\""
							isErrorsImported = true
						}

						msgCore, msgCode := constructMessageContents(i, vars, msg, prefix)
						body += fmt.Sprintf("\n\n// %s returns the error.\nfunc %s(%s) error {\n"+
							"\treturn %s\n}",
							name,
							name,
							vars,
							msgCode,
						)
						msgs = append(msgs, msgCore)
					}()
				}

				funcName := strcase.ToCamel(structName)
				body += fmt.Sprintf("\n\n// %sList returns the list of errors.\nfunc %sList() []string {\n\treturn []string{\n", funcName, funcName)
				for _, m := range msgs {
					body += fmt.Sprintf("\t\t`%s`,\n", m)
				}
				body += "\t}\n}\n"

				dst := fmt.Sprintf("%s_errmsg_gen.go", strcase.ToSnake(structName))
				if outputFilePath != "" {
					dst = outputFilePath
				}

				err = ioutil.WriteFile(dst, []byte(header+body), 0644)
				if err != nil {
					log.Fatalf("[ERROR] failed output generated code to a file")
				}
			}
		}
	}
}

func constructMessageContents(i int, varsString string, msg string, prefix string) (string, string) {
	msgCore := fmt.Sprintf("[%s%d] %s", prefix, i, msg)
	if varsString == "" {
		return msgCore, fmt.Sprintf(`errors.New("%s")`, msgCore)
	}
	varNames, err := extractVarNames(varsString)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	return msgCore, fmt.Sprintf(`fmt.Errorf("%s", %s)`, msgCore, strings.Join(varNames, ", "))
}

func extractVarNames(varsString string) ([]string, error) {
	vars := strings.Split(varsString, ",")
	varNames := make([]string, len(vars))

	for i, v := range vars {
		leaves := strings.Split(strings.TrimSpace(v), " ")
		if len(leaves) != 2 {
			return nil, fmt.Errorf("invalid syntax of vars has detected: given=%s", v)
		}
		varNames[i] = leaves[0]
	}

	return varNames, nil
}
