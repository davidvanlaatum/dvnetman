package modal

type ObjectType string

type Tag struct {
	Base        `bson:",inline"`
	Name        string       `bson:"name,omitempty"`
	Description string       `bson:"description,omitempty"`
	Color       string       `bson:"color,omitempty"`
	ObjectTypes []ObjectType `bson:"object_types,omitempty"`
}
