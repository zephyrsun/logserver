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
	DumpError(err)

	defer conn.Close()
	Dump("Listening - %s", conn.LocalAddr())

	this.Read(conn)
}


func (this *LogServer) Read(conn net.PacketConn) {
	c := make(chan []byte, runtime.NumCPU())

	buf := make([]byte, 2048) //var buf [2048]byte

	go func() {
		for {
			n , _, err := conn.ReadFrom(buf)
			if err == nil {
				c <-buf[:n]
			}
		}
	}()

	for b := range c {
		go this.Parse(b)
	}

}

// &分隔
func (this *LogServer) Parse(b []byte) {

	//arr := strings.Split(log, "&")

	sep := "&"[0]

	start := 0
	for i := 0; i < len(b); i++ {
		if b[i] == sep {
			this.writer.Write(this.Format(b[start:i]))

			start = i+1
		}
	}

	//last one
	this.writer.Write(this.Format(b[start:]))
}

func (this *LogServer) Format(b []byte) (k string, buf []byte) {
	s := bytes.SplitN(b, []byte("="), 2)

	buf = append([]byte(this.timeNow.Format("2006-01-02 15:04:05")), "|"...)

	buf = append(buf, s[1]...)

	k = string(s[0])

	return
}

// e.g.
// 1=2014-07-10 13:23:46|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||
func (this *LogServer) Write(k string, b []byte) {
	//Dump("writing: %s", b)
	this.writer.Write(k, b)
}

/*
func prof(cpuprofile string, memprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		ErrorHandler(err)

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		ErrorHandler(err)

		defer pprof.WriteHeapProfile(f)
	}
}
*/

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
