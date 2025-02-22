package lksh

import (
	"context"
	"time"
)

const (
	MegaByte    = 1024 * 1024
	DefaultExit = 255
)

type Config struct {
	Color         bool
	Ctx           context.Context
	Decorate      bool
	MaxBufferSize int64
	KeepEnvVar    []string
	LookUpPath    bool
	Timeout       time.Duration
	Pipe          bool
}

func (c *Config) AddEnvVar(v string) *Config {
	c.KeepEnvVar = append(c.KeepEnvVar, v)
	return c
}

func DefaultConfig() *Config {
	return &Config{
		LookUpPath:    true,
		Decorate:      true,
		KeepEnvVar:    []string{"PATH", "CXX", "CXXFALGS", "C", "CFLAGS", "V", "VERBOSE"},
		MaxBufferSize: 100 * MegaByte,
		Pipe:          true,
	}
}
