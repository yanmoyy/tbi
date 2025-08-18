package database

import (
	"context"

	"github.com/yanmoyy/tbi/internal/models"
	"gorm.io/gorm"
)

func (c *Client) CreateBlockList(ctx context.Context, blocks []models.Block) error {
	return c.db.WithContext(ctx).Create(blocks).Error
}

func (c *Client) CreateBlock(ctx context.Context, block models.Block) error {
	return c.db.WithContext(ctx).Create(&block).Error
}

func (c *Client) GetLastBlockInfo(ctx context.Context) (height, totalTxs int, err error) {
	var lastBlock models.Block
	err = c.db.WithContext(ctx).
		Order("height DESC").First(&lastBlock).Error
	if err == gorm.ErrRecordNotFound {
		return -1, 0, nil
	}
	if err != nil {
		return -1, 0, err
	}
	return lastBlock.Height, lastBlock.TotalTxs, nil
}

func (c *Client) ClearBlocks() error {
	return c.db.Exec("TRUNCATE TABLE blocks CASCADE").Error
}

// Deprecated: ignore data corruption
//
// // missed: need to search missing blocks in DB
// func (c *Client) GetLastHeight() (int32, bool, error) {
// 	var totalBlocks int64
// 	err := c.db.Model(&models.Block{}).Count(&totalBlocks).Error
// 	if err != nil {
// 		return -1, false, fmt.Errorf("failed to count total blocks: %w", err)
// 	}
//
// 	if totalBlocks == 0 {
// 		return -1, false, nil
// 	}
//
// 	var lastBlock models.Block
// 	err = c.db.Order("height DESC").First(&lastBlock).Error
// 	if err != nil {
// 		return -1, false, fmt.Errorf("failed to get last block: %w", err)
// 	}
//
// 	missed := int64(lastBlock.Height) != totalBlocks-1
// 	return lastBlock.Height, missed, nil
// }
//
//
// // GetMissingHeights is triggered when there is missing block on database.
// // it will be rarely triggered.
// func (c *Client) GetMissingHeights() ([]int32, error) {
// 	var minMax struct {
// 		MinHeight int32
// 		MaxHeight int32
// 	}
// 	err := c.db.Model(&models.Block{}).
// 		Select("MIN(height) as min_height, MAX(height) as max_height").
// 		Scan(&minMax).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get min/max heights: %w", err)
// 	}
// 	if minMax.MaxHeight == 0 {
// 		return []int32{}, nil
// 	}
//
// 	var existing []int32
// 	err = c.db.Model(&models.Block{}).
// 		Order("height ASC").
// 		Pluck("height", &existing).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch block heights: %w", err)
// 	}
//
// 	var missing []int32
// 	need := minMax.MinHeight
// 	for _, height := range existing {
// 		for need < height {
// 			missing = append(missing, need)
// 			need++
// 		}
// 		need = height + 1
// 	}
// 	return missing, nil
// }
