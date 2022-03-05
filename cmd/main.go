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

	// Discordセッションの権限付与
	// IntentsGuildMessages はメッセージイベントに関するもの
	dg.Identify.Intents = discordgo.IntentsGuildMessages

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
