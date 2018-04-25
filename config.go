package influxdb

import (
	"errors"
	"time"

	"github.com/devopsfaith/krakend/config"
)

const defaultBufferSize = 0

type influxConfig struct {
	address    string
	username   string
	password   string
	ttl        time.Duration
	database   string
	bufferSize int
}

func configGetter(extraConfig config.ExtraConfig) interface{} {
	value, ok := extraConfig[Namespace]
	if !ok {
		return nil
	}

	castedConfig, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}

	cfg := influxConfig{}

	if value, ok := castedConfig["address"]; ok {
		cfg.address = value.(string)
	}

	if value, ok := castedConfig["username"]; ok {
		cfg.username = value.(string)
	}

	if value, ok := castedConfig["password"]; ok {
		cfg.password = value.(string)
	}

	if value, ok := castedConfig["buffer_size"]; ok {
		if s, ok := value.(int); ok {
			cfg.bufferSize = s
		}
	}

	if value, ok := castedConfig["ttl"]; ok {
		s, ok := value.(string)

		if !ok {
			return nil
		}
		var err error
		cfg.ttl, err = time.ParseDuration(s)

		if err != nil {
			return nil
		}
	}

	if value, ok := castedConfig["db"]; ok {
		cfg.database = value.(string)
	} else {
		cfg.database = "krakend"
	}

	return cfg
}

var errNoConfig = errors.New("Unable to load custom config from the extra config")
