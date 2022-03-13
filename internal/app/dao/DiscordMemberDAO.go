package dao

import (
	"naka-disc/discord-bot-golang/internal/app/entity"
)

// DiscordのメンバーIDでメンバー情報を取得。
func GetDiscordMemberByMemberId(memberId string) (*entity.DiscordMember, bool) {
	db := getDatabaseObject()

	// メンバー情報を取得
	// 取得できる場合は、eの中にデータがセットされる
	e := entity.NewDiscordMember()
	tx := db.Limit(1).Find(&e, "discord_member_id = ?", memberId)

	// RowsAffectedが1のときはデータあり、0の時はデータなしなので失敗で返す
	if tx.RowsAffected == 0 {
		return e, false
	}

	return e, true
}

// メンバー情報を登録。
func AddDiscordMember(e *entity.DiscordMember) {
	db := getDatabaseObject()
	db.Create(&e)
}

// メンバー情報を更新。
// 更新パラメータのMapを指定して、その部分だけ上書きをかけるイメージ。
func EditDiscordMember(e *entity.DiscordMember, om map[string]interface{}) {
	db := getDatabaseObject()
	db.Model(&e).Updates(om)
}
