package logserver

import (
	"net"
	"flag"
	"time"
	"os"
	"runtime/pprof"
	"bytes"
)

type jsonConfigType map[string]string

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
	config jsonConfigType

)

func prof(cpuprofile string, memprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		panicOnError(err)

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		panicOnError(err)

		defer pprof.WriteHeapProfile(f)
	}
}

func Listen() {

	cfgFile := flag.String("c", "config.json", "set configuration file")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flag.String("memprofile", "", "write memory profile to this file")

	flag.Parse()

	prof(*cpuprofile, *memprofile)

	config = loadConfig(*cfgFile)

	sysLog = NewLogger(config["sys_log"])
	defer sysLog.Close()

	server := &LogServer{}

	server.logger = make(map[string]*Logger)

	server.initLogger()

	server.listen(config["address"])

	/*
		var buf [2048]byte
		for {
			n , _, err := conn.ReadFrom(buf[0:])
			if err != nil {
				return
			}

			str := string(buf[:n])

			dump("recv: %s", str)
		}
		*/
}

type LogServer struct{
	timeNow time.Time
	logger map[string]*Logger
}

func (this *LogServer) listen(addr string) {
	conn, err := net.ListenPacket("udp", addr)
	panicOnError(err)

	defer conn.Close()
	dump("Listening - %s", conn.LocalAddr())

	c := make(chan []byte)

	this.read(c, conn)
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

				if this.logger[k] != nil {
					this.logger[k].Close()
				}

				this.logger[k] = NewLogger(file)
			}

			//sysLog.Log("Open file:%s", file)
		}

		return true
	}

	go ticker(1, f)
}

func (this *LogServer) read(c chan []byte, conn net.PacketConn) {
	buf := make([]byte, 2048) //var buf [2048]byte

	go func() {
		for {
			n , _, err := conn.ReadFrom(buf)
			if err == nil {
				c <- buf[:n]
			}
		}
	}()

	for {
		this.split(<-c)
	}
}

// &分隔
func (this *LogServer) split(b []byte) {

	//arr := strings.Split(log, "&")

	sep := "&"[0]

	start := 0
	for i := 0; i < len(b); i++ {
		if b[i] == sep {
			this.write(b[start:i])

			start = i+1
		}
	}

	//last one
	this.write(b[start:])
}

func (this *LogServer) write(b []byte) {

	s := bytes.SplitN(b, []byte("="), 2)

	buf := append([]byte(this.timeNow.Format("2006-01-02 15:04:05")), "|"...)

	buf = append(buf, s[1]...)

	//dump("writing: %s", buf)

	this.logger[string(s[0])].Write(buf)
}
