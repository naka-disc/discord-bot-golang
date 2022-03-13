package main

import (
	"fmt"
	"naka-disc/discord-bot-golang/internal/app/dao"
	"naka-disc/discord-bot-golang/internal/app/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// 初期処理。
// エントリーポイントのmain関数の前に実行される。
func init() {
	fmt.Println("naka-disc/discord-bot-golang init")

	// TODO: DBのマイグレーション処理を暫定でここに挿入 別途コマンドとかにした方が良さげ
	dao.Migration()

}

// エントリーポイント。
// 実行したらここが実行される。
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
	dg.AddHandler(service.MessageDispatch)     // メッセージ受信時
	dg.AddHandler(service.VoiceDispatch)       // ボイスチャンネルへの入退室時
	dg.AddHandler(service.GuildAddDispatch)    // サーバーへの入室時
	dg.AddHandler(service.GuildRemoveDispatch) // サーバーへの退室時

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
