package logserver

import (
	"net"
	"flag"
	"time"
	"bytes"
	"runtime"
)

type LogWriter interface {
	Write(string, []byte) error
}

type LogServer struct{
	timeNow time.Time
	writer  LogWriter
}

func Listen() {

	cfgFile := flag.String("c", "config.json", "Set configuration file")
	w := flag.String("w", "file", "Set writer")

	flag.Parse()

	loadConfig(*cfgFile)

	New(*w).listen(Config["address"])
}

func New(w string) (s *LogServer) {

	s = &LogServer{}

	if w == "file" {
		s.writer = NewFileWriter()
	}    else {
		s.writer = NewConsoleWriter()
	}

	return s
	//make(map[string]*Writer)
}

func (this *LogServer) listen(addr string) {
	conn, err := net.ListenPacket("udp", addr)
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

func (this *LogServer) Read(conn net.PacketConn) {
	c := make(chan []byte, runtime.NumCPU())

	buf := make([]byte, 2048) //var buf [2048]byte

	go func() {
		for {
			n , _, err := conn.ReadFrom(buf)
			if err == nil {
				c <-buf[:n]
			}else {
				DumpError(err, false)
			}
		}
	}()

	for b := range c {
		this.Parse(b)
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

	this.writer.Write(string(s[0]), buf)
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
