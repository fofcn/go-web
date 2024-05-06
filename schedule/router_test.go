package schedule_test

import (
	"errors"
	"go-web/pkg/config"
	"go-web/schedule"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	loadConfig()
	r := gin.Default()
	public := r.Group("/")
	schedule.InitRouter(public)
	return r
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

func TestInitRouter(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/schedule/worker", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
