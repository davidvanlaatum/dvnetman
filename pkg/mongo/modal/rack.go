package modal

type Rack struct {
	Base     `bson:",inline"`
	Name     string `bson:"name,omitempty"`
	Location *UUID  `bson:"location,omitempty"`
	Site     *UUID  `bson:"site,omitempty"`
	Height   int    `bson:"height,omitempty"`
	Width    int    `bson:"width,omitempty"`
}
