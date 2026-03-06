package mapper

import (
	"github.com/brucechen520/kuji-go/internal/dto"
	"github.com/brucechen520/kuji-go/internal/model"
)

// 轉換單個 Prize
func MapPrizeToDetailDTO(p *model.Prize) dto.PrizeDetailDTO {
	return dto.PrizeDetailDTO{
		ID:                p.ID,
		Name:              p.Name,
		Level:             p.Level,
		RemainingQuantity: p.RemainingQuantity, // 初始數，之後由 Service 補上動態值
	}
}

// 轉換單個 Box
func MapBoxToDetailDTO(b *model.Box) dto.BoxDetailDTO {
	dtoBox := dto.BoxDetailDTO{
		ID:           b.ID,
		LocationName: b.LocationName,
		Prizes:       make([]dto.PrizeDetailDTO, len(b.Prizes)),
	}

	for i, p := range b.Prizes {
		dtoBox.Prizes[i] = MapPrizeToDetailDTO(&p)
	}
	return dtoBox
}

func MapSeriesToDetailDTO(s *model.Series) *dto.SeriesDetailDTO {
	dtoSeries := &dto.SeriesDetailDTO{
		ID:    s.ID,
		Name:  s.Name,
		Price: s.Price,
		Boxes: make([]dto.BoxDetailDTO, len(s.Boxes)),
	}

	for i, b := range s.Boxes {
		// 直接呼叫剛剛寫好的 Box Mapper
		dtoSeries.Boxes[i] = MapBoxToDetailDTO(&b)
	}
	return dtoSeries
}
