package lksh

const (
	MegaByte = 1024 * 1024
)

type Config struct {
	Pipe          bool
	MaxBufferSize int64
	Decorate      bool

	LookUpPath  bool
	KeepEnvVars []string
	AddEnvVars  map[string]string
	DefaultExit int
}

func NewConfig() *Config {
	return &Config{
		Pipe:          true,
		MaxBufferSize: 0,
		Decorate:      true,

		LookUpPath:  true,
		KeepEnvVars: []string{},
		AddEnvVars:  map[string]string{},
		DefaultExit: 255,
	}
}

func (c *Config) KeepEnvVar(v string) *Config {
	c.KeepEnvVars = append(c.KeepEnvVars, v)
	return c
}

func (c *Config) AddEnvVar(k, v string) *Config {
	c.AddEnvVars[k] = v
	return c
}
