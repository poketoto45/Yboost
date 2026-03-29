package main

import "gorm.io/gorm"

type gameTime struct {
	game1 string `gorm:"primaryKey"`
	game2 string `gorm:"primaryKey"`
	game3 string `gorm:"primaryKey"`
}
