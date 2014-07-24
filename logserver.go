package logserver

import (
	"net"
	"flag"
	"time"
	"bytes"
	"runtime"
)

const (
	bufSize = 1024 * 1024
)

type LogWriter interface {
	Write(string, []byte)
}

type FileLogWriter interface {
	Rotate(time.Time)
	Flush()
}

type LogServer struct{
	timeNow time.Time
	wr      LogWriter
	buf     []byte
}

func Listen() {

	cfgFile := flag.String("c", "config.json", "Set configuration file")

	flag.Parse()

	loadConfig(*cfgFile)

	New().listen(Config["address"])
}

func New() (s *LogServer) {

	s = &LogServer{buf:make([]byte, bufSize)}

	switch Config["writer"] {
	case "console":
		s.wr = NewConsoleWriter()
	case "file":
		fallthrough
	default:
		s.wr = NewFileWriter()
	}

	return s
	//make(map[string]*Writer)
}

func (o *LogServer) listen(addr string) {

	la, err := net.ResolveUDPAddr("udp", addr)
	DumpError(err, true)

	conn, err := net.ListenUDP("udp", la)
	DumpError(err, true)

	defer conn.Close()

	Dump("Listening - %s", conn.LocalAddr())

	o.Tick()

	o.Read(conn)
}

func (o *LogServer) Tick() {

	t := func(now time.Time) {
		o.timeNow = now
	}

	go Ticker(1*time.Second, t)

}

func (o *LogServer) Read(conn *net.UDPConn) {
	writeBuf := make(chan []byte, 15*bufSize)//, runtime.NumCPU()

	readBuf := make([]byte, 2048) //var buf [2048]byte

	/*
		count := 0
		go Ticker(1*time.Second, func(time.Time) {
				print(count)
			})
	*/


	go func() {
		for {
			n , _, err := conn.ReadFromUDP(readBuf)
			if err == nil {
				writeBuf <-readBuf[:n]
			}else {
				DumpError(err, false)
			}
		}
	}()

	flushTimer := time.Tick(1 * time.Second)
	for {
		select {
		case b := <-writeBuf:
			o.Parse(b)

		case now := <-flushTimer:
			if m, ok := o.wr.(FileLogWriter); ok {
				m.Flush()
				m.Rotate(now)
			}
		}
	}
}

// &分隔
func (o *LogServer) Parse(b []byte) {

	//arr := strings.Split(log, "&")

	sep := "&"[0]

	start := 0
	for i := 0; i < len(b); i++ {
		if b[i] == sep {
			o.Write(b[start:i])

			start = i+1
		}
	}

	//last one
	o.Write(b[start:])
}

func (o *LogServer) Write(b []byte) {

	s := bytes.SplitN(b, []byte("="), 2)

	o.buf = append([]byte(o.timeNow.Format("2006-01-02 15:04:05")), "|"...)
	o.buf = append(o.buf, s[1]...)

	o.wr.Write(string(s[0]), o.buf)
}

/*
func prof(cpuprofile string, memprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		DumpError(err, true)

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		DumpError(err, true)

		defer pprof.WriteHeapProfile(f)
	}
}
*/

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
