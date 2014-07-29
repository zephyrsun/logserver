package logserver

import (
	"net"
	"flag"
	"time"
	"runtime"
	"os/signal"
	"os"
	"syscall"
	"runtime/pprof"
)

type LogWriter interface {
	Write(string, []byte)
}

type BufLogWriter interface {
	Rotate(time.Time)
	Flush()
}

type LogServer struct{
	timeString string
	wr         LogWriter
	buf        []byte
}

func Listen() *LogServer {

	c := flag.String("c", "config.json", "Set configuration file")

	profile := flag.Bool("profile", false, "file for memory profile")

	flag.Parse()

	if *profile {
		Dump("start cpuprofile")
		cpu, err := os.Create("cpu.prof")
		DumpError(err, true)
		pprof.StartCPUProfile(cpu)
		defer pprof.StopCPUProfile()

		Dump("start memprofile")
		mem, err := os.Create("mem.prof")
		DumpError(err, true)
		defer pprof.WriteHeapProfile(mem)
	}

	s := &LogServer{buf:make([]byte, 512)}
	switch Config["writer"] {
	case "console":
		s.wr = NewConsoleWriter()
	case "file":
		s.wr = NewFileWriter()
	case "buffer":
		fallthrough
	default:
		s.wr = NewBufferWriter()
	}

	loadConfig(*c)

	s.listen(Config["address"])

	return s

}

func (o *LogServer) listen(addr string) {

	o.Tick()

	la, err := net.ResolveUDPAddr("udp", addr)
	DumpError(err, true)

	conn, err := net.ListenUDP("udp", la)
	DumpError(err, true)
	//defer conn.Close()

	Dump("Listening - %s", conn.LocalAddr())

	o.Read(conn)
}

func (o *LogServer) Tick() {

	t := func(now time.Time) {
		o.timeString = now.Format("2006-01-02 15:04:05")
	}

	t(time.Now())

	go Ticker(1*time.Second, t)

}

func (o *LogServer) Read(conn *net.UDPConn) {
	writeBuf := make(chan []byte, bufSize)//, runtime.NumCPU()

	readBuf := make([]byte, 2048) //var buf [2048]byte

	go func() {
		for {
			n, err := conn.Read(readBuf)
			//n , _, err := conn.ReadFromUDP(readBuf)
			if err == nil {
				writeBuf <-readBuf[:n]
			}else {
				DumpError(err, false)
			}
		}
	}()

	go func() {
		bufTimer := time.Tick(1 * time.Second)
		for {
			select {
			case now := <-bufTimer:
				if m, ok := o.wr.(BufLogWriter); ok {
					m.Rotate(now)
					m.Flush()
				}

			case b := <-writeBuf:
				o.Parse(b)
			}
		}
	}()

	// listen exit
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	Dump("received signal:%v", <-exitSignal)

	if m, ok := o.wr.(BufLogWriter); ok {
		Dump("Flushing data...")
		m.Flush()
	}
}

// &分隔
func (o *LogServer) Parse(b []byte) {
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
	sep := "="[0]

	i := 0
	for ; i < len(b); i++ {
		if b[i] == sep {
			break
		}
	}

	d := o.buf[:0]

	d = append(d, []byte(o.timeString)...)
	d = append(d, "|"...)
	d = append(d, b[i+1:]...)
	d = append(d, eol...)
	//d := append([]byte(o.timeString), "|"...)
	//d = append(d, b[i+1:]...)

	o.wr.Write(string(b[:i]), d)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
