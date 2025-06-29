package user

import (
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserDTO represents the data transfer object for creating a new user
type CreateUserDTO struct {
	FistName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName string `json:"last_name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// LoginUserDTO represents login credentials
type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Validate checks the CreateUserDTO for required fields and formats.
//
// Parameters:
//   - json_body - This is a pointer to an Unmarshalled version of the raw request body which you want to validate,
//     which in this case this is the raw user_json value
func (c CreateUserDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	return v.Struct(c)
}

// HashPassword hashes the user's password. This is a placeholder function.
//
// This will automatically override the original password with the hashed one
func (c *CreateUserDTO) HashPassword() error {
	logger := pkg.NewLogger()
	defer logger.Close() // Ensure the logger is closed when the function exits

	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost) // Generate the hash of the password
	if err != nil {
		logger.StdoutLogger.Error("Error hashing password", "err", err.Error()) // Log the error to stdou
		logger.FileLogger.Error("Error hashing password", "err", err.Error())   // Log the error to file
		return err
	}

	c.Password = string(hashed_pass) // Set the hased password to the actual struct
	return nil
}

// ToCreteUserParams converts the CreateUserDTO to a CreateUserParams for database operations.
func (c CreateUserDTO) ToCreteUserParams() repository.CreateUserParams {
	return repository.CreateUserParams{
		FirstName: c.FistName,
		LastName:  c.LastName,
		Email:     c.Email,
		Password:  c.Password,
	}
}

// Validate checks the LoginUserDTO for required fields and formats
func (l LoginUserDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	return v.Struct(l)
}

// VerifyPassword checks if the provided password matches the hashed password
func (l LoginUserDTO) VerifyPassword(hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(l.Password))
}
