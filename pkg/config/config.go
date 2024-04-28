package config

type RedisStore struct {
	ClientName  string   `env:"SCHEDULER_REDIS_CLIENTNAME"`
	ClusterMode string   `env:"SCHEDULER_REDIS_CLUSTERMODE"`
	Addrs       []string `env:"SCHEDULER_REDIS_ADDRS"`
	Username    string
	Password    string
	DB          int
	MaxRetries  int
	Timeout     int
	Pool        struct {
		Size        int
		MaxIdle     int
		MaxActive   int
		IdleTimeout int
	}
}

type Jwt struct {
	Secret    string
	ExpiresIn int
	Issuer    string
}

type Cookie struct {
	Name     string
	MaxAge   int
	Path     string
	Domain   string
	Security bool
	HttpOnly bool
}

type Auth struct {
	Type   string
	Cookie Cookie
	Jwt    Jwt
}

type Scheduler struct {
	LoadBalancer   string `env:"SCHEDULER_LOADBALANCER"`
	WorkerStore    string `env:"SCHEDULER_WORKERSTORE"`
	PingInterval   int    `env:"SCHEDULER_PINGINTERVAL"`
	EvictThreshold int    `env:"SCHEDULER_EVICT_THRESHOLD"`
	Redis          RedisStore
}

type Server struct {
	Addr           string `env:"SERVER_ADDR, default=:8080"`
	GinMode        string `env:"SERVER_GINMODE"`
	ReadTimeout    int    `env:"SERVER_READTIMEOUT, default=10"`
	WriteTimeout   int    `env:"SERVER_WRITETIMEOUT, default=10"`
	MaxHeaderBytes int    `env:"SERVER_MAXHEADERBYTES, default=1048576000"`
}

type Aliyun struct {
	Oss Oss
}

type Oss struct {
	EndPoint  string
	AccessKey string
	SecretKey string
}

type Config struct {
	Server    Server
	Scheduler Scheduler
	Aliyun    Aliyun
	Auth      Auth
}

var ApplicationConfig *Config = &Config{}

func GetHttpServerConfig() *Server {
	return &ApplicationConfig.Server
}

func GetAliyunOss() *Oss {
	return &ApplicationConfig.Aliyun.Oss
}

func GetScheduler() *Scheduler {
	return &ApplicationConfig.Scheduler
}

func GetAuthConfig() *Auth {
	return &ApplicationConfig.Auth
}
