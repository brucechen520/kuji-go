package prize

import (
	"github.com/brucechen520/kuji-go/internal/models"
	"github.com/brucechen520/kuji-go/internal/pkg"
	"github.com/brucechen520/kuji-go/internal/repository"
)

type Service interface {
	Draw(userID int) (*models.Prize, error)
}
type service struct {
	repo repository.PrizeRepo
	tm   *pkg.TransactionManager
}

func New(repo repository.PrizeRepo, tm *pkg.TransactionManager) Service {
	return &service{repo: repo, tm: tm}
}
