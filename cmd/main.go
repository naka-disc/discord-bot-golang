package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func init() {
	fmt.Println("naka-disc/discord-bot-golang init")
}

func main() {
	// 環境変数からToken取得
	token := os.Getenv("DISCORD_BOT_TOKEN")

	// Discordセッションを生成
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		// TODO: エラーがあれば中断
		fmt.Println("error creating Discord session,", err)
		return
	}

	// イベントのコールバックを登録する
	dg.AddHandler(messageCreate)
	dg.AddHandler(voiceStateUpdate)

	// Discordセッションの権限付与
	// TODO: 適切な権限付与の方がいいと思うが、面倒なので全部の権限を与えている
	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// WebSocketコネクションを確立し、双方向通信できるようにする
	err = dg.Open()
	if err != nil {
		// コネクションが成立できなかった場合
		fmt.Println("error opening connection,", err)
		return
	}
	// 終了時にコネクションはクローズ
	defer dg.Close()

	// キーボードからの割り込みシグナルで終了するための処理
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

/*
Dicordセッションに追加するコールバック関数。

messageCreate: メッセージを受信した時にキックされる処理。
*/
func messageCreate(session *discordgo.Session, msgCreate *discordgo.MessageCreate) {
	// DiscordセッションのユーザーID（=BotのID）と発言者のIDが一緒なら何もしない
	if msgCreate.Author.ID == session.State.User.ID {
		return
	}

	// メッセージ内容が ping だったら Pong! をメッセージとして送信
	if msgCreate.Content == "ping" {
		session.ChannelMessageSend(msgCreate.ChannelID, "Pong!")
	}
}

/*
Dicordセッションに追加するコールバック関数。

voiceStateUpdate: VoiceChannelに入退室した時にキックされる処理。
*/
func voiceStateUpdate(session *discordgo.Session, voiceState *discordgo.VoiceStateUpdate) {
	// TODO: イベントフックだけで、実処理が書けていない

	// 入退室のどちらかを取得
	// 退室時のステートはChannelIDがブランクになってるのでそれで判断
	channelId := voiceState.VoiceState.ChannelID
	isJoin := (channelId != "")

	// 入室か退室かで処理を分岐
	if isJoin {
		// 入室時の処理
		fmt.Println("Join")

	} else {
		// 退室時の処理
		fmt.Println("Leave.")

	}
}
