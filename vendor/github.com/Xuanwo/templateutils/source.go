package templateutils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

// Source represent all content in the file.
type Source struct {
	Interfaces []*Interface
	Methods    []*Method

	r    *os.File
	fset *token.FileSet
}

// ParseContent will parse content.
func ParseContent(filename string, content []byte) (s *Source, err error) {
	s = &Source{
		fset: token.NewFileSet(),
	}

	r, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	// Keep open file in source so that we can get every code block content.
	s.r = r

	f, err := parser.ParseFile(s.fset, filename, string(content), parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	for _, decl := range f.Decls {
		switch v := decl.(type) {
		case *ast.GenDecl:
			err = s.parseGenDecl(v)
			if err != nil {
				return nil, err
			}
		case *ast.FuncDecl:
			err = s.parseFuncDecl(v)
			if err != nil {
				return nil, err
			}
		default:
			// Ignore unsupported decl.
			continue
		}
	}

	return s, nil
}

func (s *Source) parseGenDecl(g *ast.GenDecl) (err error) {
	// The Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec.
	//
	// Only support Type for now
	for _, spec := range g.Specs {
		switch v := spec.(type) {
		case *ast.TypeSpec:
			err = s.parseTypeSpec(v)
			if err != nil {
				return err
			}
		default:
			continue
		}
	}
	return nil
}

func (s *Source) parseTypeSpec(t *ast.TypeSpec) error {
	// Only support interface type for now.
	switch v := t.Type.(type) {
	case *ast.InterfaceType:
		in := &Interface{
			Name:   t.Name.Name,
			Method: make([]*Method, 0),
		}

		err := in.parseInterfaceType(v)
		if err != nil {
			return err
		}

		s.Interfaces = append(s.Interfaces, in)
	default:
		return nil
	}

	return nil
}

func (s *Source) parseFuncDecl(f *ast.FuncDecl) error {
	m := &Method{
		Name: f.Name.Name,
	}

	err := m.parseFuncDecl(f)
	if err != nil {
		return err
	}

	s.Methods = append(s.Methods, m)
	return nil
}

// Interface represent an interface.
type Interface struct {
	Name   string
	Method []*Method
}

func (i *Interface) parseInterfaceType(in *ast.InterfaceType) error {
	for _, v := range in.Methods.List {
		funcType, ok := v.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		m := &Method{}
		err := m.parseFuncType(funcType)
		if err != nil {
			return err
		}

		for _, name := range v.Names {
			xm := *m
			xm.Name = name.Name

			i.Method = append(i.Method, &xm)
		}
	}

	return nil
}

// Method represent a method.
type Method struct {
	Name    string
	Recv    *Recv
	Params  FieldList
	Results FieldList
}

func (m *Method) parseFuncType(f *ast.FuncType) error {
	m.Params = make(FieldList, 0)
	m.Results = make(FieldList, 0)

	if f.Params != nil {
		err := m.Params.parseFields(f.Params)
		if err != nil {
			return err
		}
	}

	if f.Results != nil {
		err := m.Results.parseFields(f.Results)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Method) parseFuncDecl(f *ast.FuncDecl) error {
	m.Params = make(FieldList, 0)
	m.Results = make(FieldList, 0)
	m.Recv = &Recv{}

	err := m.parseFuncType(f.Type)
	if err != nil {
		return err
	}
	if f.Recv != nil {
		err = m.Recv.parseFieldList(f.Recv)
		if err != nil {
			return err
		}
	}
	return nil
}

// Recv represent a receiver
type Recv struct {
	Name string
	Type string
}

func (r *Recv) String() string {
	return fmt.Sprintf("%s %s", r.Name, r.Type)
}

func (r *Recv) parseFieldList(f *ast.FieldList) error {
	// Only support one receiver for now
	if f.NumFields() > 1 {
		return fmt.Errorf("not supported recv")
	}
	field := f.List[0]

	r.Name = field.Names[0].Name
	r.Type = formatExpr(field.Type)
	return nil
}

// FieldList is a list for Field.
type FieldList []*Field

func (f FieldList) String() string {
	s := []string{}
	for _, v := range f {
		s = append(s, v.String())
	}
	return strings.Join(s, ",")
}

func (f *FieldList) parseFields(list *ast.FieldList) error {
	for _, v := range list.List {
		field := &Field{}
		err := field.parseField(v)
		if err != nil {
			return err
		}

		*f = append(*f, field)
	}
	return nil
}

// Field represent a field.
type Field struct {
	Names []string
	Type  string
}

func (f *Field) String() string {
	s := []string{}
	for _, name := range f.Names {
		s = append(s, name)
	}
	return strings.Join(s, ",") + " " + f.Type
}

func (f *Field) parseField(field *ast.Field) error {
	for _, v := range field.Names {
		f.Names = append(f.Names, v.Name)
	}

	f.Type = formatExpr(field.Type)
	return nil
}

func formatExpr(t ast.Expr) string {
	switch v := t.(type) {
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", formatExpr(v.X), v.Sel.Name)
	case *ast.Ident:
		return v.Name
	case *ast.StarExpr:
		return "*" + formatExpr(v.X)
	case *ast.Ellipsis:
		return "..." + formatExpr(v.Elt)
	case *ast.ArrayType:
		return "[]" + formatExpr(v.Elt)
	default:
		log.Fatalf("not handled type %+#v", v)
		return ""
	}
}
