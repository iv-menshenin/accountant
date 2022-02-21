package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	ConfigStorage struct {
		prefix     string
		arguments  []string
		flagSet    *flag.FlagSet
		registered []config
	}
	config struct {
		name    string
		cmdName string
		envName string
		claim   claim
	}
	claim interface {
		initEnv(string) (bool, error)
		initCmd()
	}
	strClaim struct {
		dest *string
		cmd  *string
	}
	intClaim struct {
		dest *int64
		cmd  *int64
	}
	boolClaim struct {
		dest *bool
		cmd  *bool
	}
	durClaim struct {
		dest *time.Duration
		cmd  *time.Duration
	}
)

func New(prefix string, args ...string) *ConfigStorage {
	flagSet := flag.NewFlagSet(prefix, flag.ContinueOnError)
	if len(args) == 0 {
		args = os.Args[1:]
	}
	return &ConfigStorage{
		prefix:    prefix,
		arguments: args,
		flagSet:   flagSet,
	}
}

func (c *ConfigStorage) StringConfig(name, cmd, env string, defaultValue, usage string) *string {
	var dest = new(string)
	var cnf = c.makeConfig(name, cmd, env)
	cnf.claim = &strClaim{
		dest: dest,
		cmd:  c.flagSet.String(cnf.cmdName, defaultValue, usage),
	}
	c.registered = append(c.registered, cnf)
	return dest
}

func (c *ConfigStorage) IntegerConfig(name, cmd, env string, defaultValue int64, usage string) *int64 {
	var dest = new(int64)
	var cnf = c.makeConfig(name, cmd, env)
	cnf.claim = &intClaim{
		dest: dest,
		cmd:  c.flagSet.Int64(cnf.cmdName, defaultValue, usage),
	}
	c.registered = append(c.registered, cnf)
	return dest
}

func (c *ConfigStorage) BooleanConfig(name, cmd, env string, usage string) *bool {
	var dest = new(bool)
	var cnf = c.makeConfig(name, cmd, env)
	cnf.claim = &boolClaim{
		dest: dest,
		cmd:  c.flagSet.Bool(cnf.cmdName, false, usage),
	}
	c.registered = append(c.registered, cnf)
	return dest
}

func (c *ConfigStorage) DurationConfig(name, cmd, env string, defaultValue time.Duration, usage string) *time.Duration {
	var dest = new(time.Duration)
	var cnf = c.makeConfig(name, cmd, env)
	cnf.claim = &durClaim{
		dest: dest,
		cmd:  c.flagSet.Duration(cnf.cmdName, defaultValue, usage),
	}
	c.registered = append(c.registered, cnf)
	return dest
}

func (c *ConfigStorage) makeConfig(name, cmd, env string) config {
	if c.prefix != "" {
		cmd = fmt.Sprintf("%s.%s", c.prefix, cmd)
		env = strings.ToUpper(fmt.Sprintf("%s_%s", c.prefix, env))
	}
	return config{
		name:    name,
		cmdName: cmd,
		envName: env,
	}
}

func (c *ConfigStorage) Init() (err error) {
	if err = c.flagSet.Parse(c.arguments); err != nil {
		return err
	}
	for i := range c.registered {
		var (
			ok  bool
			cnf = c.registered[i]
		)
		if ok, err = cnf.claim.initEnv(cnf.envName); err != nil {
			return err
		}
		if !ok {
			cnf.claim.initCmd()
		}
	}
	return nil
}

func (s *strClaim) initEnv(name string) (bool, error) {
	env, ok := os.LookupEnv(name)
	if ok {
		*s.dest = env
	}
	return ok, nil
}

func (s *strClaim) initCmd() {
	*s.dest = *s.cmd
}

func (s *intClaim) initEnv(name string) (bool, error) {
	if env, ok := os.LookupEnv(name); ok {
		i, err := strconv.ParseInt(env, 10, 64)
		if err != nil {
			return false, err
		}
		*s.dest = i
		return true, nil
	}
	return false, nil
}

func (s *intClaim) initCmd() {
	*s.dest = *s.cmd
}

func (s *boolClaim) initEnv(name string) (bool, error) {
	if env, ok := os.LookupEnv(name); ok {
		i, err := strconv.ParseBool(env)
		if err != nil {
			return false, err
		}
		*s.dest = i
		return true, nil
	}
	return false, nil
}

func (s *boolClaim) initCmd() {
	*s.dest = *s.cmd
}

func (s *durClaim) initEnv(name string) (bool, error) {
	if env, ok := os.LookupEnv(name); ok {
		dur, err := time.ParseDuration(env)
		if err != nil {
			return false, err
		}
		*s.dest = dur
		return true, nil
	}
	return false, nil
}

func (s *durClaim) initCmd() {
	*s.dest = *s.cmd
}
