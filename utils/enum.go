package utils

const (
	Male = iota
	Famale
)

const (
	RoleDrift   = "drift"
	RoleManager = "manager"
	RoleUser    = "user"
)

const (
	StatusUndefined = iota
	StatusAudit
	StatusNormal
	StatusLock
	StatusDenied
)

const (
	OpSucceed = iota
	OpFailed
	OpPasswordError
	OpLoginError
	OpLogonError
	OpJsonBindError
	OpValidateError
	OpAuthError
	OpResourcesDenied
)
