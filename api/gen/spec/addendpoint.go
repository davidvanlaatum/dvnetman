package spec

import "dvnetman/api/gen/openapi"

type AddEndpointOpts struct {
	Method      string
	Path        string
	Operation   string
	Tags        []string
	Responses   map[string]openapi.Response
	Parameters  []openapi.Parameter
	RequestBody *openapi.RequestBody
	InSecure    bool
}

func (o *OpenAPI) AddEndpoint(opts AddEndpointOpts) {
	if o.Paths[opts.Path] == nil {
		o.Paths[opts.Path] = map[string]openapi.Endpoint{}
	}
	if _, found := o.Paths[opts.Path][opts.Method]; found {
		panic("endpoint " + opts.Method + " " + opts.Path + " already exists")
	}
	if _, found := opts.Responses["400"]; !found {
		opts.Responses["400"] = badRequestResponse
	}
	if _, found := opts.Responses["500"]; !found {
		opts.Responses["500"] = internalServerErrorResponse
	}
	endpoint := openapi.Endpoint{
		OperationID: opts.Operation,
		Tags:        opts.Tags,
		Responses:   opts.Responses,
		Parameters:  opts.Parameters,
		RequestBody: opts.RequestBody,
	}
	if !opts.InSecure {
		endpoint.Security = []map[string][]string{
			{
				"bearerAuth": {},
			},
		}
	}
	o.Paths[opts.Path][opts.Method] = endpoint
}
