package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string `env:"APP_ENV" envDefault:"development"`
	App    AppConfig
	Auth   AuthConfig
	DB     DBConfig
	Redis  RedisConfig
	Logger LoggerConfig
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

type AppConfig struct {
	Port              string `env:"APP_PORT" envDefault:"3000"`
	DisablePProf      bool   `env:"DISABLE_PPROF" envDefault:"false"`
	DisableSwagger    bool   `env:"DISABLE_SWAGGER" envDefault:"false"`
	DisablePrometheus bool   `env:"DISABLE_PROMETHEUS" envDefault:"false"`
	EnableCors        bool   `env:"ENABLE_CORS" envDefault:"true"`
	EnableRate        bool   `env:"ENABLE_RATE_LIMIT" envDefault:"false"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

type AuthConfig struct {
	TokenTTL time.Duration `env:"AUTH_TOKEN_TTL" envDefault:"24h"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD"`
}

func Load() (*Config, error) {
	// 1. 先嘗試從 .env 檔案載入到系統環境變數
	// 注意：路徑要對，如果是從根目錄執行，直接寫 ".env" 即可
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Cannot load .env file")
	}

	cfg := &Config{}

	// 2. 再將系統環境變數（包含剛剛載入的）填入 struct
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
