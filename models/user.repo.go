package models

type User struct {
	Id         string `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	UserName   string `json:"user_name" bson:"user_name"`
	Password   string `json:"password" bson:"password"`
	IsAdmin    bool   `json:"is_admin" bson:"is_admin"`
	IsVerified bool   `json:"is_verified" bson:"is_verified"`
}

type UserResp struct {
	Id         string `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	UserName   string `json:"user_name" bson:"user_name"`
	IsAdmin    bool   `json:"is_admin" bson:"is_admin"`
	IsVerified bool   `json:"is_verified" bson:"is_verified"`
}

type AddUser struct {
	Name     string `json:"name" bson:"name" binding:"required"`
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type Login struct {
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}
