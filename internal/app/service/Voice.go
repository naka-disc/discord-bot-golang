package service

import (
	"log"
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
		joinVoiceChannel(vs)

	} else {
		// 退室時の処理
		leaveVoiceChannel(vs)

	}
}

// ボイスチャンネルに入室した時の処理。
func joinVoiceChannel(vs *discordgo.VoiceStateUpdate) {
	log.Println("Join.")

	// 入室時は、エンティティに退室日時と滞在時間以外を設定し、登録に回す
	// TODO: voiceStateにユーザー名とかがないので登録してない
	// どこかしらでユーザー情報を取得して、DiscordMemberTableとか作ってリレーション張った方がいいかも
	e := entity.NewVoiceChannelAccessLog()
	e.DiscordMemberId = vs.VoiceState.UserID
	e.VoiceChannelId = vs.ChannelID
	e.JoinDatetime = dateutil.GetNowString()

	dao.AddVoiceChannelAccess(e)
}

// ボイスチャンネルから退室した時の処理。
func leaveVoiceChannel(vs *discordgo.VoiceStateUpdate) {
	log.Println("Leave.")

	// 入室時のデータを取得
	channelId := vs.BeforeUpdate.ChannelID
	memberId := vs.VoiceState.UserID
	e, ok := dao.GetVoiceChannelAccessLogForNotLeave(channelId, memberId)
	if !ok {
		// 入室時データなし ここに来ることはほぼない想定
		log.Printf("leaveVoiceChannel error. ChId: %s, MemId: %s", channelId, memberId)
		return
	}

	// 退室日時と滞在時間を更新
	// 滞在時間を計算できなかった場合も更新(退室日時がわかっていればあとから計算して後付けが可能なため)
	leavetime := dateutil.GetNowString()
	staySecond, _ := dateutil.DiffSecond(e.JoinDatetime, leavetime)
	om := map[string]interface{}{"leave_datetime": leavetime, "stay_second": staySecond}
	dao.EditVoiceChannelAccess(e, om)
}
