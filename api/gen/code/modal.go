package code

import (
	"dvnetman/pkg/utils"
	. "github.com/dave/jennifer/jen"
)

type GoType struct {
	pkg     string
	name    string
	nested  *GoType
	pointer bool
}

func (t *GoType) ToCode() *Statement {
	p := func(c *Statement) *Statement {
		if t.pointer {
			return Op("*").Add(c)
		}
		return c
	}
	if t.pkg != "" {
		return p(Qual(t.pkg, t.name))
	}
	switch t.name {
	case "string":
		return p(String())
	case "int":
		return p(Int())
	case "bool":
		return p(Bool())
	case "float64":
		return p(Float64())
	case "array":
		return p(Index().Add(t.nested.ToCode()))
	default:
		return p(Id(t.name))
	}
}

type Field struct {
	name     string
	wireName string
	required bool
	goType   *GoType
}

type Modal struct {
	name   string
	fields []*Field
}

func (m *Modal) ToCode() *Statement {
	return Type().Id(m.name).StructFunc(
		func(g *Group) {
			for _, f := range m.fields {
				jsonTag := f.wireName
				if !f.required {
					jsonTag += ",omitzero"
				}
				g.Id(f.name).Add(f.goType.ToCode()).Tag(map[string]string{"json": jsonTag})
			}
		},
	)
}

type Enum struct {
	name   string
	items  []string
	GoType *GoType
}

func (e *Enum) ToCode() *Statement {
	return Do(
		func(s *Statement) {
			s.Type().Id(e.name).Add(e.GoType.ToCode()).Line()
			s.Var().DefsFunc(
				func(g *Group) {
					for _, item := range e.items {
						g.Id(e.name + utils.UCFirst(item)).Id(e.name).Op("=").Lit(item)
					}
				},
			)
		},
	)
}

type APIFuncParam struct {
	name     string
	wireName string
	goType   *GoType
	required bool
	in       string
}

type APIFunc struct {
	path   string
	method string
	name   string
	params []*APIFuncParam
}

func (f *APIFunc) hasRequiredParams() bool {
	for _, p := range f.params {
		if p.required && p.in != "path" && p.in != "body" {
			return true
		}
	}
	return false
}

func (f *APIFunc) requiredParams() []*APIFuncParam {
	var required []*APIFuncParam
	for _, p := range f.params {
		if p.required && p.in != "path" && p.in != "body" {
			required = append(required, p)
		}
	}
	return required
}

func (f *APIFunc) hasPathParams() bool {
	for _, p := range f.params {
		if p.in == "path" {
			return true
		}
	}
	return false
}

func (f *APIFunc) hasQueryParams() bool {
	for _, p := range f.params {
		if p.in == "query" {
			return true
		}
	}
	return false
}

func (f *APIFunc) hasHeaderParams() bool {
	for _, p := range f.params {
		if p.in == "header" {
			return true
		}
	}
	return false
}
