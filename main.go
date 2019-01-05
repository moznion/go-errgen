package errgen

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/moznion/go-struct-custom-tag-parser"
)

// Run generates code for errors from a struct that defines errors.
func Run(typ string, prefix string, outputFilePath string) {
	dir := "."
	p, err := build.Default.ImportDir(dir, 0)
	if err != nil {
		log.Fatalf("[ERROR] cannot process directory %s: %s", dir, err)
	}

	for _, goFile := range p.GoFiles {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, goFile, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("[ERROR] failed parsing a go file: filename=%s, err=%s", goFile, err)
		}

		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				// TODO
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					// TODO
					continue
				}
				structName := typeSpec.Name.Name
				if typ != structName {
					// TODO
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

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					// TODO
					continue
				}

				i := 0
				isFmtImported := false
				isErrorsImported := false
				for _, field := range structType.Fields.List {
					func() {
						defer func() {
							i++
						}()

						tagValue := field.Tag.Value[1 : len(field.Tag.Value)-1]
						tagKeyValue, err := tagparser.Parse(tagValue, true)
						if err != nil {
							// TODO be fatalf
							return
						}

						if tagKeyValue["obsoleted"] != "" {
							return
						}

						msg := tagKeyValue["errmsg"]
						if msg == "" {
							// TODO
							return
						}
						name := field.Names[0].Name

						vars := tagKeyValue["vars"]
						if vars != "" && !isFmtImported {
							header += "\nimport \"fmt\""
							isFmtImported = true
						} else if !isErrorsImported {
							header += "\nimport \"errors\""
							isErrorsImported = true
						}

						body += fmt.Sprintf("\n\nfunc %s(%s) error {\n"+
							"\treturn %s\n}",
							name,
							vars,
							constructMessageContents(i, vars, msg, prefix),
						)
					}()
				}
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

func constructMessageContents(i int, varsString string, msg string, prefix string) string {
	if varsString == "" {
		return fmt.Sprintf(`errors.New("[%s%d] %s")`, prefix, i, msg)
	}
	varNames, err := extractVarNames(varsString)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	return fmt.Sprintf(`fmt.Errorf("[%s%d] %s", %s)`, prefix, i, msg, strings.Join(varNames, ", "))
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
