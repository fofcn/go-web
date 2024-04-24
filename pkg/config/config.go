package config

type RedisStore struct {
	ClientName  string `env:"SCHEDULER_REDIS_CLIENTNAME"`
	ClusterMode string `env:"SCHEDULER_REDIS_CLUSTERMODE"`
}

type Config struct {
	// Http server config
	HttpServerConfig struct {
		Addr           string `env:"SERVER_ADDR, default=:8080"`
		GinMode        string `env:"SERVER_GINMODE"`
		ReadTimeout    int    `env:"SERVER_READTIMEOUT, default=10"`
		WriteTimeout   int    `env:"SERVER_WRITETIMEOUT, default=10"`
		MaxHeaderBytes int    `env:"SERVER_MAXHEADERBYTES, default=1048576000"`
	} `yaml:"server"`

	SchedulerConfig struct {
		LoadBalancer    string `env:"SCHEDULER_LOADBALANCER"`
		WorkerStoreType string `env:"SCHEDULER_WORKERSTORE"`
		PingInterval    int    `env:"SCHEDULER_PINGINTERVAL"`
		EvictThreshold  int    `env:"SCHEDULER_EVICT_THRESHOLD"`
		Store           RedisStore
	}
}

var ApplicationConfig Config

func GetHttpServerConfig() *struct {
	Addr           string `env:"SERVER_ADDR, default=:8080"`
	GinMode        string `env:"SERVER_GINMODE"`
	ReadTimeout    int    `env:"SERVER_READTIMEOUT, default=10"`
	WriteTimeout   int    `env:"SERVER_WRITETIMEOUT, default=10"`
	MaxHeaderBytes int    `env:"SERVER_MAXHEADERBYTES, default=1048576000"`
} {
	return &ApplicationConfig.HttpServerConfig
}
