/**
* @author: gongquanlin
* @since: 2024/3/23
* @desc:
 */

package domain

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var Db *gorm.DB

func InitDB() {
	// 连接 SQLite 数据库
	var err error
	Db, err = gorm.Open(sqlite.Open("webhook.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移模式，如果没有TokenScript表则创建
	err = Db.AutoMigrate(&TokenScript{}, &User{})
	if err != nil {
		panic("failed to migrate model")
	}
}

// TokenScript 结构体用于映射数据库中的表结构
type TokenScript struct {
	ID     uint   `gorm:"primaryKey" json:"id,omitempty"`
	UserId int    `gorm:"not null" json:"user_id,omitempty"`
	Token  string `gorm:"unique" json:"token,omitempty"`
	Script string `json:"script,omitempty"`
}

type User struct {
	ID           int       `gorm:"primaryKey"`
	Account      string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	LastLoginAt  time.Time `gorm:"autoUpdateTime"`
}
