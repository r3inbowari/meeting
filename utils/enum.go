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

const (
	ContentFeedback = iota
)

const (
	RoomNotStart = 0 //未开始
	RoomCarryOn  = 1 //进行中
	RoomHasEnd   = 2 //已结束
)
