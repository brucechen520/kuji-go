package prize

import (
	"github.com/brucechen520/kuji-go/internal/models"
)

func (s *service) Draw(userID int) (*models.Prize, error) {
	// 複雜的計算邏輯寫在這裡，這檔案就只會有這一個大 Function
	// 如果還有別的功能，就再開別的檔案實作，保證每個檔案不超過 100 行
	return s.repo.GetByID(1) // 這裡只是示範，實際邏輯會更複雜
}
