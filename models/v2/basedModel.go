package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	BaseModel struct {
		CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
		CreatedBy string    `json:"createdBy" gorm:"column:created_by"`
		UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
		UpdatedBy string    `json:"updatedBy" gorm:"column:updated_by"`
	}

	BaseModelSoftDelete struct {
		BaseModel
		DeletedAt gorm.DeletedAt `json:"deleteAt" gorm:"column:deleted_at"`
		DeletedBy *string        `json:"deleteBy" gorm:"column:deleted_by"`
	}
)
