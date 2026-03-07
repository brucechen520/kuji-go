package client

import (
	"context"
	"errors"
	"testing"

	"github.com/brucechen520/kuji-go/internal/config"
	"github.com/brucechen520/kuji-go/internal/model"
	repomock "github.com/brucechen520/kuji-go/internal/repository/postgre/client/mock"
	redismock "github.com/brucechen520/kuji-go/internal/repository/redis/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetSeriesById_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repomock.NewMockSeriesRepository(ctrl)
	mockRedis := redismock.NewMockKujiStore(ctrl)

	service := NewSeriesService(mockRepo, mockRedis, &config.AuthConfig{})

	ctx := context.Background()
	seriesID := uint(1)
	boxID := uint(10)

	// Mock data
	mockSeries := &model.Series{
		Name:  "Test Series",
		Price: 100,
		Boxes: []model.Box{
			{
				SeriesID:     seriesID,
				LocationName: "Store A",
				Prizes: []model.Prize{
					{Level: "A", Name: "Prize A"},
					{Level: "B", Name: "Prize B"},
				},
			},
		},
	}
	mockSeries.ID = seriesID
	mockSeries.Boxes[0].ID = boxID
	mockSeries.Boxes[0].Prizes[0].ID = 1001
	mockSeries.Boxes[0].Prizes[1].ID = 1002

	mockInv := map[string]int{
		"1001": 5,
		"1002": 10,
	}

	// Expectations
	// 1. GetSeriesMeta from Redis hits successfully
	mockRedis.EXPECT().GetSeriesMeta(ctx, seriesID).Return(mockSeries, nil)

	// 2. GetBoxInventory from Redis hits successfully
	mockRedis.EXPECT().GetBoxInventory(ctx, boxID).Return(mockInv, nil)

	// Execute
	dto, err := service.GetSeriesById(ctx, seriesID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, dto)
	assert.Equal(t, seriesID, dto.ID)
	assert.Equal(t, "Test Series", dto.Name)

	require.Len(t, dto.Boxes, 1)
	assert.Equal(t, boxID, dto.Boxes[0].ID)

	require.Len(t, dto.Boxes[0].Prizes, 2)
	assert.Equal(t, 5, dto.Boxes[0].Prizes[0].RemainingQuantity)
	assert.Equal(t, 10, dto.Boxes[0].Prizes[1].RemainingQuantity)
}

func TestGetSeriesById_CacheMiss_DBHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repomock.NewMockSeriesRepository(ctrl)
	mockRedis := redismock.NewMockKujiStore(ctrl)

	service := NewSeriesService(mockRepo, mockRedis, &config.AuthConfig{})

	ctx := context.Background()
	seriesID := uint(2)
	boxID := uint(20)

	// Mock data
	mockSeries := &model.Series{
		Name:  "Test Series 2",
		Price: 200,
		Boxes: []model.Box{
			{
				SeriesID:     seriesID,
				LocationName: "Store B",
				Prizes: []model.Prize{
					{Level: "A", Name: "Prize A"},
				},
			},
		},
	}
	mockSeries.ID = seriesID
	mockSeries.Boxes[0].ID = boxID
	mockSeries.Boxes[0].Prizes[0].ID = 2001

	dbPrizes := []model.Prize{
		{RemainingQuantity: 7},
	}
	dbPrizes[0].ID = 2001

	// Expectations
	// 1. GetSeriesMeta from Redis misses
	mockRedis.EXPECT().GetSeriesMeta(ctx, seriesID).Return(nil, errors.New("redis nil"))

	// 2. Fallback to DB
	mockRepo.EXPECT().GetSeriesById(ctx, seriesID).Return(mockSeries, nil)

	// 3. Write back to Redis
	mockRedis.EXPECT().SetSeriesMeta(ctx, seriesID, mockSeries).Return(nil)

	// 4. GetBoxInventory from Redis misses (returns nil without error per kuji_store.go)
	mockRedis.EXPECT().GetBoxInventory(ctx, boxID).Return(nil, nil)

	// 5. Fallback to DB for inventory
	mockRepo.EXPECT().GetBoxInventoryById(ctx, boxID).Return(dbPrizes, nil)

	// 6. Write inventory back to Redis
	mockRedis.EXPECT().SetBoxInventory(ctx, boxID, map[string]int{"2001": 7}).Return(nil)

	// Execute
	dto, err := service.GetSeriesById(ctx, seriesID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, dto)
	assert.Equal(t, seriesID, dto.ID)
	assert.Equal(t, "Test Series 2", dto.Name)

	require.Len(t, dto.Boxes, 1)
	require.Len(t, dto.Boxes[0].Prizes, 1)
	assert.Equal(t, 7, dto.Boxes[0].Prizes[0].RemainingQuantity)
}

func TestGetSeriesById_DBNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repomock.NewMockSeriesRepository(ctrl)
	mockRedis := redismock.NewMockKujiStore(ctrl)

	service := NewSeriesService(mockRepo, mockRedis, &config.AuthConfig{})

	ctx := context.Background()
	seriesID := uint(3)

	// Expectations
	// Redis misses
	mockRedis.EXPECT().GetSeriesMeta(ctx, seriesID).Return(nil, errors.New("redis nil"))

	// DB completely misses (returns nil, nil)
	mockRepo.EXPECT().GetSeriesById(ctx, seriesID).Return(nil, nil)

	// Execute
	dto, err := service.GetSeriesById(ctx, seriesID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, "series 3 not found", err.Error())
	assert.Nil(t, dto)
}

func TestGetSeriesById_InventoryDBNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repomock.NewMockSeriesRepository(ctrl)
	mockRedis := redismock.NewMockKujiStore(ctrl)

	service := NewSeriesService(mockRepo, mockRedis, &config.AuthConfig{})

	ctx := context.Background()
	seriesID := uint(4)
	boxID := uint(40)

	// Mock data
	mockSeries := &model.Series{
		Name:  "Test Series 4",
		Boxes: []model.Box{
			{
				SeriesID:     seriesID,
				LocationName: "Store C",
				Prizes: []model.Prize{
					{Level: "A", Name: "Prize A"},
				},
			},
		},
	}
	mockSeries.ID = seriesID
	mockSeries.Boxes[0].ID = boxID
	mockSeries.Boxes[0].Prizes[0].ID = 4001

	// Expectations
	// 1. GetSeriesMeta from Redis hits successfully
	mockRedis.EXPECT().GetSeriesMeta(ctx, seriesID).Return(mockSeries, nil)

	// 2. GetBoxInventory misses
	mockRedis.EXPECT().GetBoxInventory(ctx, boxID).Return(nil, nil)

	// 3. Fallback to DB for inventory but returns empty
	mockRepo.EXPECT().GetBoxInventoryById(ctx, boxID).Return([]model.Prize{}, nil)

	// Execute
	dto, err := service.GetSeriesById(ctx, seriesID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, "Prizes 4 not found", err.Error())
	assert.Nil(t, dto)
}
