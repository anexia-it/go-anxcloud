package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const xxgenerated_object_identifier = `// DO NOT EDIT, auto generated

package {{.Package}}

import (
	"context"
)
{{range .IdentifiedObjects }}
// GetIdentifier returns the primary identifier of a {{.ObjectName}} object
func (x *{{.ObjectName}}) GetIdentifier(ctx context.Context) (string, error) {
	return x.{{.IdentifyingField}}, nil
}
{{end}}`

type identifiedObject struct {
	ObjectName       string
	IdentifyingField string
}

func init() {
	tools["object-identifier"] = func() {
		packages, err := filepath.Glob("./pkg/apis/*/v1")
		if err != nil {
			log.Fatalf("failed to list packages: %s", err)
		}

		t := template.Must(template.New("").Parse(xxgenerated_object_identifier))

		for _, pkg := range packages {
			objects := getIdentifiedObjectsFromPackage(pkg)

			out, err := os.OpenFile(path.Join(pkg, "xxgenerated_object_identifier.go"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				log.Fatalf("failed to open file: %s", err)
			}

			err = t.Execute(out, map[string]interface{}{
				"Package":           "v1",
				"IdentifiedObjects": objects,
			})

			if err != nil {
				log.Fatalf("failed to execute template: %s", err)
			}
		}
	}
}

func getIdentifiedObjectsFromPackage(pkg string) []identifiedObject {
	fileSet := token.NewFileSet()
	packages, err := parser.ParseDir(fileSet, pkg, func(fi fs.FileInfo) bool {
		return fi.Name() != "xxgenerated_object_identifier.go"
	}, parser.AllErrors|parser.ParseComments)
	if err != nil {
		log.Fatalf("unable to parse %s directory", pkg)
	}

	var v1pkg *ast.Package
	for _, pkg := range packages {
		if pkg.Name == "v1" {
			v1pkg = pkg
			break
		}
	}
	if v1pkg == nil {
		log.Fatal("unable to find v1 package")
	}

	alreadyIdentifiedObjects := filterAlreadyIdentifiedObjects(v1pkg.Files)
	identifiableObjects := filterIdentifiableObjects(v1pkg.Files, alreadyIdentifiedObjects)

	sort.Slice(identifiableObjects, func(i, j int) bool {
		return identifiableObjects[i].ObjectName < identifiableObjects[j].ObjectName
	})

	return identifiableObjects
}

func filterAlreadyIdentifiedObjects(files map[string]*ast.File) []string {
	alreadyIdentifiedObjects := make([]string, 0)
	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Name.String() == "GetIdentifier" && fn.Recv != nil && len(fn.Recv.List) == 1 {
					if r, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
						alreadyIdentifiedObjects = append(alreadyIdentifiedObjects, r.X.(*ast.Ident).Name)
					}
				}
			}

			return true
		})
	}
	return alreadyIdentifiedObjects
}

func filterIdentifiableObjects(files map[string]*ast.File, skipObjects []string) []identifiedObject {
	identifiableObjects := make([]identifiedObject, 0)

	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			if gen, ok := n.(*ast.GenDecl); ok && gen.Tok == token.TYPE && len(gen.Specs) > 0 {
				if spec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
					if contains(skipObjects, spec.Name.String()) {
						return true
					}
					if t, ok := spec.Type.(*ast.StructType); ok {
						for _, field := range t.Fields.List {
							if field.Tag != nil && strings.Contains(field.Tag.Value, "anxcloud:\"identifier\"") {
								identifiableObjects = append(identifiableObjects, identifiedObject{
									ObjectName:       spec.Name.String(),
									IdentifyingField: field.Names[0].String(),
								})
							}
						}
					}
				}
			}

			return true
		})
	}

	return identifiableObjects
}

func contains(haystack []string, needle string) bool {
	i := sort.SearchStrings(haystack, needle)
	return i < len(haystack) && haystack[i] == needle
}
