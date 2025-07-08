package errcode

// Internal Errors
const (
	ErrInternalFailure       = "IIF001" // Unexpected internal error
	ErrServiceUnavailable    = "ISU001" // Downstream service not available
)

// Input/Validation Errors
const (
	ErrInvalidInput          = "IIP001" // Invalid parameters
	ErrMissingRequiredField  = "IIP002" // Required field missing
	ErrInvalidFormat         = "IIP003" // Incorrect format (e.g., email, UUID)
	ErrTooLarge              = "IIP004" // Payload too large
	ErrTooManyRequests       = "IIP005" // Rate limit exceeded
)

// Auth & Access Errors
const (
	ErrUnauthorized          = "IUA001" // Not logged in or token missing
	ErrTokenExpired          = "IUA002" // Token expired
	ErrInvalidToken          = "IUA003" // Token malformed or revoked
	ErrForbidden             = "IFB001" // No permission for resource
)

// Resource Errors
const (
	ErrNotFound              = "INF001" // Entity not found
	ErrAlreadyExists         = "ICF001" // Duplicate entity
	ErrConflict              = "ICF002" // State conflict
)

// Business Logic Errors
const (
	ErrUserNotVerified       = "BUS001" // Email not verified
	ErrAccountDisabled       = "BUS002" // Suspended account
	ErrInsufficientBalance   = "BUS003" // For billing or point systems
)

// Dependency Errors
const (
	ErrDependencyFailure     = "DEP001" // External service call failed
	ErrTimeout               = "DEP002" // Timeout calling downstream
)

