package roles

// Role constants
const (
	Admin     = "admin"
	Moderator = "moderator"
	User      = "user"
)

// GetAllRoles returns all valid roles
func GetAllRoles() []string {
	return []string{Admin, Moderator, User}
}

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	switch role {
	case Admin, Moderator, User:
		return true
	default:
		return false
	}
}

// HasAdminAccess checks if the role has admin-level access
func HasAdminAccess(role string) bool {
	return role == Admin
}

// HasModeratorAccess checks if the role has moderator-level access or higher
func HasModeratorAccess(role string) bool {
	return role == Admin || role == Moderator
}

// CanManageUsers checks if the role can manage users
func CanManageUsers(role string) bool {
	return role == Admin
}

// CanManageFiles checks if the role can manage any files
func CanManageFiles(role string) bool {
	return role == Admin || role == Moderator
}

// CanViewAllUsers checks if the role can view all users
func CanViewAllUsers(role string) bool {
	return role == Admin || role == Moderator
}