package entity

// Discordのメンバー情報エンティティ。
type DiscordMember struct {
	Id                uint   `gorm:"primarykey"`
	DiscordMemberId   string // Discordユーザーに与えられる一意のID
	DiscordMemberName string // Discordユーザー名
	JoinCount         uint   // 入室回数 繰り返し入退室する場合を考えて
	IsBot             bool   // Botかどうかの判定 trueならBot
	IsStay            bool   // 滞在しているかどうかの判定 trueなら滞在中
	FirstJoinDatetime string // 最初入室日時
	LastJoinDatetime  string // 最終入室日時
}

// コンストラクタ用処理。
// この関数を実行してインスタンスを生成すること。
func NewDiscordMember() *DiscordMember {
	ret := new(DiscordMember)
	return ret
}
