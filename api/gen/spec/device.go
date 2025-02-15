package spec

import "dvnetman/api/gen/openapi"

func init() {
	RegisterBuilder(BuildDevice)
}

func BuildDevice(api *OpenAPI) {
	api.AddObject(
		AddObjectOpts{
			Name: "device",
			Tags: []string{"Device"},
			SearchModal: addCommonProps(
				&openapi.Schema{
					Type: "object",
					Properties: map[string]openapi.Schema{
						"name": {
							Type: "string",
						},
						"description": {
							Type: "string",
						},
						"status": {
							Type: "string",
						},
						"site": {
							Ref: objectRef,
						},
						"location": {
							Ref: objectRef,
						},
						"rack": {
							Ref: objectRef,
						},
						"deviceType": {
							Ref: objectRef,
						},
					},
					Required: []string{"id"},
				},
			),
			GetModal: addCommonProps(
				&openapi.Schema{
					Type: "object",
					Properties: map[string]openapi.Schema{
						"name": {
							Type: "string",
						},
						"description": {
							Type: "string",
						},
						"status": {
							Type: "string",
						},
						"site": {
							Ref: objectRef,
						},
						"location": {
							Ref: objectRef,
						},
						"rack": {
							Ref: objectRef,
						},
						"rackFace": {
							Type: "string",
							Enum: []string{"front", "rear"},
						},
						"rackPosition": {
							Type: "number",
						},
						"deviceType": {
							Ref: objectRef,
						},
						"ports": {
							Type: "array",
							Items: &openapi.Schema{
								Ref: "#/components/schemas/DevicePort",
							},
						},
					},
				},
			),
			Additional: map[string]*openapi.Schema{
				"DevicePort": {
					Type: "object",
					Properties: map[string]openapi.Schema{
						"id": {
							Type:     "string",
							Format:   "uuid",
							ReadOnly: true,
						},
						"name": {
							Type: "string",
						},
					},
				},
			},
		},
	)
}
