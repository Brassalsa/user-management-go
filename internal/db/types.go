package db

type User struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
	Avatar   string `json:"avatar" bson:"avatar"`
	Role     string `json:"role" bson:"role"`
}
