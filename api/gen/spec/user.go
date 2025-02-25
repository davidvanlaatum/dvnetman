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
								Ref: "#/components/schemas/User",
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
					"fullName": {
						Type: "string",
					},
					"fullNameRegex": {
						Type: "string",
					},
					"email": {
						Type:   "string",
						Format: "email",
					},
					"username": {
						Type: "string",
					},
					"usernameRegex": {
						Type: "string",
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
					"username": {
						Type: "string",
					},
				},
				Required: []string{"id"},
			},
			GetModal: addCommonProps(
				&openapi.Schema{
					Type: "object",
					Properties: map[string]openapi.Schema{
						"username": {
							Type: "string",
						},
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
