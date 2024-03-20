package main

import (
	"go/ast"
	"go/token"
	"html/template"
	"log"
	"sort"
	"strings"
)

type identifiableObject struct {
	ObjectName       string
	IdentifyingField string
}

func (gen *ObjectGenerator) generateGetIdentifier() {
	const templates = `
{{range .IdentifiedObjects }}
// GetIdentifier returns the primary identifier of a {{.ObjectName}} object
func (o {{.ObjectName}}) GetIdentifier(ctx context.Context) (string, error) {
	return o.{{.IdentifyingField}}, nil
}
{{end}}`

	t := template.Must(template.New("").Parse(templates))

	objectsToIdentify := make([]identifiableObject, 0, len(gen.identifiableObjects)-len(gen.alreadyIdentifiedObjects))
	for _, obj := range gen.identifiableObjects {
		if contains(gen.alreadyIdentifiedObjects, obj.ObjectName) {
			continue
		}
		objectsToIdentify = append(objectsToIdentify, obj)
	}

	err := t.Execute(gen.out, map[string]interface{}{
		"IdentifiedObjects": objectsToIdentify,
	})
	if err != nil {
		log.Fatalf("Error executing template for object identifier: %v", err)
	}
}

func (gen *ObjectGenerator) findAlreadyIdentifiedObjects(file *ast.File, fset *token.FileSet, source string) {
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name.String() == "GetIdentifier" && fn.Recv != nil && len(fn.Recv.List) == 1 {
				if r, ok := fn.Recv.List[0].Type.(*ast.Ident); ok {
					gen.alreadyIdentifiedObjects = append(gen.alreadyIdentifiedObjects, r.Name)
				}
			}
		}

		return true
	})
}

func (gen *ObjectGenerator) findIdentifiableObjects(file *ast.File, fset *token.FileSet, source string) {
	ast.Inspect(file, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE && len(genDecl.Specs) > 0 {
			if spec, ok := genDecl.Specs[0].(*ast.TypeSpec); ok {
				if t, ok := spec.Type.(*ast.StructType); ok {
					if fieldName, ok := findIdentifyingFieldInObject(t); ok {
						gen.identifiableObjects = append(gen.identifiableObjects, identifiableObject{
							ObjectName:       spec.Name.String(),
							IdentifyingField: fieldName,
						})
					}
				}
			}
		}

		return true
	})
}

func findIdentifyingFieldInObject(t *ast.StructType) (string, bool) {
	for _, field := range t.Fields.List {
		if field.Tag != nil && strings.Contains(field.Tag.Value, "anxcloud:\"identifier\"") {
			return field.Names[0].String(), true
		}
	}
	return "", false
}

func contains(haystack []string, needle string) bool {
	i := sort.SearchStrings(haystack, needle)
	return i < len(haystack) && haystack[i] == needle
}
