package main

import (
	"context"
	"errors"
	"go-web/auth"
	"go-web/converter"
	"go-web/file"
	"go-web/index"
	"go-web/pdf"
	"go-web/pkg/config"
	"go-web/pkg/middleware"
	"go-web/pkg/router"
	"go-web/schedule"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/viper"
)

func main() {
	log.SetOutput(os.Stdout)

	cfg, err := loadConfig()
	if err != nil {
		log.Println("load config error:", err)
		return
	}

	config.ApplicationConfig = cfg

	ctx := context.Background()
	if err := envconfig.Process(ctx, config.ApplicationConfig); err != nil {
		log.Println("proceesing env config error:", err)
		return
	}

	log.Println(cfg.Server.Addr)

	server := prepareServer()

	go func() {
		startServer(server)
	}()

	shutodwn := make(chan os.Signal, 1)
	signal.Notify(shutodwn, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutodwn
	if sig != syscall.SIGKILL {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownServer(server, ctx)
	}

}

func loadConfig() (*config.Config, error) {
	viper.SetConfigName("go-web.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return nil, errors.New("go-web.yaml cannot be found")
		} else {
			return nil, errors.New("parse config file error")
		}
	}

	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, errors.New("parse config file error")
	}
	return &cfg, nil
}

func prepareServer() *http.Server {
	// to set gin Mode, either you can use env or code
	// - using env:    export GIN_MODE=release
	// - using code:    gin.SetMode(gin.ReleaseMode)
	// if envValue, isExisting := os.LookupEnv("GIN_MODE"); isExisting {
	//     gin.SetMode(envValue)
	// } else {
	//     gin.SetMode(gin.DebugMode)
	// }
	gin.SetMode(config.GetHttpServerConfig().GinMode)

	log.Println("listen on : ", config.GetHttpServerConfig().Addr)

	r := router.CreateEngine()
	server := &http.Server{
		Addr:           config.GetHttpServerConfig().Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(config.GetHttpServerConfig().ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.GetHttpServerConfig().WriteTimeout) * time.Second,
		MaxHeaderBytes: config.GetHttpServerConfig().MaxHeaderBytes,
	}

	corsCfg := cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Cookie"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	public := r.Group("/")
	public.Use(cors.New(corsCfg))
	public.Use(middleware.ConfigContext(), middleware.OptionalToken())
	index.InitRouter(public)
	schedule.InitRouter(public)
	auth.InitRouter(public)
	pdf.InitRouter(public)

	private := public.Group("/")
	private.Use(cors.New(corsCfg))
	private.Use(middleware.MustAuth())
	file.InitRouter(private)
	file.InitRouterFile(private)
	converter.InitRouter(private)

	return server
}

func startServer(srv *http.Server) {
	log.Println("start server")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("startup http server error", err)
	}
	log.Println("server started")
}

func shutdownServer(srv *http.Server, ctx context.Context) {
	log.Println("server shutdown start")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("shutdown http server error", err)
	}
	log.Println("server shutdown")
}
