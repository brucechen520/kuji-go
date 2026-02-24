package service // 定義 service 套件，負責核心商業邏輯

import (
	"context"                 // 引入 context 套件，用於傳遞請求上下文
	"kuji-go/internal/models" // 引入 models
	"kuji-go/internal/pkg"

	// 引入 repository 層，用於存取資料庫
	"github.com/redis/go-redis/v9" // 引入 Redis 套件
)

// PrizeRepository 定義了 PrizeService 所需的資料庫操作
// 這樣 Service 就不需要依賴具體的 Repository 實作 (解耦)
type PrizeRepository interface {
	GetPrizesByBoxID(ctx context.Context, boxID string) ([]models.Prize, error)
}

// PrizeService 結構體定義了獎品相關的服務
type PrizeService struct {
	Repo               PrizeRepository         // 改為依賴介面
	RDB                *redis.Client           // 持有 Redis Client 的依賴，用於操作快取
	TransactionManager *pkg.TransactionManager // 持有 Transaction Manager 的依賴，用於處理事務
}

// NewPrizeService 是建構函式，用於初始化 PrizeService
// 接收 Repository 和 Redis 作為參數 (依賴注入)
func NewPrizeService(repo PrizeRepository, rdb *redis.Client, txManager *pkg.TransactionManager) *PrizeService {
	return &PrizeService{Repo: repo, RDB: rdb, TransactionManager: txManager} // 回傳初始化後的 Service 指標
}

// GetPrizes 實作獲取獎品列表的商業邏輯
func (s *PrizeService) GetPrizes(ctx context.Context, boxID string) ([]string, error) {
	// 這裡未來可以加入快取邏輯 (Cache-Aside Pattern)
	// 例如：先查 Redis，沒資料再查 DB (s.Repo.GetPrizesByBoxID)
	return []string{"A賞-火拳艾斯", "B賞-魯夫"}, nil // 目前先回傳範例資料
}

// Draw 實作抽獎的商業邏輯
func (s *PrizeService) Draw(ctx context.Context) (string, error) {
	// 這裡未來會實作：檢查庫存 -> 計算機率 -> 扣庫存 -> 寫入 Redis 等複雜邏輯
	return "恭喜抽中 A 賞！", nil // 回傳抽獎結果
}
