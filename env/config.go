package env

// Flags contains values of flags which are important in the whole API
type Flags struct {
	BindAddress         string
	APIVersion          string
	LogFormatterType    string
	SessionDuration     int
	ClassicRegistration bool
}
