package spec

import "dvnetman/api/gen/openapi"

func init() {
	RegisterBuilder(BuildUser)
}

func BuildUser(api *OpenAPI) {
	api.AddEndpoint(
		AddEndpointOpts{
			Method:    "get",
			Path:      "/api/v1/user/current",
			Operation: "GetCurrentUser",
			Tags:      []string{"User"},
			Responses: map[string]openapi.Response{
				"200": {
					Description: "Get current user",
					Content: map[string]openapi.Content{
						json: {
							Schema: openapi.Schema{
								Ref: api.AddSchema(
									"CurrentUser",
									&openapi.Schema{
										Type: "object",
										Properties: map[string]openapi.Schema{
											"displayName": {
												Type: "string",
											},
											"email": {
												Type: "string",
											},
											"externalProvider": {
												Type: "string",
											},
											"externalID": {
												Type: "string",
											},
											"loggedIn": {
												Type: "boolean",
											},
											"profileImageURL": {
												Type:   "string",
												Format: "uri",
											},
										},
									},
								),
							},
						},
					},
				},
			},
			InSecure: true,
		},
	)
	api.AddEndpoint(
		AddEndpointOpts{
			Method:    "get",
			Path:      "/api/v1/user/providers",
			Operation: "GetUserProviders",
			Tags:      []string{"User"},
			InSecure:  true,
			Responses: map[string]openapi.Response{
				"200": {
					Description: "Get user providers",
					Content: map[string]openapi.Content{
						json: {
							Schema: openapi.Schema{
								Type: "array",
								Items: &openapi.Schema{
									Ref: api.AddSchema(
										"UserProvider", &openapi.Schema{
											Type: "object",
											Properties: map[string]openapi.Schema{
												"provider": {
													Type: "string",
												},
												"displayName": {
													Type: "string",
												},
												"loginURL": {
													Type:   "string",
													Format: "uri",
												},
												"loginButtonImageURL": {
													Type:   "string",
													Format: "uri",
												},
											},
											Required: []string{"provider", "displayName", "loginURL"},
										},
									),
								},
							},
						},
					},
				},
			},
		},
	)
	api.AddObject(
		AddObjectOpts{
			Name: "user",
			Tags: []string{"User"},
			SearchBody: &openapi.Schema{
				Type: "object",
				Properties: map[string]openapi.Schema{
					"ids": {
						Type: "array",
						Items: &openapi.Schema{
							Type:   "string",
							Format: "uuid",
						},
					},
					"firstName": {
						Type: "string",
					},
					"firstNameRegex": {
						Type: "string",
					},
					"lastName": {
						Type: "string",
					},
					"lastNameRegex": {
						Type: "string",
					},
					"displayName": {
						Type: "string",
					},
					"displayNameRegex": {
						Type: "string",
					},
					"email": {
						Type:   "string",
						Format: "email",
					},
					"fields": {
						Type: "array",
						Items: &openapi.Schema{
							Type: "string",
						},
					},
				},
			},
			SearchModal: &openapi.Schema{
				Type: "object",
				Properties: map[string]openapi.Schema{
					"id": {
						Type:   "string",
						Format: "uuid",
					},
					"firstName": {
						Type: "string",
					},
					"lastName": {
						Type: "string",
					},
					"email": {
						Type: "string",
					},
				},
				Required: []string{"id"},
			},
			GetModal: addCommonProps(
				&openapi.Schema{
					Type: "object",
					Properties: map[string]openapi.Schema{
						"firstName": {
							Type: "string",
						},
						"lastName": {
							Type: "string",
						},
						"email": {
							Type:   "string",
							Format: "email",
						},
						"password": {
							Type: "string",
						},
						"externalProvider": {
							Type: "string",
						},
						"externalID": {
							Type: "string",
						},
					},
				},
			),
		},
	)
}
