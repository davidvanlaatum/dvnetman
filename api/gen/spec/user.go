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
			SearchModal: &openapi.Schema{
				Type: "object",
				Properties: map[string]openapi.Schema{
					"id": {
						Type:   "string",
						Format: "uuid",
					},
					"displayName": {
						Type: "string",
					},
				},
			},
			GetModal: &openapi.Schema{
				Type: "object",
				Properties: map[string]openapi.Schema{
					"id": {
						Type:   "string",
						Format: "uuid",
					},
					"username": {
						Type: "string",
					},
					"email": {
						Type: "string",
					},
					"password": {
						Type: "string",
					},
					"externalProvider": {
						Type: "string",
					},
				},
			},
		},
	)
}
