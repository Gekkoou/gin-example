package global

import "gorm.io/plugin/soft_delete"

type MODEL struct {
	CreatedAt int                   `gorm:"column:created_at;type:int(11) unsigned;comment:创建时间;NOT NULL;autoCreateTime" json:"created_at"`
	UpdatedAt int                   `gorm:"column:updated_at;type:int(11) unsigned;comment:更新时间;NOT NULL;autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;type:int(11) unsigned;comment:删除时间;NOT NULL;default:0" json:"-"`
}
