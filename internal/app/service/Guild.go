package service

import (
	"naka-disc/discord-bot-golang/internal/app/dao"
	"naka-disc/discord-bot-golang/internal/app/entity"
	"naka-disc/discord-bot-golang/internal/app/util/dateutil"

	"github.com/bwmarrin/discordgo"
)

// サーバー入室時のディスパッチ処理。
// 内容に応じて処理を分岐する。
func GuildAddDispatch(s *discordgo.Session, gma *discordgo.GuildMemberAdd) {
	// FIXME: 起動確認ができていないため、何らかの方法でテストを試したい

	// MEMO: メールアドレスも取れるは取れる(gma.User.Email)
	// ただし個人情報にあたるので、収集するのはリスキー いらないなら取らないのが無難

	// DiscordのメンバーIDから、メンバー情報取得
	e, ok := dao.GetDiscordMemberByMemberId(gma.User.ID)

	// あれば更新、なければ新規登録で分岐
	if ok {
		// データが既にある場合、一部情報を最新化して更新
		joinCount := e.JoinCount + 1
		om := map[string]interface{}{"join_count": joinCount, "is_stay": true, "last_join_datetime": dateutil.GetNowString()}
		dao.EditDiscordMember(e, om)

	} else {
		// 新規参加メンバーの場合、新規でデータを登録
		e := entity.NewDiscordMember()
		e.DiscordMemberId = gma.User.ID
		e.DiscordMemberName = gma.User.Username
		e.JoinCount = 1
		e.IsBot = gma.User.Bot
		e.IsStay = true
		e.FirstJoinDatetime = dateutil.GetNowString()
		e.LastJoinDatetime = dateutil.GetNowString()
		dao.AddDiscordMember(e)

	}
}
