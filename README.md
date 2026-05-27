## pterodactyl-go

Discord Bot を 1 プロセスで動かしつつ、機能を以下の 2 つに分割しています。

- `pterodactyl-bot/`: Pterodactyl パネル操作（コマンド系）
- `linkfixer/`: URL 書き換え（LinkFixer）

どちらも `main.go` から同じ Discord セッションへ登録され、1つの Bot として動作します。

### memo

.env は VPS に保管

### Docker image (GHCR)

- `ghcr.io/oudentabetai/pterodactyl-go:latest`
- `ghcr.io/oudentabetai/pterodactyl-go:<release-tag>` (例: `v1.0.5`)
