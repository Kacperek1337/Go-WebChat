package views

//RegisterBody model
type RegisterBody struct {
	Username string `form:"username" binding:"required,username"`
	Email    string `form:"email" biding:"required,email"`
	Password string `form:"password"`
}

//LoginBody struct
type LoginBody struct {
	Email    string `form:"email" biding:"required,email"`
	Password string `form:"password"`
}
