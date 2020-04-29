package logger

type Logger interface {
	Init()
	Write([]byte)
}
