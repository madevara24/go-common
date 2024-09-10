package database

type Configuration struct {
	Key     string `json:"key" yaml:"key"` // Key is connection key identification
	Dsn     string `json:"dsn" yaml:"dsn"` // Dsn configuration using uri
	SqlxKey string
}

func NewConfiguration(dsn, sqlxKey string) *Configuration {
	return &Configuration{
		Dsn:     dsn,
		SqlxKey: sqlxKey,
	}
}
