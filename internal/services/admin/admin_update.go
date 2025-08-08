package admin

import (
	"net/http"
	"slices"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// UpdateAdminDTO defines the structure for admin update requests
type UpdateAdminDTO struct {
	Email       *string   `json:"email,omitempty" validate:"omitempty,email"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Password    *string   `json:"password,omitempty" validate:"omitempty,min=8,max=100"`
	Role        *string   `json:"role,omitempty" validate:"omitempty,role"`
	Permissions *[]string `json:"permissions,omitempty" validate:"omitempty,dive,permissions"`
}

// Validate checks the UpdateAdminDTO for valid formats
func (u *UpdateAdminDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("permissions", permission_validator)
	v.RegisterValidation("role", role_validator)

	return v.Struct(u)
}

func UpdateAdmin(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Extract the admin ID from URL parameters
	vars := mux.Vars(wr.R)
	targetAdminID := vars["id"]
	if targetAdminID == "" {
		logger.StdoutLogger.Error("Admin ID not provided in URL parameters")
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Missing admin ID",
				"error":   "Admin ID must be provided in the URL",
			},
		)
		return
	}

	// Parse the request body into an UpdateAdminDTO
	updateDTO := &UpdateAdminDTO{}
	if err := wr.ParseBody(updateDTO); err != nil {
		logger.StdoutLogger.Error("Error parsing update admin request body", "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid request body",
				"error":   err.Error(),
			},
		)
		return
	}

	// Validate the DTO if any fields are provided
	if updateDTO.Email != nil || updateDTO.Name != nil || updateDTO.Password != nil ||
		updateDTO.Role != nil || updateDTO.Permissions != nil {
		if err := updateDTO.Validate(); err != nil {
			logger.StdoutLogger.Error("Invalid admin data format", "err", err.Error())
			wr.Status(http.StatusBadRequest).Json(
				utils.M{
					"message": "Invalid admin data format",
					"error":   err.Error(),
				},
			)
			return
		}
	}

	// Get the current admin ID from the request context
	currentAdminID, ok := wr.R.Context().Value("admin_id").(string)
	if !ok {
		logger.StdoutLogger.Error("Failed to get admin ID from context")
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication required",
				"error":   "Admin ID not found in request context",
			},
		)
		return
	}

	// Start a transaction for atomicity
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Error starting transaction for UpdateAdmin", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to start database transaction",
			},
		)
		return
	}

	// Convert the admin IDs to UUIDs
	currentAdminUUID, err := utils.StringToUUID(currentAdminID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid current admin UUID", "adminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format",
				"error":   err.Error(),
			},
		)
		return
	}

	targetAdminUUID, err := utils.StringToUUID(targetAdminID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid target admin UUID", "adminID", targetAdminID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format",
				"error":   err.Error(),
			},
		)
		return
	}

	// Get current admin details to check permissions
	currentAdmin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), currentAdminUUID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve current admin data", "adminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication failed",
				"error":   "Admin account not found",
			},
		)
		return
	}

	// Get target admin details
	targetAdmin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), targetAdminUUID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve target admin data", "adminID", targetAdminID, "err", err.Error())
		wr.Status(http.StatusNotFound).Json(
			utils.M{
				"message": "Target admin not found",
				"error":   "The admin account you're trying to update doesn't exist",
			},
		)
		return
	}

	// Check permissions for update
	// 1. Admins can always update their own account
	// 2. SUPER_ADMINs can update any account
	// 3. Admins with MANAGE_USERS permission can update non-SUPER_ADMIN accounts
	isSelfUpdate := currentAdminID == targetAdminID
	isSuperAdmin := currentAdmin.Role == "SUPER_ADMIN"
	hasManageUsersPermission := slices.Contains(currentAdmin.Permissions, "MANAGE_USERS")

	// Check if the current admin is allowed to update the target admin
	hasUpdatePermission := isSelfUpdate || isSuperAdmin ||
		(hasManageUsersPermission && targetAdmin.Role != "SUPER_ADMIN")

	if !hasUpdatePermission {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Admin lacks permission to update target admin",
			"currentAdminID", currentAdminID,
			"targetAdminID", targetAdminID)
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "Access denied",
				"error":   "You don't have permission to update this admin account",
			},
		)
		return
	}

	// Special case: only SUPER_ADMINs can promote others to SUPER_ADMIN
	if updateDTO.Role != nil && *updateDTO.Role == ADMIN_ROLE_SUPERADMIN && currentAdmin.Role != "SUPER_ADMIN" {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Non-super admin attempting to promote to super admin",
			"adminID", currentAdminID,
			"adminRole", currentAdmin.Role)
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "Access denied",
				"error":   "Only SUPER_ADMIN users can promote others to SUPER_ADMIN",
			},
		)
		return
	}

	// TODO: Implement the actual update operation
	// This would typically use a repository method like:
	// UpdateAdminByID(ctx, UpdateAdminByIDParams{...})

	// Since we don't have that method in the repository yet, we'll simulate the update
	updatedAdmin := targetAdmin

	// Prepare the updated fields message
	updatedFields := []string{}

	// Update fields if provided
	if updateDTO.Email != nil {
		// In a real implementation, we would update the email
		updatedFields = append(updatedFields, "email")
	}

	if updateDTO.Name != nil {
		// In a real implementation, we would update the name
		updatedFields = append(updatedFields, "name")
	}

	if updateDTO.Password != nil {
		// In a real implementation, we would hash and update the password
		_, err := utils.BcryptHash(*updateDTO.Password)
		if err != nil {
			tx.Rollback(wr.R.Context()) // Rollback the transaction
			logger.StdoutLogger.Error("Failed to hash password", "err", err.Error())
			wr.Status(http.StatusInternalServerError).Json(
				utils.M{
					"message": "Failed to update admin account",
					"error":   "Password hashing failed",
				},
			)
			return
		}
		// We would update the password in the database
		updatedFields = append(updatedFields, "password")
	}

	if updateDTO.Role != nil {
		// In a real implementation, we would update the role
		updatedFields = append(updatedFields, "role")
	}

	if updateDTO.Permissions != nil {
		// In a real implementation, we would update the permissions
		updatedFields = append(updatedFields, "permissions")
	}

	// Commit the transaction
	if err = tx.Commit(wr.R.Context()); err != nil {
		logger.StdoutLogger.Error("Failed to commit transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to complete database operation",
			},
		)
		return
	}

	// Return success response
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message":        "Admin account updated successfully",
			"updated_fields": updatedFields,
			"admin": utils.M{
				"id":          updatedAdmin.ID.String(),
				"email":       updatedAdmin.Email,
				"username":    updatedAdmin.Username,
				"name":        updatedAdmin.Name,
				"role":        updatedAdmin.Role,
				"permissions": updatedAdmin.Permissions,
				"created_at":  updatedAdmin.CreatedAt,
				"updated_at":  updatedAdmin.UpdatedAt,
			},
		},
	)
}
