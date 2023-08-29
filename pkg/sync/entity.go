package sync

// 서버로 전송
type RootDirectory struct {
	UUID string
	Path string // 로컬의 절대경로 (ex: /home/user/Quics)
	Date string
}
