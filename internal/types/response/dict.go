package response

import (
	"LeotureWeb/internal/model"
	"time"

	"github.com/google/uuid"
)

// DictResp 字典响应
type DictResp struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Category  string       `json:"category"`
	Key       string       `json:"key"`
	Val       string       `json:"val"`
	Sort      int          `json:"sort"`
	IsDefault bool         `json:"isDefault"`
	Status    model.Status `json:"status"`
	Remark    string       `json:"remark"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

// ToDictResp 模型转响应
func ToDictResp(m *model.Dict) *DictResp {
	return &DictResp{
		ID:        m.ID,
		Name:      m.Name,
		Category:  m.Category,
		Key:       m.Key,
		Val:       m.Val,
		Sort:      m.Sort,
		IsDefault: m.IsDefault,
		Status:    m.Status,
		Remark:    m.Remark,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
