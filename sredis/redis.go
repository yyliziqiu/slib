package sredis

import (
	"github.com/go-redis/redis/v8"
)

var (
	_configs map[string]Config

	_clis map[string]*redis.Client
	_clus map[string]*redis.ClusterClient
)

func Init(configs ...Config) error {
	_configs = make(map[string]Config, 16)
	for _, config := range configs {
		conf := config.Default()
		_configs[conf.Id] = conf
	}

	_clis = make(map[string]*redis.Client, 16)
	_clus = make(map[string]*redis.ClusterClient, 16)
	for _, config := range configs {
		cli, clu := New(config)
		if cli != nil {
			_clis[config.Id] = cli
		}
		if clu != nil {
			_clus[config.Id] = clu
		}
	}

	return nil
}

func New(config Config) (*redis.Client, *redis.ClusterClient) {
	switch config.Mode {
	case ModeCluster:
		return nil, NewCluster(config)
	case ModeSentinel:
		return NewFailover(config), nil
	case ModeSentinelCluster:
		return nil, NewFailoverCluster(config)
	default:
		return NewSingle(config), nil
	}
}

func NewSingle(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:               config.Addr,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.Db,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	})
}

func NewCluster(config Config) *redis.ClusterClient {
	ops := &redis.ClusterOptions{
		Addrs:              config.Addrs,
		Username:           config.Username,
		Password:           config.Password,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}

	switch config.ReadPreference {
	case "ReadOnly":
		ops.ReadOnly = true
	case "RouteByLatency":
		ops.RouteByLatency = true
	case "RouteRandomly":
		ops.RouteRandomly = true
	}

	return redis.NewClusterClient(ops)
}

func NewFailover(config Config) *redis.Client {
	ops := &redis.FailoverOptions{
		MasterName:         config.MasterName,
		SentinelAddrs:      config.SentinelAddrs,
		SentinelPassword:   config.SentinelPassword,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.Db,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}

	switch config.ReadPreference {
	case "SlaveOnly":
		ops.SlaveOnly = true
	}

	return redis.NewFailoverClient(ops)
}

func NewFailoverCluster(config Config) *redis.ClusterClient {
	ops := &redis.FailoverOptions{
		MasterName:         config.MasterName,
		SentinelAddrs:      config.SentinelAddrs,
		SentinelPassword:   config.SentinelPassword,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.Db,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}

	switch config.ReadPreference {
	case "SlaveOnly":
		ops.SlaveOnly = true
	case "RouteByLatency":
		ops.RouteByLatency = true
	case "RouteRandomly":
		ops.RouteRandomly = true
	}

	return redis.NewFailoverClusterClient(ops)
}

func Finally() {
	for _, cli := range _clis {
		_ = cli.Close()
	}
	for _, clu := range _clus {
		_ = clu.Close()
	}
}

func GetCli(id string) *redis.Client {
	return _clis[id]
}

func GetCliDefault() *redis.Client {
	return GetCli(DefaultId)
}

func GetClu(id string) *redis.ClusterClient {
	return _clus[id]
}

func GetCluDefault() *redis.ClusterClient {
	return GetClu(DefaultId)
}
