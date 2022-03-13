package dao

import (
	"naka-disc/discord-bot-golang/internal/app/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GORMのデータベースオブジェクトを取得。
// daoパッケージ内のみでの使用に留めるため、外には公開しない。
func getDatabaseObject() *gorm.DB {
	// TODO: 接続エラーは発生しない想定で組んであるため、エラーハンドリングしていない
	db, _ := gorm.Open(sqlite.Open("database/database.sqlite"), &gorm.Config{})
	return db
}

// マイグレーション。
// エンティティを作成したらここに追加することを忘れないように。
func Migration() {
	db := getDatabaseObject()
	db.AutoMigrate(entity.NewDiscordMember())
	db.AutoMigrate(entity.NewVoiceChannelAccessLog())
}
