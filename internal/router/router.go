package router

import (
	"errors"
	"fmt"

	"github.com/brucechen520/kuji-go/internal/pkg"
	"github.com/brucechen520/kuji-go/internal/pkg/core"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type resource struct {
	mux                core.Mux
	logger             *zap.Logger
	db                 *gorm.DB
	rdb                *redis.Client
	transactionManager *pkg.TransactionManager
}

type Server struct {
	Mux                core.Mux
	Db                 pkg.DbRepo
	Rdb                pkg.RedisRepo
	TransactionManager *pkg.TransactionManager
}

func NewHTTPServer(logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	r := new(resource)
	r.logger = logger

	// 1. 初始化你的 DB 和 Redis (原本在 app.go 做的事)
	db, err := pkg.NewDB()
	if err != nil {
		return nil, fmt.Errorf("db init failed: %w", err)
	}

	// 2. 檢查 d 是不是 nil (防禦性程式碼)
	if db == nil {
		return nil, errors.New("db instance is nil")
	}
	transactionManager := pkg.NewTransactionManager(db.GetDbW())
	rdb, _ := pkg.NewRedis()

	// 3. 初始化他的核心引擎
	mux, err := core.New(logger, core.WithDisableSwagger(), core.WithEnableRate())
	if err != nil {
		panic(err)
	}

	r.mux = mux
	r.db = db.GetDbW()     // 這裡直接傳 GORM 的 DB 連線，讓 Repository 可以使用
	r.rdb = rdb.GetRedis() // 這裡直接傳 Redis 客戶端，讓 Repository 可以使用
	r.transactionManager = transactionManager

	// 设置 API 路由
	setApiRouter(r)

	s := new(Server)
	s.Mux = mux
	s.Db = db   // 這裡直接傳 repo 讓 main.go 可以關閉 DB 連線
	s.Rdb = rdb // 這裡直接傳 repo 讓 main.go 可以關閉 Redis 連線
	s.TransactionManager = r.transactionManager

	return s, nil
}
