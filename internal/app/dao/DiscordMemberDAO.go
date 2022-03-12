package dao

import (
	"naka-disc/discord-bot-golang/internal/app/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Discordのメンバー情報DAO。
type DiscordMemberDAO struct {
	DB *gorm.DB
}

// コンストラクタ用処理。
// この関数を実行してインスタンスを生成すること。
func NewDiscordMemberDAO() *DiscordMemberDAO {
	ret := new(DiscordMemberDAO)

	// TODO: BaseDAO作って、SQLiteのオープンはそっちに寄せたい
	// TODO: DB接続失敗は想定しづらいので、エラー時対処はスキップ
	db, _ := gorm.Open(sqlite.Open("database/database.sqlite"), &gorm.Config{})

	ret.DB = db
	return ret
}

// DiscordのメンバーIDでメンバー情報を取得。
func (instance *DiscordMemberDAO) GetDiscordMemberByMemberId(memberId string) (*entity.DiscordMember, bool) {
	getEntity := entity.NewDiscordMember()
	tx := instance.DB.Limit(1).Find(&getEntity, "discord_member_id = ?", memberId)

	// RowsAffectedが1のときはある、0の時はないので失敗で返す
	if tx.RowsAffected == 0 {
		return getEntity, false
	}

	return getEntity, true
}

// メンバー情報を新規登録。
func (instance *DiscordMemberDAO) AddDiscordMember(saveEntity *entity.DiscordMember) {
	instance.DB.Create(&saveEntity)
}

// メンバー情報を更新。
// 更新パラメータのMapを指定して、その部分だけ上書きをかけるイメージ。
func (instance *DiscordMemberDAO) EditDiscordMember(saveEntity *entity.DiscordMember, overrideMap map[string]interface{}) {
	instance.DB.Model(&saveEntity).Updates(overrideMap)
}
