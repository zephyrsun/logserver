package logger

type Logger interface {
	Write([]byte)
	Close()
}
