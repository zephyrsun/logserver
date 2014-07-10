package logserver

import (
	"net"
	"flag"
	"time"
	"os"
	"runtime/pprof"
	"bytes"
	"runtime"
)

var logType = map[string]string {
	"1":"login",
	"2":"act",
	"3":"pay",
	"4":"item",
	"5":"error",
	"6":"funel",
	"7":"att",
}

var (
	sysLog *Logger
)

type LogServer struct{
	timeNow time.Time
	logger map[string]*Logger
}

func Listen() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	cfgFile := flag.String("c", "config.json", "Set configuration file")
	cpuprofile := flag.String("cpuprofile", "", "Write cpu profile to file")
	memprofile := flag.String("memprofile", "", "Write memory profile to file")

	flag.Parse()

	prof(*cpuprofile, *memprofile)

	loadConfig(*cfgFile)

	New().listen(config["address"])
}

func New() (server *LogServer) {
	sysLog = NewLogger(config["sys_log"])
	defer sysLog.Close()

	server = &LogServer{}

	server.logger = make(map[string]*Logger)

	server.initLogger()

	return
}

func (this *LogServer) listen(addr string) {
	conn, err := net.ListenPacket("udp", addr)
	PanicOnError(err)

	defer conn.Close()
	Dump("Listening - %s", conn.LocalAddr())

	c := make(chan []byte)

	this.Read(c, conn)
}

func (this *LogServer) initLogger() {

	var lastHour int = -1

	f := func(now time.Time) bool {
		this.timeNow = now

		h := this.timeNow.Hour() //this.timeNow.Format("2006-01-02-03")
		if lastHour != h {

			lastHour = h

			for k, v := range logType {
				file := config["save_dir"] + v + "_" + this.timeNow.Format("2006-01-02-15") + ".log"

				l, ok := this.logger[k]
				if ok {
					l.Close()
				}
				//Dump("Open file:%s", file)

				this.logger[k] = NewLogger(file)
			}
		}

		return true
	}

	f(time.Now())

	go ticker(1, f)
}

func (this *LogServer) Read(c chan []byte, conn net.PacketConn) {

	go func() {
		buf := make([]byte, 2048) //var buf [2048]byte
		for {
			n , _, err := conn.ReadFrom(buf)
			if err == nil {
				c <- buf[:n]
			}
		}
	}()

	for {
		this.Split(<-c)
	}
}

// &分隔
func (this *LogServer) Split(b []byte) {

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

// e.g.
// 1=2014-07-10 13:23:46|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||
func (this *LogServer) Write(b []byte) {

	s := bytes.SplitN(b, []byte("="), 2)

	buf := append([]byte(this.timeNow.Format("2006-01-02 15:04:05")), "|"...)

	buf = append(buf, s[1]...)

	//Dump("writing: %s", buf)

	this.logger[string(s[0])].Write(buf)
}

func prof(cpuprofile string, memprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		PanicOnError(err)

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		PanicOnError(err)

		defer pprof.WriteHeapProfile(f)
	}
}
