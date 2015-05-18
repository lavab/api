package env

// Flags contains values of flags which are important in the whole API
type Flags struct {
	BindAddress      string
	APIVersion       string
	LogFormatterType string
	ForceColors      bool
	EmailDomain      string

	SessionDuration int

	RedisAddress  string
	RedisDatabase int
	RedisPassword string

	RethinkDBAddress  string
	RethinkDBKey      string
	RethinkDBDatabase string

	LookupdAddress string
	NSQdAddress    string

	YubiCloudID  string
	YubiCloudKey string

	LogglyToken string

	SlackURL      string
	SlackLevels   string
	SlackChannel  string
	SlackIcon     string
	SlackUsername string

	BloomFilter string
	BloomCount  uint
}
