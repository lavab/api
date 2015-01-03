package env

// Flags contains values of flags which are important in the whole API
type Flags struct {
	BindAddress      string
	APIVersion       string
	LogFormatterType string
	ForceColors      bool

	SessionDuration     int
	ClassicRegistration bool
	UsernameReservation bool

	RedisAddress  string
	RedisDatabase int
	RedisPassword string

	RethinkDBAddress  string
	RethinkDBKey      string
	RethinkDBDatabase string

	NATSAddress string

	YubiCloudID  string
	YubiCloudKey string

	LogglyToken string
}
