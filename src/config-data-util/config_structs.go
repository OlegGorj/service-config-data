// Package config-data-util defines all of the data structures used in the service

package config_data_util

import (
	"config-data-util/environment"
)

// An Environment represents the configuration of a particular environment this can be one of
// {sandbox, dev}

type MappingToEnv map[string]*environment.Environment
