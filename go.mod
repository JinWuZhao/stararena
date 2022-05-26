module github.com/jinwuzhao/stararena

go 1.18

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/JinWuZhao/bilidanmu v0.0.0-20220507054202-d5d89dcf88ea
	github.com/JinWuZhao/sc2client v0.0.0-20220512142018-b760770e746b
	github.com/prataprc/goparsec v0.0.0-20211219142520-daac0e635e7e
	go.uber.org/atomic v1.9.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)

replace github.com/JinWuZhao/sc2client => ../sc2client
