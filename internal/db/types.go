package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserRegister struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Name     string `json:"name" bson:"name"`
	Password string `bson:"password"`
	Avatar   string `json:"avatar" bson:"avatar"`
	Role     string `json:"role" bson:"role"`
}

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	UserRegister `json:",inline" bson:",inline"`
}

type UserWithoutPassword struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Name     string             `json:"name" bson:"name"`
	Avatar   string             `json:"avatar" bson:"avatar"`
	Role     string             `json:"role" bson:"role"`
}
