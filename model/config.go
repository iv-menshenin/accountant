package model

import "time"

type (
	Config interface {
		Init() error
		StringConfig(name, cmd, env string, defaultValue, usage string) *string
		IntegerConfig(name, cmd, env string, defaultValue int64, usage string) *int64
		BooleanConfig(name, cmd, env string, usage string) *bool
		DurationConfig(name, cmd, env string, defaultValue time.Duration, usage string) *time.Duration
	}
)
