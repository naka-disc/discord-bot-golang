package service

import "github.com/bwmarrin/discordgo"

func MessageDispatch(s *discordgo.Session, mc *discordgo.MessageCreate) {
	// DiscordセッションのユーザーID（=BotのID）と発言者のIDが一緒なら何もしない
	if mc.Author.ID == s.State.User.ID {
		return
	}

	// メッセージ内容に応じたディスパッチ
	switch mc.Content {
	case "ping":
		s.ChannelMessageSend(mc.ChannelID, "Pong!")

	}
}
