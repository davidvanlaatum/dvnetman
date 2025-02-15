package modal

type IPPrefix struct {
	ObjectId string `bson:"_id,omitempty"`
	ID       string `bson:"id,omitempty"`
	Prefix   string `bson:"prefix,omitempty"`
}
