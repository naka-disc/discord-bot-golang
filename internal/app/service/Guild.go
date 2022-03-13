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
	memDAO := dao.NewDiscordMemberDAO()
	getEntity, ok := memDAO.GetDiscordMemberByMemberId(gma.User.ID)

	// あれば更新、なければ新規登録で分岐
	if ok {
		// データが既にある場合
		// 一部情報を最新化して更新
		joinCount := getEntity.JoinCount + 1
		overrideMap := map[string]interface{}{"join_count": joinCount, "is_stay": true, "last_join_datetime": dateutil.GetNowString()}
		memDAO.EditDiscordMember(getEntity, overrideMap)

	} else {
		// 新規参加メンバーの場合
		saveEntity := entity.NewDiscordMember()
		saveEntity.DiscordMemberId = gma.User.ID
		saveEntity.DiscordMemberName = gma.User.Username
		saveEntity.DiscordMemberDiscriminator = gma.User.Discriminator
		saveEntity.JoinCount = 1
		saveEntity.IsBot = gma.User.Bot
		saveEntity.IsStay = true
		saveEntity.FirstJoinDatetime = dateutil.GetNowString()
		saveEntity.LastJoinDatetime = dateutil.GetNowString()
		memDAO.AddDiscordMember(saveEntity)

	}
}
