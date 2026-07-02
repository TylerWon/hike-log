package models

import (
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	HikeID  uint
	SrcUrl  string
	Caption *string
}
