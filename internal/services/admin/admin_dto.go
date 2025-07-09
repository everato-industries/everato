package admin

// AdminLoginDTO represents the data transfer object for admin login.
//
// Either of the email or the username must be passed otherwise the validation
// will throw an error
type AdminLoginDTO struct {
	Email    string `json:"email" validate:"email"`
	UserName string `json:"username" validate:"min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}
