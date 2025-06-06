package entity

type Effect struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	DSPType     string `bson:"dsp_type" json:"dsp_type"`
}
