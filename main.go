package errgen

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"

	"github.com/iancoleman/strcase"
	"github.com/moznion/go-struct-custom-tag-parser"
)

func Run(typ string, prefix string) {
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

package %s`,
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
				for _, field := range structType.Fields.List {
					func() {
						defer func() {
							i++
						}()

						tagValue := field.Tag.Value[1 : len(field.Tag.Value)-1]
						tagKeyValue, err := tagparser.Parse(tagValue, true)
						if err != nil {
							// TODO
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
							header += "\n\nimport \"fmt\""
							isFmtImported = true
						}

						body += fmt.Sprintf("\n\nfunc %s(%s) string {\n"+
							"\treturn %s\n}",
							name,
							func() string {
								if vars == "" {
									return ""
								}
								return fmt.Sprintf("%s string", vars)
							}(),
							func() string {
								if vars == "" {
									return fmt.Sprintf("`[%s%d] %s`", prefix, i, msg)
								}
								return fmt.Sprintf(`fmt.Sprintf("[%s%d] %s", %s)`, prefix, i, msg, vars)
							}(),
						)
					}()
				}
				err = ioutil.WriteFile(fmt.Sprintf("%s_err_gen.go", strcase.ToSnake(structName)), []byte(header+body), 0644)
				if err != nil {
					log.Fatalf("[ERROR] failed output generated code to a file")
				}
			}
		}
	}
}
