package responses

var (
	PathUsername = "data.username"
	PathPassword = "data.password"
	PathEmail    = "data.email"

	UsernameExists   = "username already exists"
	UsernameTooShort = "username too short"
	EmailExists      = "email address already exists"
	EmailIncorrect   = "email address incorrect"
	PasswordTooShort = "password too short"

	WrongUsernamePassword = "wrong username or password"
	UserInactive          = "user is inactive"

	InsufficientPermissionOrNonExistentResource = "user has insufficient permissions for this action or the resource does not exist"
)
