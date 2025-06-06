package entity

type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Name     string `bson:"name" json:"name"`
	Login    string `bson:"login" json:"login"`
	Password string `bson:"password,omitempty" json:"-"`
}
