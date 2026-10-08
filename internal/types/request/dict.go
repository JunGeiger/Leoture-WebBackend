package request

import (
	"LeotureWeb/internal/model"
)

// DictCreate 创建字典请求
type DictCreate struct {
	Name      string       `json:"name" binding:"required,max=50"`
	Category  string       `json:"category" binding:"required,max=50"`
	Key       string       `json:"key" binding:"required,max=100"`
	Val       string       `json:"val" binding:"required,max=200"`
	Sort      int          `json:"sort" binding:"-"`
	IsDefault bool         `json:"is_default" binding:"-"`
	Status    model.Status `json:"status" binding:"required,oneof=1 2"`
	Remark    string       `json:"remark" binding:"max=255"`
}

// ToModel 转换为数据库模型
func (r *DictCreate) ToModel() *model.Dict {
	return &model.Dict{
		Name:      r.Name,
		Category:  r.Category,
		Key:       r.Key,
		Val:       r.Val,
		Sort:      r.Sort,
		IsDefault: r.IsDefault,
		Status:    r.Status,
		Remark:    r.Remark,
	}
}

// DictUpdate 更新字典请求
type DictUpdate struct {
	ID        string       `json:"id" binding:"required,uuid"`
	Name      string       `json:"name" binding:"required,max=50"`
	Category  string       `json:"category" binding:"required,max=50"`
	Key       string       `json:"key" binding:"required,max=100"`
	Val       string       `json:"val" binding:"required,max=200"`
	Sort      int          `json:"sort" binding:"-"`
	IsDefault bool         `json:"is_default" binding:"-"`
	Status    model.Status `json:"status" binding:"required,oneof=1 2"`
	Remark    string       `json:"remark" binding:"max=255"`
}

// ToModel 转换为数据库模型
func (r *DictUpdate) ToModel() *model.Dict {
	return &model.Dict{
		Name:      r.Name,
		Category:  r.Category,
		Key:       r.Key,
		Val:       r.Val,
		Sort:      r.Sort,
		IsDefault: r.IsDefault,
		Status:    r.Status,
		Remark:    r.Remark,
	}
}

// DictQuery 字典查询请求
type DictQuery struct {
	Name     string       `form:"name" binding:"-"`
	Category string       `form:"category" binding:"-"`
	Key      string       `form:"key" binding:"-"`
	Status   model.Status `form:"status" binding:"-,oneof=1 2"`
	PageParams
	DateRange
}
