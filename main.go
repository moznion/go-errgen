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
	g "github.com/moznion/gowrtr/generator"
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

				root := g.NewRoot(
					g.NewComment(" This package was auto generated."),
					g.NewComment(" DO NOT EDIT BY YOUR HAND!"),
					g.NewNewline(),
					g.NewPackage(pkgName),
					g.NewNewline(),
					g.NewImport("github.com/pkg/errors"),
				)

				listFuncReturnItems := make([]string, 0)
				funcNames := make([]string, 0)
				identifiers := make([]string, 0)

				i := 1
				for _, field := range structType.Fields.List {
					func() {
						defer func() {
							i++
						}()

						tagValue := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

						if tagValue.Get("obsoleted") != "" {
							return
						}

						funcName := field.Names[0].Name
						funcNames = append(funcNames, funcName)
						msg := tagValue.Get("errmsg")
						if msg == "" {
							log.Printf("[WARN] `errmsg` tag is missing at `%s` field", funcName)
							return
						}

						vars := tagValue.Get("vars")
						msgCore, msgCode, identifier := constructMessageContents(i, vars, msg, prefix)
						identifiers = append(identifiers, identifier)

						root = root.AddStatements(
							g.NewNewline(),
							g.NewCommentf(" %s returns the error.", funcName),
						)

						funcSig := g.NewFuncSignature(funcName).ReturnTypes("error")
						wrapFuncSig := g.NewFuncSignature(funcName + "Wrap").ReturnTypes("error")
						for _, v := range strings.Split(vars, ",") {
							if v == "" {
								continue
							}
							p := strings.Split(strings.TrimSpace(v), " ")
							if len(p) != 2 {
								log.Fatalf("invalid syntax of vars has detected: given=%s", v)
							}
							funcSig = funcSig.AddParameters(g.NewFuncParameter(p[0], p[1]))
							wrapFuncSig = wrapFuncSig.AddParameters(g.NewFuncParameter(p[0], p[1]))
						}
						wrapFuncSig = wrapFuncSig.AddParameters(g.NewFuncParameter("err", "error"))

						root = root.AddStatements(
							g.NewFunc(nil, funcSig, g.NewReturnStatement(msgCode)),
							g.NewFunc(nil, wrapFuncSig, g.NewReturnStatement(
								fmt.Sprintf("errors.Wrap(%s, err.Error())", msgCode),
							)),
						)

						listFuncReturnItems = append(listFuncReturnItems, msgCore)
					}()
				}

				funcName := strcase.ToCamel(structName)

				root = root.AddStatements(
					g.NewNewline(),
					g.NewRawStatementf("type %sType int", funcName),
					g.NewNewline(),
					g.NewRawStatement("const ("),
				)

				for i, fname := range funcNames {
					root = root.AddStatements(
						g.NewNewline(),
						g.NewRawStatementf("%sType", fname).WithNewline(false),
					)
					if i == 0 {
						root = root.AddStatements(g.NewRawStatementf(" %sType = iota", funcName).WithNewline(false))
					}
				}
				root = root.AddStatements(
					g.NewNewline(),
					g.NewRawStatementf("%sUnknownType", funcName),
					g.NewNewline(),
					g.NewRawStatement(")"),
				)

				root = root.AddStatements(
					g.NewNewline(),
					g.NewCommentf(" List%s returns the list of errors.", funcName),
					g.NewFunc(
						nil,
						g.NewFuncSignature("List"+funcName).ReturnTypes("[]string"),
						// TODO use composite literal
						g.NewReturnStatement(fmt.Sprintf("%#v", listFuncReturnItems)),
					),
					g.NewNewline(),
					g.NewCommentf(" Identify%s checks the identity of an error", funcName),
				)

				identifyFuncSig := g.NewFuncSignature("Identify" + funcName).
					Parameters(g.NewFuncParameter("err", "error")).
					ReturnTypes(funcName + "Type")

				switchStmt := g.NewSwitch("")
				for i, fname := range funcNames {
					switchStmt = switchStmt.AddCase(g.NewCase(
						fmt.Sprintf(`strings.HasPrefix(errStr, "%s")`, identifiers[i]),
						g.NewReturnStatement(fmt.Sprintf("%sType", fname)),
					))
				}
				switchStmt = switchStmt.Default(
					g.NewDefaultCase(g.NewReturnStatement(fmt.Sprintf("%sUnknownType", funcName))),
				)

				root = root.AddStatements(
					g.NewFunc(
						nil,
						identifyFuncSig,
						g.NewRawStatementf("errStr := err.Error()"),
						switchStmt,
					),
				)

				generated, err := root.Gofmt("-s").Goimports().Generate(0)
				if err != nil {
					log.Fatalf("[ERROR] failed to generate code: err=%s", err)
				}

				dst := fmt.Sprintf("%s_errmsg_gen.go", strcase.ToSnake(structName))
				if outputFilePath != "" {
					dst = outputFilePath
				}

				err = ioutil.WriteFile(dst, []byte(generated), 0644)
				if err != nil {
					log.Fatalf("[ERROR] failed output generated code to a file")
				}
			}
		}
	}
}

func constructMessageContents(i int, varsString string, msg string, prefix string) (string, string, string) {
	identifier := fmt.Sprintf("[%s%d]", prefix, i)
	msgCore := fmt.Sprintf("%s %s", identifier, msg)
	if varsString == "" {
		return msgCore, fmt.Sprintf("errors.New(`%s`)", msgCore), identifier
	}
	varNames, err := extractVarNames(varsString)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	return msgCore, fmt.Sprintf("fmt.Errorf(`%s`, %s)", msgCore, strings.Join(varNames, ", ")), identifier
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
