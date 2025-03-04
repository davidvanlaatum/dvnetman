package code

import (
	"bufio"
	"bytes"
	"context"
	"dvnetman/api/gen/openapi"
	"dvnetman/pkg/file"
	"dvnetman/pkg/logger"
	"dvnetman/pkg/utils"
	. "github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

const netHttpPkg = "net/http"
const muxPkg = "github.com/gorilla/mux"
const uuidPkg = "github.com/google/uuid"
const errorsPkg = "github.com/pkg/errors"
const utilsPkg = "dvnetman/pkg/utils"
const loggerPkg = "dvnetman/pkg/logger"
const responseStruct = "Response"

type CodeGen struct {
	modals    map[string]*Modal
	enums     map[string]*Enum
	files     map[string]*File
	funcs     map[string]*APIFunc
	funcOrder []*APIFunc
	apis      []string
}

func NewCodeGen() *CodeGen {
	return &CodeGen{}
}

func (c *CodeGen) getFile(name string) *File {
	if c.files == nil {
		c.files = map[string]*File{}
	}
	if c.files[name] == nil {
		f := NewFile("openapi")
		c.files[name] = f
		f.HeaderComment("Code generated by dvnetman. DO NOT EDIT.")
	}
	return c.files[name]
}

func (c *CodeGen) Generate(ctx context.Context, api *openapi.OpenAPI) (err error) {
	if err = c.generateModals(ctx, api); err != nil {
		return
	}
	if err = c.generateAPI(ctx, api); err != nil {
		return
	}
	for k, v := range c.modals {
		f := c.getFile(k)
		f.Add(v.ToCode())
	}
	e := c.getFile("enum")
	for _, v := range c.enums {
		e.Add(v.ToCode())
	}
	if err = c.renderAPI(); err != nil {
		return
	}
	c.renderErrorConverter()
	c.renderErrors()
	return nil
}

func makeIdentifier(s string) string {
	s = strings.ReplaceAll(s, "-", "")

	for {
		i := strings.Index(s, "_")
		if i == -1 {
			break
		}
		s = s[:i] + utils.UCFirst(s[i+1:])
	}

	return s
}

func (c *CodeGen) generateAPIMethod(path, method string, e openapi.Endpoint) (err error) {
	if _, found := c.funcs[e.OperationID]; found {
		return errors.Errorf("Duplicate operation ID %s", e.OperationID)
	}
	f := &APIFunc{
		name:   utils.UCFirst(e.OperationID),
		path:   path,
		api:    utils.UCFirst(e.Tags[0]) + "API",
		method: method,
	}
	if !slices.Contains(c.apis, utils.UCFirst(e.Tags[0])+"API") {
		c.apis = append(c.apis, utils.UCFirst(e.Tags[0])+"API")
		sort.Strings(c.apis)
	}
	for _, p := range e.Parameters {
		t := c.determineGoTypeFor(p.Schema)
		if !p.Required && t.name != "array" {
			t.pointer = true
		}
		f.params = append(
			f.params, &APIFuncParam{
				name:     utils.UCFirst(makeIdentifier(p.Name)),
				wireName: p.Name,
				goType:   t,
				required: p.Required,
				in:       p.In,
			},
		)
	}
	if e.RequestBody != nil {
		for _, p := range utils.MapSortedByKey(e.RequestBody.Content, utils.MapSortedByKeyString) {
			t := c.determineGoTypeFor(&p.Value.Schema)
			t.pointer = true
			f.params = append(
				f.params, &APIFuncParam{
					name:     "Body",
					goType:   t,
					required: e.RequestBody.Required,
					in:       "body",
				},
			)
		}
	}
	c.funcs[f.name] = f
	c.funcOrder = append(c.funcOrder, f)
	return
}

func (c *CodeGen) generateAPI(ctx context.Context, api *openapi.OpenAPI) (err error) {
	c.funcs = map[string]*APIFunc{}
	for _, v := range api.Paths.Sorted() {
		logger.Info(ctx).Msgf("Generating API for %s", v.Key)
		for _, m := range v.Value {
			if err = c.generateAPIMethod(v.Key, m.Key, m.Value); err != nil {
				return
			}
		}
	}
	return nil
}

func (c *CodeGen) renderOptsStructs() {
	for _, apiFunc := range c.funcOrder {
		if len(apiFunc.params) > 0 {
			f := c.getFile(apiFunc.api)
			f.Type().Id(apiFunc.name + "Opts").StructFunc(
				func(g *Group) {
					for _, p := range apiFunc.params {
						g.Id(p.name).Add(p.goType.ToCode())
					}
				},
			)
		}
	}
}

func (c *CodeGen) renderAPIInterface() {
	for _, a := range c.apis {
		af := c.getFile(a)
		af.Type().Id(a).InterfaceFunc(
			func(g *Group) {
				for _, v := range utils.Filter(
					utils.MapSortedByKey(c.funcs, utils.MapSortedByKeyString),
					func(v utils.SortedMapEntries[string, *APIFunc]) bool {
						return v.Value.api == a
					},
				) {
					g.Id(v.Key).Params(
						Id("ctx").Qual("context", "Context"), Do(
							func(g *Statement) {
								if len(v.Value.params) > 0 {
									g.Id("opts").Op("*").Id(v.Key + "Opts")
								}
							},
						),
					).Params(
						Id("res").Op("*").Id(responseStruct), Err().Error(),
					)
				}
			},
		)
	}
}

func (c *CodeGen) renderResponse() {
	f := c.getFile("response")
	f.Type().Id(responseStruct).StructFunc(
		func(g *Group) {
			g.Id("Code").Int()
			g.Id("Headers").Qual(netHttpPkg, "Header")
			g.Id("Object").Interface()
		},
	)

	f.Func().Params(Id("res").Op("*").Id(responseStruct)).Id("Write").Params(
		Id("r").Op("*").Qual(netHttpPkg, "Request"),
		Id("w").Qual(netHttpPkg, "ResponseWriter"),
	).Params(
		Err().Error(),
	).BlockFunc(
		func(g *Group) {
			g.Var().Id("data").Index().Byte()
			g.If(Id("res").Dot("Object").Op("!=").Nil()).BlockFunc(
				func(g *Group) {
					g.Id("w").Dot("Header").Call().Dot("Set").Call(Lit("Content-Type"), Lit("application/json"))
					g.If(
						List(Id("data"), Err()).Op("=").Qual("encoding/json", "MarshalIndent").Call(
							Id("res").Dot("Object"), Lit(""), Lit("  "),
						),
						Err().Op("!=").Nil(),
					).BlockFunc(
						func(g *Group) {
							g.Return(Qual(errorsPkg, "Wrap").Call(Err(), Lit("failed to marshal response data")))
						},
					)
					g.Id("w").Dot("Header").Call().Dot("Set").Call(
						Lit("Content-Length"), Qual("strconv", "FormatUint").Call(
							Uint64().Call(Len(Id("data"))),
							Lit(10),
						),
					)
				},
			)
			g.If(Id("res").Dot("Code").Op("==").Lit(0)).Block(
				Return(
					Qual(
						errorsPkg, "New",
					).Call(Lit("response code not set")),
				),
			)
			g.For(List(Id("k"), Id("v")).Op(":=").Range().Id("res").Dot("Headers")).BlockFunc(
				func(g *Group) {
					g.For(List(Id("_"), Id("value")).Op(":=").Range().Id("v")).BlockFunc(
						func(g *Group) {
							g.Id("w").Dot("Header").Call().Dot("Add").Call(Id("k"), Id("value"))
						},
					)
				},
			)
			g.Id("w").Dot("WriteHeader").Call(Id("res").Dot("Code"))
			g.If(Id("data").Op("!=").Nil()).BlockFunc(
				func(g *Group) {
					g.List(Id("_"), Err()).Op("=").Id("w").Dot("Write").Call(Id("data"))
				},
			)
			g.Return()
		},
	)
}

func (c *CodeGen) renderAPIFuncPathParam(v *APIFunc, p *APIFuncParam, g *Group) {
	if p.goType.name == "UUID" {
		g.If(
			List(Id("opts").Dot(p.name), Err()).Op("=").Qual(
				uuidPkg, "Parse",
			).Call(Id("vars").Index(Lit(p.wireName))), Err().Op("!=").Nil(),
		).BlockFunc(
			func(g *Group) {
				g.Id("h").Dot("errors").Dot("ErrorHandler").Call(
					Id("w"), Id("r"),
					Qual(errorsPkg, "WithStack").Call(Id("NewPathParamError").Call(Lit(p.name), Err())),
				)
				g.Return()
			},
		)
	} else {
		panic(errors.Errorf("Unhandled path param type %s", p.goType.name))
	}
}

func (c *CodeGen) renderAPIFuncParam(v *APIFunc, p *APIFuncParam, g *Group) {
	switch p.in {
	case "path":
		c.renderAPIFuncPathParam(v, p, g)
	case "body":
		c.renderAPIFuncBodyParam(v, p, g)
	case "query":
	case "header":
	default:
		panic(errors.Errorf("Unhandled param type %s", p.in))
	}
}

func (c *CodeGen) renderAPIFunc(v *APIFunc) *Statement {
	f := c.getFile(v.api)
	return f.Func().Params(Id("h").Op("*").Id(v.api+"Handler")).Id(v.name).Params(
		Id("w").Qual(netHttpPkg, "ResponseWriter"),
		Id("r").Op("*").Qual(netHttpPkg, "Request"),
	).BlockFunc(
		func(g *Group) {
			g.Var().Id("res").Op("*").Id(responseStruct)
			g.Var().Err().Error()
			c.renderAPIFuncParams(g, v)
			g.If(
				List(
					Id("res"), Err(),
				).Op("=").Id("h").Dot("service").Dot(v.name).Call(
					Id("r").Dot("Context").Call(), Do(
						func(s *Statement) {
							if len(v.params) > 0 {
								s.Id("opts")
							}
						},
					),
				),
				Err().Op("!=").Nil(),
			).BlockFunc(
				func(g *Group) {
					g.Id("h").Dot("errors").Dot("ErrorHandler").Call(Id("w"), Id("r"), Err())
				},
			).Else().If(Id("res").Op("==").Nil()).BlockFunc(
				func(g *Group) {
					g.Id("h").Dot("errors").Dot("ErrorHandler").Call(
						Id("w"), Id("r"), Qual(errorsPkg, "Errorf").Call(Lit("no response returned")),
					)
				},
			).Else().If(Err().Op("=").Id("res").Dot("Write").Call(Id("r"), Id("w")), Err().Op("!=").Nil()).BlockFunc(
				func(g *Group) {
					g.Qual(
						loggerPkg, "Info",
					).Call(Id("r").Dot("Context").Call()).Dot("Err").Call(Err()).Dot("Msg").Call(Lit("error writing response"))
				},
			)
		},
	)
}

func (c *CodeGen) renderAPIFuncParams(g *Group, v *APIFunc) {
	if len(v.params) > 0 {
		g.Id("opts").Op(":=").Op("&").Id(v.name + "Opts").Block()
		if v.hasPathParams() {
			g.Id("vars").Op(":=").Qual(muxPkg, "Vars").Call(Id("r"))
		}
		if v.hasRequiredParams() {
			g.Var().List(
				utils.MapTo(
					v.requiredParams(), func(p *APIFuncParam) Code {
						return Id(p.name + "Present")
					},
				)...,
			).Bool()
		}
		for _, p := range v.params {
			c.renderAPIFuncParam(v, p, g)
		}
		if v.hasQueryParams() {
			c.renderAPIFuncQueryParams(g, v)
		}
		if v.hasHeaderParams() {
			c.renderAPIFuncHeaderParams(g, v)
		}
	}
}

func (c *CodeGen) renderAPIFuncHeaderParams(g *Group, v *APIFunc) {
	g.For(
		List(Id("k"), Id("v")).Op(":=").Range().Id("r").Dot("Header").BlockFunc(
			func(g *Group) {
				g.Switch(Id("k")).BlockFunc(
					func(g *Group) {
						for _, p := range v.params {
							if p.in != "header" {
								continue
							}
							g.Case(Lit(p.wireName)).BlockFunc(c.renderAPIFuncHeaderParam(p))
						}
					},
				)
			},
		),
	)
}

func (c *CodeGen) renderAPIFuncHeaderParam(p *APIFuncParam) func(*Group) {
	return func(g *Group) {
		ptr := func(s *Statement) *Statement {
			if p.goType.pointer {
				return Qual(utilsPkg, "ToPtr").Call(s)
			}
			return s
		}
		if p.goType.name == "string" {
			g.Id("opts").Dot(p.name).Op("=").Add(ptr(Id("v").Index(Lit(0))))
		} else if p.goType.name == "int" {
			g.Var().Id("x").Int()
			g.If(
				List(Id("x"), Err()).Op("=").Qual(
					"strconv", "Atoi",
				).Call(Id("v").Index(Lit(0))), Err().Op("!=").Nil(),
			).BlockFunc(
				func(g *Group) {
					g.Id("h").Dot("errors").Dot("ErrorHandler").Call(
						Id("w"), Id("r"),
						Qual(errorsPkg, "WithStack").Call(
							Id("NewQueryParamError").Call(
								Lit(p.wireName), Err(),
							),
						),
					)
					g.Return()
				},
			)
			g.Id("opts").Dot(p.name).Op("=").Add(
				ptr(Id("x")),
			)
		} else if p.goType.name == "Time" {
			g.Var().Id("t").Qual("time", "Time")
			g.If(
				List(Id("t"), Err()).Op("=").Qual(
					"time", "Parse",
				).Call(Qual("time", "RFC1123"), Id("v").Index(Lit(0))), Err().Op("!=").Nil(),
			).BlockFunc(
				func(g *Group) {
					g.Id("h").Dot("errors").Dot("ErrorHandler").Call(
						Id("w"), Id("r"),
						Qual(errorsPkg, "WithStack").Call(
							Id("NewQueryParamError").Call(
								Lit(p.wireName), Err(),
							),
						),
					)
					g.Return()
				},
			)
			g.Id("opts").Dot(p.name).Op("=").Add(ptr(Id("t")))
		} else {
			panic(errors.Errorf("Unhandled query param type %s", p.goType.name))
		}
	}
}

func (c *CodeGen) renderAPIFuncQueryParams(g *Group, v *APIFunc) {
	g.For(
		List(Id("k"), Id("v")).Op(":=").Range().Id("r").Dot("URL").Dot("Query").Call().BlockFunc(
			func(g *Group) {
				g.Switch(Id("k")).BlockFunc(
					func(g *Group) {
						for _, p := range v.params {
							if p.in != "query" {
								continue
							}
							g.Case(Lit(p.wireName)).BlockFunc(c.renderAPIFuncQueryParam(p))
						}
					},
				)
			},
		),
	)
}

func (c *CodeGen) renderAPIFuncQueryParam(p *APIFuncParam) func(g *Group) {
	return func(g *Group) {
		ptr := func(s *Statement) *Statement {
			if p.goType.pointer {
				return Qual(utilsPkg, "ToPtr").Call(s)
			}
			return s
		}
		if p.goType.name == "string" {
			g.Id("opts").Dot(p.name).Op("=").Add(ptr(Id("v").Index(Lit(0))))
		} else if p.goType.name == "int" {
			g.Var().Id("x").Int()
			g.If(
				List(Id("x"), Err()).Op("=").Qual(
					"strconv", "Atoi",
				).Call(Id("v").Index(Lit(0))), Err().Op("!=").Nil(),
			).BlockFunc(
				func(g *Group) {
					g.Id("h").Dot("errors").Dot("ErrorHandler").Call(
						Id("w"), Id("r"),
						Qual(errorsPkg, "WithStack").Call(
							Id("NewQueryParamError").Call(
								Lit(p.wireName), Err(),
							),
						),
					)
					g.Return()
				},
			)
			g.Id("opts").Dot(p.name).Op("=").Add(
				ptr(Id("x")),
			)
		} else if p.goType.name == "array" && p.goType.nested.name == "string" {
			g.Id("opts").Dot(p.name).Op("=").Append(Id("opts").Dot(p.name), ptr(Id("v")).Op("..."))
		} else if p.goType.name == "array" && p.goType.nested.name == "UUID" {
			g.If(
				List(Id("opts").Dot(p.name), Err()).Op("=").Qual(utilsPkg, "MapErr").Call(Id("v"), Qual(uuidPkg, "Parse")),
				Err().Op("!=").Nil(),
			).BlockFunc(
				func(g *Group) {
					g.Id("h").Dot("errors").Dot("ErrorHandler").Call(
						Id("w"), Id("r"), Qual(errorsPkg, "WithStack").Call(Id("NewQueryParamError").Call(Lit(p.wireName), Err())),
					)
					g.Return()
				},
			)
		} else {
			panic(errors.Errorf("Unhandled query param type %s: %s", p.goType.name, p.goType.GoString()))
		}
	}
}

func (c *CodeGen) renderRouter() {
	for _, api := range c.apis {
		f := c.getFile(api)
		f.Func().Id("Attach"+api).Call(
			Id("service").Id(api), Id("errors").Id("ErrorHandler"), Id("router").Op("*").Qual(muxPkg, "Router"),
		).BlockFunc(
			func(g *Group) {
				g.Id("handler").Op(":=").Op("&").Id(api + "Handler").Values(
					Dict{
						Id("service"): Id("service"),
						Id("errors"):  Id("errors"),
					},
				)
				for _, v := range c.funcOrder {
					if v.api == api {
						g.Id("router").
							Dot("Methods").Call(Qual(netHttpPkg, "Method"+utils.UCFirst(v.method))).
							Dot("Path").Call(Lit(v.muxPath())).
							Dot("Name").Call(Lit(v.name)).
							Dot("HandlerFunc").Call(Id("handler").Dot(v.name))
					}
				}
			},
		)
	}
}

func (c *CodeGen) renderHandler() {
	for _, a := range c.apis {
		f := c.getFile(a)
		f.Type().Id(a + "Handler").StructFunc(
			func(g *Group) {
				g.Id("service").Id(a)
				g.Id("errors").Id("ErrorHandler")
			},
		)
	}
	for _, v := range c.funcOrder {
		c.renderAPIFunc(v)
	}
}

func (c *CodeGen) renderAPI() (err error) {
	c.renderAPIInterface()
	c.renderOptsStructs()
	c.renderResponse()
	c.renderHandler()
	c.renderRouter()
	return nil
}

func (c *CodeGen) determineGoTypeFor(schema *openapi.Schema) *GoType {
	switch schema.Type {
	case "string":
		if schema.Format == "date-time" {
			return &GoType{pkg: "time", name: "Time"}
		} else if schema.Format == "uuid" {
			return &GoType{pkg: "github.com/google/uuid", name: "UUID"}
		} else if schema.Format != "" && schema.Format != "email" && schema.Format != "uri" {
			panic(errors.Errorf("Unhandled format %s", schema.Format))
		} else if schema.Enum != nil {
			return &GoType{name: "enum"}
		}
		return &GoType{name: "string"}
	case "integer":
		return &GoType{name: "int"}
	case "boolean":
		return &GoType{name: "bool"}
	case "array":
		return &GoType{name: "array", nested: c.determineGoTypeFor(schema.Items)}
	case "number":
		return &GoType{name: "float64"}
	default:
		if schema.Ref != "" {
			return &GoType{name: utils.UCFirst(strings.TrimPrefix(schema.Ref, "#/components/schemas/"))}
		}
		panic(errors.Errorf("Unhandled schema %#v", schema))
	}
}

func (c *CodeGen) isBasicType(t *GoType) bool {
	if t.name == "string" || t.name == "int" || t.name == "bool" || t.name == "float64" {
		return true
	}
	if t.pkg == uuidPkg && t.name == "UUID" {
		return true
	}
	return false
}

func (c *CodeGen) generateModals(ctx context.Context, api *openapi.OpenAPI) (err error) {
	c.modals = map[string]*Modal{}
	c.enums = map[string]*Enum{}

	for _, v := range utils.MapSortedByKey(api.Components.Schemas, utils.MapSortedByKeyString) {
		logger.Info(ctx).Msgf("Generating modal for %s", v.Key)
		x := &Modal{name: v.Key}

		for _, p := range utils.MapSortedByKey(v.Value.Properties, utils.MapSortedByKeyString) {
			t := c.determineGoTypeFor(&p.Value)
			f := &Field{name: utils.UCFirst(p.Key), wireName: p.Key, goType: t}
			f.required = v.Value.Required != nil && utils.Contains(v.Value.Required, p.Key)
			if t.name != "array" && !f.required {
				t.pointer = true
			}
			if t.name == "array" && !c.isBasicType(t.nested) {
				t.nested.pointer = true
			}
			x.fields = append(x.fields, f)
			if f.goType.name == "enum" {
				e := &Enum{name: v.Key + utils.UCFirst(p.Key)}
				e.items = p.Value.Enum
				e.GoType = &GoType{name: "string"}
				f.goType.name = e.name
				c.enums[e.name] = e
			}
		}
		c.modals[v.Key] = x
	}
	return nil
}

func (c *CodeGen) renderAPIFuncBodyParam(v *APIFunc, p *APIFuncParam, g *Group) {
	g.Id("decoder").Op(":=").Qual("encoding/json", "NewDecoder").Call(Id("r").Dot("Body"))
	g.Id("decoder").Dot("DisallowUnknownFields").Call()
	g.If(
		Err().Op("=").Id("decoder").Dot("Decode").Call(
			Op("&").Id("opts").Dot(p.name),
		), Err().Op("!=").Nil(),
	).BlockFunc(
		func(g *Group) {
			g.Id("h").Dot("errors").Dot("ErrorHandler").Call(
				Id("w"), Id("r"),
				Qual(errorsPkg, "WithStack").Call(
					Id("NewBodyParamError").Call(Err()),
				),
			)
			g.Return()
		},
	)
	g.If(Id("decoder").Dot("More").Call()).BlockFunc(
		func(g *Group) {
			g.Id("h").Dot("errors").Dot("ErrorHandler").Call(
				Id("w"), Id("r"),
				Qual(errorsPkg, "WithStack").Call(
					Id("NewBodyParamError").Call(
						Qual(errorsPkg, "New").Call(Lit("unexpected data after body")),
					),
				),
			)
			g.Return()
		},
	)
}

func (c *CodeGen) renderErrorConverter() {
	f := c.getFile("errors")
	f.Type().Id("ErrorConverterFunc").Func().Params(Error()).Op("*").Id("Response")
	f.Var().Id("errorConverters").Index().Id("ErrorConverterFunc")
	f.Func().Id("RegisterErrorConverter").Params(Id("ec").Id("ErrorConverterFunc")).BlockFunc(
		func(g *Group) {
			g.Id("errorConverters").Op("=").Append(Id("errorConverters"), Id("ec"))
		},
	)
	f.Type().Id("ErrorConverter").Struct()
	f.Func().Params(Id("ec").Op("*").Id("ErrorConverter")).Id("ErrorHandler").Params(
		Id("w").Qual(netHttpPkg, "ResponseWriter"), Id("r").Op("*").Qual(netHttpPkg, "Request"), Err().Error(),
	).BlockFunc(
		func(g *Group) {
			g.For(List(Id("_"), Id("converter"))).Op(":=").Range().Id("errorConverters").BlockFunc(
				func(g *Group) {
					g.If(Id("res").Op(":=").Id("converter").Call(Err()), Id("res").Op("!=").Nil()).BlockFunc(
						func(g *Group) {
							g.If(
								Err().Op(":=").Id("res").Dot("Write").Call(Id("r"), Id("w")), Err().Op("!=").Nil(),
							).BlockFunc(
								func(g *Group) {
									g.Qual(
										loggerPkg, "Error",
									).Call(Id("r").Dot("Context").Call()).Dot("Err").Call(Err()).Dot("Msg").Call(Lit("error writing error response"))
									g.Return()
								},
							)
							g.Return()
						},
					)
				},
			)
			g.Qual(loggerPkg, "Error").Call(Id("r").Dot("Context").Call()).Dot("Msg").Call(Lit("no error converter found"))
			g.Qual(netHttpPkg, "Error").Call(
				Id("w"),
				Err().Dot("Error").Call(), Qual(netHttpPkg, "StatusInternalServerError"),
			)
		},
	)
}

func (c *CodeGen) registerErrorHandler(f *File, errorType string, statusCode string, code string) {
	varName := utils.LCFirst(errorType)
	f.Func().Id("init").Params().BlockFunc(
		func(g *Group) {
			g.Id("RegisterErrorConverter").Call(
				Func().Params(Err().Error()).Params(Op("*").Id("Response")).BlockFunc(
					func(g *Group) {
						g.Var().Id(varName).Op("*").Id(errorType)
						g.If(Id("ok").Op(":=").Qual(errorsPkg, "As").Call(Err(), Id(varName)), Id("ok")).BlockFunc(
							func(g *Group) {
								g.Return(
									Op("&").Id("Response").Values(
										Dict{
											Id("Code"): Qual(netHttpPkg, statusCode),
											Id("Object"): Id("APIErrorModal").Values(
												Dict{
													Id("Errors"): Index().Op("*").Id("ErrorMessage").Values(
														Values(
															Dict{
																Id("Code"):    Lit(code),
																Id("Message"): Err().Dot("Error").Call(),
															},
														),
													),
												},
											),
										},
									),
								)
							},
						)
						g.Return(Nil())
					},
				),
			)
		},
	)
}

func (c *CodeGen) renderErrors() {
	f := c.getFile("errors")

	f.Type().Id("ErrorHandler").InterfaceFunc(
		func(g *Group) {
			g.Id("ErrorHandler").Params(
				Id("w").Qual(netHttpPkg, "ResponseWriter"), Id("r").Op("*").Qual(netHttpPkg, "Request"), Err().Error(),
			)
		},
	)

	f.Type().Id("PathParamError").StructFunc(
		func(g *Group) {
			g.Id("name").String()
			g.Err().Error()
		},
	)

	f.Func().Params(Id("e").Id("PathParamError")).Id("Error").Params().String().BlockFunc(
		func(g *Group) {
			g.Return(
				Qual("fmt", "Sprintf").Call(
					Lit("invalid path param %s: %v"), Id("e").Dot("name"), Id("e").Dot("err"),
				),
			)
		},
	)

	f.Func().Id("NewPathParamError").Params(
		Id("name").String(),
		Id("err").Error(),
	).Error().BlockFunc(
		func(g *Group) {
			g.Return(
				Op("&").Id("PathParamError").Values(
					Dict{
						Id("name"): Id("name"),
						Err():      Err(),
					},
				),
			)
		},
	)

	c.registerErrorHandler(f, "PathParamError", "StatusBadRequest", "invalid_path_param")

	f.Type().Id("QueryParamError").StructFunc(
		func(g *Group) {
			g.Id("name").String()
			g.Err().Error()
		},
	)

	f.Func().Params(Id("e").Id("QueryParamError")).Id("Error").Params().String().BlockFunc(
		func(g *Group) {
			g.Return(
				Qual("fmt", "Sprintf").Call(
					Lit("invalid Query param %s: %v"), Id("e").Dot("name"), Id("e").Dot("err"),
				),
			)
		},
	)

	f.Func().Id("NewQueryParamError").Params(
		Id("name").String(),
		Id("err").Error(),
	).Error().BlockFunc(
		func(g *Group) {
			g.Return(
				Op("&").Id("QueryParamError").Values(
					Dict{
						Id("name"): Id("name"),
						Err():      Err(),
					},
				),
			)
		},
	)

	c.registerErrorHandler(f, "QueryParamError", "StatusBadRequest", "invalid_query_param")

	f.Type().Id("BodyParamError").StructFunc(
		func(g *Group) {
			g.Err().Error()
		},
	)

	f.Func().Params(Id("e").Id("BodyParamError")).Id("Error").Params().String().BlockFunc(
		func(g *Group) {
			g.Return(
				Qual("fmt", "Sprintf").Call(
					Lit("invalid body param: %v"), Id("e").Dot("err"),
				),
			)
		},
	)

	f.Func().Id("NewBodyParamError").Params(
		Id("err").Error(),
	).Error().BlockFunc(
		func(g *Group) {
			g.Return(
				Op("&").Id("BodyParamError").Values(
					Dict{
						Err(): Err(),
					},
				),
			)
		},
	)

	c.registerErrorHandler(f, "BodyParamError", "StatusBadRequest", "invalid_body_param")
}

func (c *CodeGen) isGeneratedFile(path string) (ok bool, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	if s.Scan() {
		return strings.Contains(s.Text(), "Code generated by dvnetman"), nil
	}
	err = s.Err()
	return
}

func (c *CodeGen) removeOldFiles(ctx context.Context, path string, files map[string]string) (err error) {
	var dir []os.DirEntry
	if dir, err = os.ReadDir(path); err != nil {
		return errors.Wrap(err, "failed to read directory")
	}
	for _, entry := range dir {
		if !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		if _, ok := files[entry.Name()]; !ok {
			if ok, err = c.isGeneratedFile(filepath.Join(path, entry.Name())); err != nil {
				logger.Error(ctx).Err(err).Msgf("Failed to check file %s", entry.Name())
				return
			} else if ok {
				logger.Info(ctx).Msgf("Removing file %s", entry.Name())
				if err = os.Remove(filepath.Join(path, entry.Name())); err != nil {
					return errors.Wrap(err, "failed to remove file")
				}
			} else {
				logger.Debug(ctx).Msgf("Not deleting extra file %s, not generated", entry.Name())
			}
		}
	}
	return
}

func (c *CodeGen) WriteFiles(ctx context.Context, path string) (err error) {
	files := map[string]string{}
	for _, v := range utils.MapSortedByKey(c.files, utils.MapSortedByKeyString) {
		logger.Info(ctx).Msgf("Rendering file %s.go", v.Key)
		buf := &bytes.Buffer{}
		if err = v.Value.Render(buf); err != nil {
			return
		}
		files[v.Key+".go"] = buf.String()
	}
	for _, v := range utils.MapSortedByKey(files, utils.MapSortedByKeyString) {
		logger.Info(ctx).Msgf("Writing file %s", v.Key)
		func() {
			var f file.FileUpdate
			if f, err = file.NewFileUpdate(filepath.Join(path, v.Key), file.OnlyGenerated); err != nil {
				return
			}
			defer f.Close()
			_, err = f.WriteString(v.Value)
		}()
		if err != nil {
			return
		}
	}
	err = c.removeOldFiles(ctx, path, files)
	return
}
