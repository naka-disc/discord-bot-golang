package dao

import "naka-disc/discord-bot-golang/internal/app/entity"

// 未退室のボイスチャンネル入退室ログを取得する。
// ユーザーID、ボイスチャンネルIDの一致と、退室日時が空という条件でデータを取れば、入室時のデータが取れる。
func GetVoiceChannelAccessLogForNotLeave(channelId string, memberId string) (*entity.VoiceChannelAccessLog, bool) {
	db := getDatabaseObject()

	// ボイスチャンネルへの入退室情報を取得
	// 取得できる場合は、eの中にデータがセットされる
	e := entity.NewVoiceChannelAccessLog()
	tx := db.Limit(1).Find(&e, "discord_member_id = ? AND voice_channel_id = ? AND leave_datetime = ''", memberId, channelId)

	// RowsAffectedが1のときはデータあり、0の時はデータなしなので失敗で返す
	if tx.RowsAffected == 0 {
		return e, false
	}

	return e, true
}

// ボイスチャンネル入退室ログを登録する。
// 登録時は退室日時と滞在時間が登録されない。
func AddVoiceChannelAccess(e *entity.VoiceChannelAccessLog) {
	db := getDatabaseObject()
	db.Create(&e)
}

// ボイスチャンネル入退室ログを更新する。
// 更新パラメータのMapを指定して、その部分だけ上書きをかけるイメージ。
func EditVoiceChannelAccess(e *entity.VoiceChannelAccessLog, mi map[string]interface{}) {
	db := getDatabaseObject()
	db.Model(&e).Updates(mi)
}
