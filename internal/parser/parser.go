package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type StructField struct {
	Name string
	Type string
}

type StructWithFields struct {
	Name    string
	Fields  []StructField
	Package string
}

func (s StructWithFields) NameWithPackage() string {
	return fmt.Sprintf("%s.%s", s.Package, s.Name)
}

func (s StructWithFields) GetTableName() string {
	return strings.ToLower(s.Name) + "s"
}

func (s StructWithFields) FieldsToSqlRepresentation() string {
	var res string
	for i, v := range s.Fields {
		res += strings.ToLower(v.Name)
		if i != len(s.Fields)-1 {
			res += ","
		}
	}

	return res
}

func (s StructWithFields) FieldArgsToString() string {
	var res string
	for i, v := range s.Fields {
		res += strings.ToLower(string(s.Name[0])) + "." + v.Name
		if i != len(s.Fields)-1 {
			res += ", "
		}
	}

	return res
}

func (s StructWithFields) AsArgument() string {
	return strings.ToLower(string(s.Name[0])) + " " + s.NameWithPackage()
}

func (s StructWithFields) RepositoryName() string {
	return s.Name + "Repository"
}

func ParseGoFile(filename string) ([]StructWithFields, error) {
	var result []StructWithFields
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, file, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("Error parsing file: %v", err)
	}

	packageName := node.Name.Name

	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						var tmpFields []StructField
						for _, field := range structType.Fields.List {
							fieldType := fieldTypeToString(field.Type, fset)
							for _, name := range field.Names {
								tmpFields = append(tmpFields, StructField{name.Name, fieldType})
							}
						}
						tmp := StructWithFields{Name: typeSpec.Name.Name, Fields: tmpFields, Package: packageName}
						result = append(result, tmp)
					}
				}
			}
		}
	}
	return result, nil
}

func fieldTypeToString(expr ast.Expr, fset *token.FileSet) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		pkgName := t.X.(*ast.Ident).Name
		typeName := t.Sel.Name
		return fmt.Sprintf("%s.%s", pkgName, typeName)
	case *ast.ArrayType:
		return "[]" + fieldTypeToString(t.Elt, fset)
	case *ast.StarExpr:
		return "*" + fieldTypeToString(t.X, fset)
	default:
		// If the type is not defined in this package, return its original representation
		return strings.TrimSpace(expr.(*ast.Ident).Name)
	}
}
