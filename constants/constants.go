package constants

const (
	Yaml               = "yaml"
	Gzip               = "gzip"
	Redis              = "redis"
	ReadDatabase       = "read-database"
	WriteDatabase      = "write-database"
	GoroutineThreshold = "goroutine-threshold"
	Kafka              = "kafka"
)

// Type CD of user
const (
	TypeCD_B2C = 1
	TypeCD_B2B = 2
)

// Status of user
const (
	StatusCD_Active      = 1
	StatusCD_Preapproval = 2
	StatusCD_Suspended   = 3
	StatusCD_Inactive    = 4
	StatusCD_Locked      = 5
	StatusCD_Joining     = 6
)

// Status of user poll
const (
	Pool_B2B_User          = 7
	Pool_B2B_Master        = 6
	Pool_B2B_Branch_master = 63
)

// Branch default
const Head_Office = "본사"

// SignInCount
const SignInCount = 0

const Branch_ID_Default = 21
const Organization_ID_Default = 17

const (
	UserIDPayload = "UserIDPayload"
)
