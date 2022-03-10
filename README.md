# discord-bot-golang
- DiscordのBot。
- Go言語で作成。


## 環境変数設定
```sh
# Botのトークン
export DISCORD_BOT_TOKEN=
```

## ビルド&実行
```sh
go mod tidy
go build -o . cmd/*.go
./main
# go build -o . cmd/*.go && ./main
```

## 依存関係パッケージ情報
- [github.com/bwmarrin/discordgo](https://pkg.go.dev/github.com/bwmarrin/discordgo)
  - BSD-3-Clause
  - [github.com/gorilla/websocket](https://pkg.go.dev/github.com/gorilla/websocket)
    - BSD-2-Clause
  - [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)
    - BSD-3-Clause
  - [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys)
    - BSD-3-Clause
- [gorm.io/gorm](https://pkg.go.dev/gorm.io/gorm)
  - MIT
- [gorm.io/driver/sqlite](https://pkg.go.dev/gorm.io/driver/sqlite)
  - MIT



