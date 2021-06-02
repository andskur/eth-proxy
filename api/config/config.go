package config

// Scheme represent top-level
// microservice Application config scheme
type Scheme struct {
	Log  *Log
	Api  *Addr
	Http *Addr
}

// Log represent Application logger params
type Log struct {
	// Level is logger available levels:
	// Trace, Debug, Info, Warning, Error,
	// Fatal and Panic (default Info)
	Level string
}

// Addr represent standard address
// structure with host and port
type Addr struct {
	Host    string
	Port    int
	Timeout string
}
