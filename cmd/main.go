package main

import (
	"fmt"
	"naka-disc/discord-bot-golang/internal/app/dao"
	"naka-disc/discord-bot-golang/internal/app/entity"
	"naka-disc/discord-bot-golang/internal/app/service"
	"naka-disc/discord-bot-golang/internal/app/util/dateutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 初期処理。
// エントリーポイントのmain関数の前に実行される。
func init() {
	fmt.Println("naka-disc/discord-bot-golang init")

	// TODO: DBのマイグレーション処理を暫定でここに挿入 別途コマンドとかにした方が良さげ
	db, err := gorm.Open(sqlite.Open("database/database.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(entity.NewVcAccessLog())
	db.AutoMigrate(entity.NewDiscordMember())
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
	dg.AddHandler(service.MessageDispatch) // メッセージ受信時
	dg.AddHandler(service.VoiceDispatch)   // ボイスチャンネルへの入退室時
	dg.AddHandler(guildMemberAdd)          // サーバーへの入室時

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

// Dicordセッションに追加するコールバック関数。
// guildMemberAdd: サーバーに入室した時にキックされる処理。
func guildMemberAdd(session *discordgo.Session, gma *discordgo.GuildMemberAdd) {
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
