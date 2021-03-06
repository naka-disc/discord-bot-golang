package entity

// ボイスチャンネルへの入退室用エンティティ。
type VoiceChannelAccessLog struct {
	Id                         uint   `gorm:"primarykey"`
	DiscordMemberId            string // Discordユーザーに与えられる一意のID
	DiscordMemberName          string // Discordユーザー名
	DiscordMemberDiscriminator string // Discordユーザーの#0000の番号 正直いらない
	VoiceChannelId             string // ボイスチャンネルID
	JoinDatetime               string // 入室日時
	LeaveDatetime              string // 退室日時
	StaySecond                 uint   // 滞在時間(秒)
}

// コンストラクタ用処理。
// この関数を実行してインスタンスを生成すること。
func NewVoiceChannelAccessLog() *VoiceChannelAccessLog {
	ret := new(VoiceChannelAccessLog)
	return ret
}
