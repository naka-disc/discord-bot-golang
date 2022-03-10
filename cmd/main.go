package main

import (
	"fmt"
	"naka-disc/discord-bot-golang/internal/app/entity"
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
	db.AutoMigrate(entity.NewVcAccessLogs())
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

// Dicordセッションに追加するコールバック関数。
// messageCreate: メッセージを受信した時にキックされる処理。
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

// Dicordセッションに追加するコールバック関数。
// voiceStateUpdate: VoiceChannelに入退室した時にキックされる処理。
func voiceStateUpdate(session *discordgo.Session, voiceState *discordgo.VoiceStateUpdate) {
	// 入退室のどちらかを取得
	// 退室時のステートはChannelIDがブランクになってるのでそれで判断
	isJoin := (voiceState.VoiceState.ChannelID != "")

	// 入室か退室かで処理を分岐
	if isJoin {
		// 入室時の処理
		fmt.Println("Join.")

		// 入室時は退室日時と滞在時間以外を登録
		// TODO: voiceStateにユーザー名とかがないので登録してない
		// どこかしらでユーザー情報を取得して、DiscordMemberTableとか作ってリレーション張った方がいいかも
		entity := entity.NewVcAccessLogs()
		entity.DiscordMemberId = voiceState.VoiceState.UserID
		entity.VoiceChannelId = voiceState.ChannelID
		entity.JoinDatetime = dateutil.GetNowString()

		db, err := gorm.Open(sqlite.Open("database/database.sqlite"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		db.Create(&entity)

	} else {
		// 退室時の処理
		fmt.Println("Leave.")

		// 退室時は、直近の入室のみデータを引っ張ってきて、それに退室情報を追加して更新をかける
		db, err := gorm.Open(sqlite.Open("database/database.sqlite"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// ユーザーID、VCIDの一致と、退室日時が空という条件でデータを取れば、入室時のデータが取れる
		// 入室がなくて退室にきた、というケースは想定しない というかあり得ない
		entity := entity.NewVcAccessLogs()
		userId := voiceState.VoiceState.UserID
		channelId := voiceState.BeforeUpdate.ChannelID
		d := db.Limit(1).Find(&entity, "discord_member_id = ? AND voice_channel_id = ? AND leave_datetime = ''", userId, channelId)

		// RowsAffectedが1のときはある、0の時はないので終了
		if d.RowsAffected == 0 {
			fmt.Println("Record Not Found")
			return
		}

		// 退室日時と、滞在時間を
		// TODO: 日付文字列をstringに変換 エラーは可能性としてあるが、一旦無視してる リファクタリング時に考える
		leavetime := dateutil.GetNowString()
		staySecond, ok := dateutil.DiffSecond(entity.JoinDatetime, leavetime)
		if !ok {
			// TODO: 失敗時どうするか
		}

		// 退室日時と滞在時間を入れて更新かける
		db.Model(&entity).Updates(
			map[string]interface{}{"leave_datetime": leavetime, "stay_second": staySecond})

	}
}
