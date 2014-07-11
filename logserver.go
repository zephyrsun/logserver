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

	this.Read(conn)
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

				//new first, copy after
				nf := NewLogger(file)

				of, ok := this.logger[k]
				if ok {
					of.Close()
				}
				//Dump("Open file:%s", file)

				this.logger[k] = nf
			}
		}

		return true
	}

	f(time.Now())

	go ticker(1, f)
}

func (this *LogServer) Read(conn net.PacketConn) {

	max := runtime.NumCPU()
	i := 0

	var a [64]chan []byte

	// maximum i is len(a)-1
	for ; i < max; i++ {
		a[i] = make(chan []byte)

		go func() {
			this.Parse(<-a[i], this.Write)
		}()
	}

	i = 0
	buf := make([]byte, 2048) //var buf [2048]byte
	for {
		n , _, err := conn.ReadFrom(buf)
		if err == nil {
			a[i] <-buf[:n]

			i++
			if i == max {
				i = 0
			}
		}
	}
}

// &分隔
func (this *LogServer) Parse(b []byte, callback func(string, []byte)) {

	//arr := strings.Split(log, "&")

	sep := "&"[0]

	start := 0
	for i := 0; i < len(b); i++ {
		if b[i] == sep {
			callback(this.Format(b[start:i]))

			start = i+1
		}
	}

	//last one
	callback(this.Format(b[start:]))
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
	//Dump("writing: %s", buf)
	this.logger[k].Write(b)
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

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
