package entity

type Effect struct {
	ID      string `bson:"_id,omitempty" json:"id"`
	Slug    string `bson:"slug" json:"slug"`
	Name    string `bson:"name" json:"name"`
	DSPType string `bson:"dsp_type" json:"dsp_type"`
}
