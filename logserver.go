package logserver

import (
	"net"
	"flag"
	"time"
	"bytes"
	"runtime"
)

type LogWriter interface {
	Write(string, []byte) (int, error)
}

type LogServer struct{
	timeNow time.Time
	wr      LogWriter
}

func Listen() {

	cfgFile := flag.String("c", "config.json", "Set configuration file")

	flag.Parse()

	loadConfig(*cfgFile)

	New().listen(Config["address"])
}

func New() (s *LogServer) {

	s = &LogServer{}

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

func (this *LogServer) listen(addr string) {

	la, err := net.ResolveUDPAddr("udp", addr)
	DumpError(err, true)

	conn, err := net.ListenUDP("udp", la)
	DumpError(err, true)

	defer conn.Close()

	Dump("Listening - %s", conn.LocalAddr())

	this.Tick()

	this.Read(conn)
}

func (this *LogServer) Tick() {

	t := func(now time.Time) {
		this.timeNow = now
	}

	go Ticker(1*time.Second, t)

}

func (this *LogServer) Read(conn *net.UDPConn) {
	ch := make(chan []byte, 4096)//, runtime.NumCPU()

	buf := make([]byte, 2048) //var buf [2048]byte

	go func() {
		for {
			this.Parse(<-ch)
		}
	}()

	for {
		n , _, err := conn.ReadFromUDP(buf)
		if err == nil {
			ch <-buf[:n]
		}else {
			DumpError(err, false)
		}
	}
}

// &分隔
func (this *LogServer) Parse(b []byte) {

	//arr := strings.Split(log, "&")

	sep := "&"[0]

	start := 0
	for i := 0; i < len(b); i++ {
		if b[i] == sep {
			this.Write(b[start:i])

			start = i+1
		}
	}

	//last one
	this.Write(b[start:])
}

func (this *LogServer) Write(b []byte) {

	s := bytes.SplitN(b, []byte("="), 2)

	buf := append([]byte(this.timeNow.Format("2006-01-02 15:04:05")), "|"...)
	buf = append(buf, s[1]...)

	_, err := this.wr.Write(string(s[0]), buf)
	DumpError(err, false)
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
