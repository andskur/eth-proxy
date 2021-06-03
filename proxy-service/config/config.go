package config

// Scheme represent top-level
// microservice Application config scheme
type Scheme struct {
	Log      *Log
	Grpc     *Grpc
	Ethereum *Ethereum
}

// Log represent Application logger params
type Log struct {
	// Level is logger available levels:
	// Trace, Debug, Info, Warning, Error,
	// Fatal and Panic (default Info)
	Level string
}

// Ethereum represent Ethereum
// node connection params
type Ethereum struct {
	Addr string
	Wss  bool
}

// Grpc represent Grpc server params
type Grpc struct {
	Addr    `mapstructure:",squash"`
	Timeout string
}

// Addr represent standard address
// structure with host and port
type Addr struct {
	Host string
	Port int
}
