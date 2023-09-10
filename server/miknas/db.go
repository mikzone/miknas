package miknas

import "time"

type PrimaryUintModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PrimaryStringModel struct {
	Sid       string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
