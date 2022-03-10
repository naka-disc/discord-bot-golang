package dao

import (
	"naka-disc/discord-bot-golang/internal/app/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ボイスチャンネルへの入退室記録DAO。
type VcAccessLogDAO struct {
	DB *gorm.DB
}

// コンストラクタ用処理。
// この関数を実行してインスタンスを生成すること。
func NewVcAccessLogDAO() *VcAccessLogDAO {
	ret := new(VcAccessLogDAO)

	// TODO: BaseDAO作って、SQLiteのオープンはそっちに寄せたい
	// TODO: DB接続失敗は想定しづらいので、エラー時対処はスキップ
	db, _ := gorm.Open(sqlite.Open("database/database.sqlite"), &gorm.Config{})

	ret.DB = db
	return ret
}

// 未退室の入退室ログを取得する。
// ユーザーID、VCIDの一致と、退室日時が空という条件でデータを取れば、入室時のデータが取れる。
func (instance *VcAccessLogDAO) GetVcAccessLogForNotLeave(channelId string, memberId string) (*entity.VcAccessLog, bool) {
	getEntity := entity.NewVcAccessLog()
	tx := instance.DB.Limit(1).Find(&getEntity, "discord_member_id = ? AND voice_channel_id = ? AND leave_datetime = ''", memberId, channelId)

	// RowsAffectedが1のときはある、0の時はないので失敗で返す
	if tx.RowsAffected == 0 {
		return getEntity, false
	}

	return getEntity, true
}

// ボイスチャンネルへの入退室記録を新規登録する。
func (instance *VcAccessLogDAO) AddVcAccessLog(saveEntity *entity.VcAccessLog) {
	instance.DB.Create(&saveEntity)
}

// ボイスチャンネルへの入退室記録を更新する。
// 更新パラメータのMapを指定して、その部分だけ上書きをかけるイメージ。
func (instance *VcAccessLogDAO) EditVcAccessLog(saveEntity *entity.VcAccessLog, overrideMap map[string]interface{}) {
	instance.DB.Model(&saveEntity).Updates(overrideMap)
}
