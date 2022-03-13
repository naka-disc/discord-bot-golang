package service

import (
	"fmt"
	"naka-disc/discord-bot-golang/internal/app/dao"
	"naka-disc/discord-bot-golang/internal/app/entity"
	"naka-disc/discord-bot-golang/internal/app/util/dateutil"

	"github.com/bwmarrin/discordgo"
)

// ボイスチャンネル入退室時のディスパッチ処理。
// 内容に応じて処理を分岐する。
func VoiceDispatch(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	// 入退室のどちらかを取得
	// 退室時のステートはChannelIDがブランクになってるのでそれで判断
	isJoin := (vs.VoiceState.ChannelID != "")

	// 入室か退室かで処理を分岐
	if isJoin {
		// 入室時の処理
		fmt.Println("Join.")

		// 入室時は退室日時と滞在時間以外を登録
		// TODO: voiceStateにユーザー名とかがないので登録してない
		// どこかしらでユーザー情報を取得して、DiscordMemberTableとか作ってリレーション張った方がいいかも
		saveEntity := entity.NewVcAccessLog()
		saveEntity.DiscordMemberId = vs.VoiceState.UserID
		saveEntity.VoiceChannelId = vs.ChannelID
		saveEntity.JoinDatetime = dateutil.GetNowString()

		valDAO := dao.NewVcAccessLogDAO()
		valDAO.AddVcAccessLog(saveEntity)

	} else {
		// 退室時の処理
		fmt.Println("Leave.")

		// 退室時は、直近の入室のみデータを引っ張ってきて、それに退室情報を追加して更新をかける
		memberId := vs.VoiceState.UserID
		channelId := vs.BeforeUpdate.ChannelID
		valDAO := dao.NewVcAccessLogDAO()
		getEntity, ok := valDAO.GetVcAccessLogForNotLeave(channelId, memberId)
		if !ok {
			// データ取得できなかったので終了
			return
		}

		// 退室日時と滞在時間を更新
		// 滞在時間を計算できなかった場合も更新 退室日時がわかっていればあとから計算して後付けが可能なため
		leavetime := dateutil.GetNowString()
		staySecond, _ := dateutil.DiffSecond(getEntity.JoinDatetime, leavetime)
		overrideMap := map[string]interface{}{"leave_datetime": leavetime, "stay_second": staySecond}
		valDAO.EditVcAccessLog(getEntity, overrideMap)
	}
}
