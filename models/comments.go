package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Name      string `gorm:"type:text" json:"name"`
	IsPrivate bool   `gorm:"default:false" json:"is_private"`
	Password  string `json:"password"`
	Comment   string `gorm:"type:text" json:"comment"`
}

func (comment Comment) TableName() string {
	return "Comment"
}
