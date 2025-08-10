package admin

import (
	"errors"
	"net/http"
	"slices"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

func CreateAdmin(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Parse the request body into a CreateAdminDTO
	adminDTO := NewCreateAdminDTO()
	if err := wr.ParseBody(adminDTO); err != nil {
		logger.StdoutLogger.Error("Error parsing create admin request body", "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid request body",
				"error":   err.Error(),
			},
		)
		return
	}

	// Validate the DTO
	if err := adminDTO.Validate(); err != nil {
		logger.StdoutLogger.Error("Invalid admin data format", "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin data format",
				"error":   err.Error(),
			},
		)
		return
	}

	// Get the current admin ID from the request context to ensure we have proper authorization
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
		logger.StdoutLogger.Error("Error starting transaction for CreateAdmin", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error",
				"error":   "Failed to start database transaction",
			},
		)
		return
	}

	// Convert the admin ID to UUID
	adminUUID, err := utils.StringToUUID(currentAdminID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid admin UUID", "adminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format",
				"error":   err.Error(),
			},
		)
		return
	}

	// Verify that the current admin exists and has permission to create other admins
	currentAdmin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), adminUUID)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to retrieve admin data", "adminID", currentAdminID, "err", err.Error())
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication failed",
				"error":   "Admin account not found",
			},
		)
		return
	}

	// Check if the current admin has permission to create other admins
	// Only SUPER_ADMIN can create other admins, or admins with MANAGE_USERS permission
	hasPermission := false
	if currentAdmin.Role == "SUPER_ADMIN" {
		hasPermission = true
	} else {
		if slices.Contains(currentAdmin.Permissions, "MANAGE_USERS") {
			hasPermission = true
		}
	}

	if !hasPermission {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Admin lacks permission to create other admins", "adminID", currentAdminID)
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "Access denied",
				"error":   "Insufficient permissions to create admin accounts",
			},
		)
		return
	}

	// Ensure SUPER_ADMIN role is only assignable by existing SUPER_ADMINs
	if adminDTO.Role == ADMIN_ROLE_SUPERADMIN && currentAdmin.Role != "SUPER_ADMIN" {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Non-super admin attempting to create super admin",
			"adminID", currentAdminID,
			"adminRole", currentAdmin.Role)
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "Access denied",
				"error":   "Only SUPER_ADMIN users can create other SUPER_ADMIN accounts",
			},
		)
		return
	}

	// Hash the password for the new admin
	hashedPassword, err := utils.BcryptHash(adminDTO.Password)
	if err != nil {
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to hash password", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to create admin account",
				"error":   "Password hashing failed",
			},
		)
		return
	}

	// Convert string role to SuperUserRole enum
	var role repository.SuperUserRole
	switch adminDTO.Role {
	case ADMIN_ROLE_SUPERADMIN:
		role = repository.SuperUserRoleSUPERADMIN
	case ADMIN_ROLE_ADMIN:
		role = repository.SuperUserRoleADMIN
	case ADMIN_ROLE_EDITOR:
		role = repository.SuperUserRoleEDITOR
	default:
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Invalid role provided", "role", adminDTO.Role)
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid role",
				"error":   "The provided role is not recognized",
			},
		)
		return
	}

	// Convert string permissions to Permissions enum array
	var permissions []repository.Permissions
	for _, perm := range adminDTO.Permissions {
		switch perm {
		case ADMIN_PERMISSION_MANAGE_EVENTS:
			permissions = append(permissions, repository.PermissionsMANAGEEVENTS)
		case ADMIN_PERMISSION_CREATE_EVENT:
			permissions = append(permissions, repository.PermissionsCREATEEVENT)
		case ADMIN_PERMISSION_EDIT_EVENT:
			permissions = append(permissions, repository.PermissionsEDITEVENT)
		case ADMIN_PERMISSION_DELETE_EVENT:
			permissions = append(permissions, repository.PermissionsDELETEEVENT)
		case ADMIN_PERMISSION_VIEW_EVENT:
			permissions = append(permissions, repository.PermissionsVIEWEVENT)
		case ADMIN_PERMISSION_MANAGE_BOOKINGS:
			permissions = append(permissions, repository.PermissionsMANAGEBOOKINGS)
		case ADMIN_PERMISSION_CREATE_BOOKING:
			permissions = append(permissions, repository.PermissionsCREATEBOOKING)
		case ADMIN_PERMISSION_EDIT_BOOKING:
			permissions = append(permissions, repository.PermissionsEDITBOOKING)
		case ADMIN_PERMISSION_DELETE_BOOKING:
			permissions = append(permissions, repository.PermissionsDELETEBOOKING)
		case ADMIN_PERMISSION_VIEW_BOOKING:
			permissions = append(permissions, repository.PermissionsVIEWBOOKING)
		case ADMIN_PERMISSION_MANAGE_USERS:
			permissions = append(permissions, repository.PermissionsMANAGEUSERS)
		case ADMIN_PERMISSION_VIEW_REPORTS:
			permissions = append(permissions, repository.PermissionsVIEWREPORTS)
		default:
			tx.Rollback(wr.R.Context()) // Rollback the transaction
			logger.StdoutLogger.Error("Invalid permission provided", "permission", perm)
			wr.Status(http.StatusBadRequest).Json(
				utils.M{
					"message": "Invalid permission",
					"error":   "The provided permission is not recognized: " + perm,
				},
			)
			return
		}
	}

	permStrings := make([]string, len(permissions))
	for i, p := range permissions {
		permStrings[i] = string(p)
	}

	newAdmin, err := repo.WithTx(tx).CreateAdminIfNotExists(
		wr.R.Context(),
		repository.CreateAdminIfNotExistsParams{
			Column1: adminDTO.UserName,
			Column2: adminDTO.Name,
			Column3: adminDTO.Email,
			Column4: hashedPassword,
			Column5: role,
			Column6: permStrings, // pass as []string
		},
	)

	// Check for specific error where the admin already exists
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// This means the admin already exists
			tx.Rollback(wr.R.Context()) // Rollback the transaction
			logger.StdoutLogger.Error("Admin already exists", "email", adminDTO.Email, "username", adminDTO.UserName)
			wr.Status(http.StatusConflict).Json(
				utils.M{
					"message": "Admin already exists",
					"error":   "An admin with this email or username already exists",
				},
			)
			return
		}

		// Some other error occurred
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		logger.StdoutLogger.Error("Failed to create admin account", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to create admin account",
				"error":   err.Error(),
			},
		)
		return
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

	// Return success response with the created admin (excluding password)
	wr.Status(http.StatusCreated).Json(
		utils.M{
			"message": "Admin account created successfully",
			"admin": utils.M{
				"id":          newAdmin.ID.String(),
				"email":       newAdmin.Email,
				"username":    newAdmin.Username,
				"name":        newAdmin.Name,
				"role":        newAdmin.Role,
				"permissions": newAdmin.Permissions,
				"created_at":  newAdmin.CreatedAt,
				"updated_at":  newAdmin.UpdatedAt,
			},
		},
	)
}
