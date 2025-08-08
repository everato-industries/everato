package admin

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// Enum values for the super user permissions
const (
	ADMIN_PERMISSION_MANAGE_EVENTS   = "MANAGE_EVENTS"
	ADMIN_PERMISSION_CREATE_EVENT    = "CREATE_EVENT"
	ADMIN_PERMISSION_EDIT_EVENT      = "EDIT_EVENT"
	ADMIN_PERMISSION_DELETE_EVENT    = "DELETE_EVENT"
	ADMIN_PERMISSION_VIEW_EVENT      = "VIEW_EVENT"
	ADMIN_PERMISSION_MANAGE_BOOKINGS = "MANAGE_BOOKINGS"
	ADMIN_PERMISSION_CREATE_BOOKING  = "CREATE_BOOKING"
	ADMIN_PERMISSION_EDIT_BOOKING    = "EDIT_BOOKING"
	ADMIN_PERMISSION_DELETE_BOOKING  = "DELETE_BOOKING"
	ADMIN_PERMISSION_VIEW_BOOKING    = "VIEW_BOOKING"
	ADMIN_PERMISSION_MANAGE_USERS    = "MANAGE_USERS"
	ADMIN_PERMISSION_VIEW_REPORTS    = "VIEW_REPORTS"
)

// Enum values for the super user role
const (
	ADMIN_ROLE_SUPERADMIN = "SUPER_ADMIN" // These users have 100% access to the dashboard and can override anything
	ADMIN_ROLE_ADMIN      = "ADMIN"       // These users have limited access to the dashboard and can only manage events and bookings
	ADMIN_ROLE_EDITOR     = "EDITOR"      // These users can only edit events and bookings, but cannot create or delete them
)

// AdminLoginDTO represents the data transfer object for admin login.
//
// Either of the email or the username must be passed otherwise the validation
// will throw an error
type AdminLoginDTO struct {
	Email    string `json:"email" validate:"email"`                     // Super user's email
	UserName string `json:"username" validate:"max=20"`                 // Super user's username
	Password string `json:"password" validate:"required,min=8,max=100"` // Super user's password
}

// Initialize a new instance of AdminLoginDTO
func NewAdminLoginDTO() *AdminLoginDTO {
	return &AdminLoginDTO{
		Email:    "",
		UserName: "",
		Password: "",
	}
}

// Validate checks the AdminLoginDTO for required fields and formats.
func (a *AdminLoginDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	// At least one of the email or username must be provided
	if a.Email == "" && a.UserName == "" {
		return errors.New("either of the username or the email must be provided")
	}

	// Validate the struct
	return v.Struct(a)
}

// Hashpassword method hashes the password with bcrypt algorithm
func (a *AdminLoginDTO) HashPassword() error {
	if a.Password == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err // Return the error if hashing fails
	}

	a.Password = string(hashedPassword) // Replace the plaintext password with the hashed one
	return nil
}

func (a *AdminLoginDTO) VerifyPassword(hashedPassword string) (bool, error) {
	if a.Password == "" {
		return false, errors.New("password cannot be empty")
	}

	if hashedPassword == "" {
		return false, errors.New("hashed password cannot be empty")
	}

	// Compare the hashed password with the provided password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(a.Password))
	if err != nil {
		return false, err // Return the error if comparison fails
	}

	return true, nil // Password matches successfully
}

// CreateAdminDTO represents the DTO for creating anothersub-admin account
//
// All of the fields are necessary
type CreateAdminDTO struct {
	Email       string   `json:"email" validate:"required,email"`                        // Admin's email address
	UserName    string   `json:"username" validate:"required,min=3,max=20"`              // Admin's username
	Name        string   `json:"name" validate:"required,min=2,max=100"`                 // Admin's name
	Password    string   `json:"password" validate:"required,min=8,max=100"`             // Admin's password
	Role        string   `json:"role" validate:"required,role"`                          // Admin's role (superadmin or admin)
	Permissions []string `json:"permissions" validate:"required,min=1,dive,permissions"` // Admin's permissions
}

// Initialize a new instance of CreateAdminDTO
func NewCreateAdminDTO() *CreateAdminDTO {
	return &CreateAdminDTO{
		Email:       "",
		Name:        "",
		UserName:    "",
		Password:    "",
		Role:        ADMIN_ROLE_EDITOR, // Default role is ADMIN
		Permissions: []string{},        // Default permissions is an empty slice
	}
}

// Permission validator validates the permissions field in CreateAdminDTO
//
// It must be a non-empty slice with at least one permission
// and all of the permissions are follows
//   - MANAGE_EVENTS
//   - CREATE_EVENT
//   - EDIT_EVENT
//   - DELETE_EVENT
//   - VIEW_EVENT
//   - MANAGE_BOOKINGS
//   - CREATE_BOOKING
//   - EDIT_BOOKING
//   - DELETE_BOOKING
//   - VIEW_BOOKING
//   - MANAGE_USERS
//   - VIEW_REPORTS
func permission_validator(fl validator.FieldLevel) bool {
	// First, get the field value as interface{}
	fieldValue := fl.Field().Interface()

	// Initialize an empty slice to store permissions
	var permissions []string

	// Handle different types using type assertions
	switch v := fieldValue.(type) {
	case string:
		// Single string case
		permissions = []string{v}
	case []string:
		// Direct []string case from CreateAdminDTO
		permissions = v
	case *[]string:
		// Pointer case from UpdateAdminDTO
		if v == nil {
			return true // nil pointer is valid for "omitempty"
		}
		permissions = *v
	default:
		// Unsupported type
		return false
	}

	// Handle empty permissions list
	if len(permissions) == 0 {
		return true // Empty is allowed with "omitempty"
	}

	// Define valid permissions
	validPermissions := map[string]bool{
		ADMIN_PERMISSION_MANAGE_EVENTS:   true,
		ADMIN_PERMISSION_CREATE_EVENT:    true,
		ADMIN_PERMISSION_EDIT_EVENT:      true,
		ADMIN_PERMISSION_DELETE_EVENT:    true,
		ADMIN_PERMISSION_VIEW_EVENT:      true,
		ADMIN_PERMISSION_MANAGE_BOOKINGS: true,
		ADMIN_PERMISSION_CREATE_BOOKING:  true,
		ADMIN_PERMISSION_EDIT_BOOKING:    true,
		ADMIN_PERMISSION_DELETE_BOOKING:  true,
		ADMIN_PERMISSION_VIEW_BOOKING:    true,
		ADMIN_PERMISSION_MANAGE_USERS:    true,
		ADMIN_PERMISSION_VIEW_REPORTS:    true,
	}

	// Check each permission
	for _, permission := range permissions {
		if !validPermissions[permission] {
			return false // Invalid permission found
		}
	}

	return true // All permissions are valid
}

// Role validator validates the role field in CreateAdminDTO
//
// It must be a non empty string, of either of the followings
//   - SUPER_USER
//   - ADMIN
//   - EDITOR
func role_validator(fl validator.FieldLevel) bool {
	// First, get the field value as interface{}
	fieldValue := fl.Field().Interface()

	// Initialize a role variable
	var role string

	// Handle different types using type assertions
	switch v := fieldValue.(type) {
	case string:
		// Direct string case from CreateAdminDTO
		role = v
	case *string:
		// Pointer case from UpdateAdminDTO
		if v == nil {
			return true // nil pointer is valid for "omitempty"
		}
		role = *v
	default:
		// Unsupported type
		return false
	}

	// Handle empty role
	if role == "" {
		return true // Empty is allowed with "omitempty"
	}

	// Define valid roles
	validRoles := map[string]bool{
		ADMIN_ROLE_SUPERADMIN: true,
		ADMIN_ROLE_ADMIN:      true,
		ADMIN_ROLE_EDITOR:     true,
	}

	return validRoles[role] // Returns true if the role is valid, false otherwise
}

// Validate checks the AdminLoginDTO for required fields and formats.
//
// The validation ensures:
// // - Email is in valid format
// // - Username meets minimum and maximum length requirements
// // - Password meets minimum security requirements
// // Returns:
// //   - nil if validation passes
// //   - A validation error detailing which fields failed validation and why
func (a *CreateAdminDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("permissions", permission_validator)
	v.RegisterValidation("role", role_validator)

	return v.Struct(a)
}
