// Package user provides authentication and user management services for the Everato platform.
// It handles user creation, verification, authentication, and profile management.
package user

import (
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserDTO represents the data transfer object for creating a new user.
// It defines the structure and validation rules for user registration data.
type CreateUserDTO struct {
	FistName string `json:"first_name" validate:"required,min=2,max=50"` // User's first name (2-50 chars)
	LastName string `json:"last_name" validate:"required,min=2,max=50"`  // User's last name (2-50 chars)
	Email    string `json:"email" validate:"required,email"`             // User's email address (must be valid format)
	Password string `json:"password" validate:"required,min=8,max=100"`  // User's password (8-100 chars)

	AdminUsername *string `json:"admin_username" validate:"omitempty,min=2,max=50"` // Optional admin username for user creation
	AdminEmail    *string `json:"admin_email" validate:"omitempty,email"`           // Optional admin email for user creation
	AdminPassword *string `json:"admin_password" validate:"omitempty,min=8"`        // Optional admin password for user creation
}

// LoginUserDTO represents the data transfer object for user login.
// It defines the structure and validation rules for authentication credentials.
type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email"`    // User's email address for login
	Password string `json:"password" validate:"required,min=8"` // User's password (minimum 8 chars)
}

// Validate checks the CreateUserDTO for required fields and formats.
// It ensures all user registration data meets the defined validation rules.
//
// The validation ensures:
// - All required fields are present
// - Name fields meet minimum and maximum length requirements
// - Email is in valid format
// - Password meets minimum security requirements
//
// Returns:
//   - nil if validation passes
//   - A validation error detailing which fields failed validation and why
func (c CreateUserDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	return v.Struct(c)
}

// HashPassword securely hashes the user's password using bcrypt.
// This method modifies the DTO by replacing the plaintext password with its hashed version.
//
// The hashing process:
// 1. Uses bcrypt with the default cost factor for security
// 2. Converts the plaintext password to a secure hash
// 3. Replaces the plaintext password in the DTO with the generated hash
// 4. Logs any errors that occur during hashing
//
// Returns:
//   - nil if hashing succeeds
//   - An error if hashing fails
func (c *CreateUserDTO) HashPassword() error {
	logger := pkg.NewLogger()
	defer logger.Close() // Ensure the logger is closed when the function exits

	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost) // Generate the hash of the password
	if err != nil {
		logger.StdoutLogger.Error("Error hashing password", "err", err.Error()) // Log the error to stdout
		logger.FileLogger.Error("Error hashing password", "err", err.Error())   // Log the error to file
		return err
	}

	c.Password = string(hashed_pass) // Set the hashed password to the actual struct
	return nil
}

// ToCreteUserParams converts the CreateUserDTO to a CreateUserParams structure
// for database operations. This method transforms the validated DTO into the format
// required by the repository layer.
//
// Returns:
//   - A repository.CreateUserParams structure ready for database insertion
func (c CreateUserDTO) ToCreteUserParams() repository.CreateUserParams {
	return repository.CreateUserParams{
		FirstName: c.FistName,
		LastName:  c.LastName,
		Email:     c.Email,
		Password:  c.Password, // At this point, Password should already be hashed
	}
}

// Validate checks the LoginUserDTO for required fields and formats.
// It ensures the login credentials meet the defined validation rules.
//
// The validation ensures:
// - Email is provided and in valid format
// - Password is provided and meets minimum length requirement
//
// Returns:
//   - nil if validation passes
//   - A validation error detailing which fields failed validation and why
func (l LoginUserDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	return v.Struct(l)
}

// VerifyPassword securely checks if the provided plaintext password matches a stored hash.
// It uses bcrypt's secure comparison function to prevent timing attacks.
//
// Parameters:
//   - hashedPassword: The bcrypt hash stored in the database
//
// Returns:
//   - nil if passwords match
//   - An error if passwords don't match or comparison fails
func (l LoginUserDTO) VerifyPassword(hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(l.Password))
}
