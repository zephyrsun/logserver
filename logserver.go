package logserver

import (
	"net"
	"flag"
	"time"
)

type jsonConfigType map[string]string

var (
	sysLog *Logger
	config jsonConfigType
)

func Run() {

	sysLog = NewLogger(config["sys_log"])

	defer sysLog.Close()

	server := &LogServer{}

	server.Listen(config["address"])

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
	conn    net.PacketConn
	timeNow time.Time
	log *Logger
	logFlag string
	data    chan []byte
}

func (this *LogServer) Listen(addr string) {

	conn, err := net.ListenPacket("udp", addr)

	sysLog.Error(err)

	this.conn = conn

	defer conn.Close()

	this.initLogger()

	dump("Listening - %s", addr)

	this.read()
}

func (this *LogServer) initLogger() {

	var hourNow int

	f := func(now time.Time) bool {
		this.timeNow = now

		h := this.timeNow.Hour() //this.timeNow.Format("2006-01-02-03")

		if hourNow != h {

			hourNow = h

			if this.log != nil {
				this.log.Close()
			}

			file := config["save_path"] + this.timeNow.Format("2006-01-02-03")

			//dump("%s", file)

			this.log = NewLogger(file)
		}

		return true
	}

	f(time.Now())

	go ticker(1, f)
}

func (this *LogServer) read() {
	buf := make([]byte, 2048) //var buf [2048]byte

	for {
		n , _, err := this.conn.ReadFrom(buf)

		if err == nil {
			this.parse(buf[:n])
		}
	}
}

func (this *LogServer) parse(d []byte) {

	//arr := strings.Split(log, "&")

	sep := "&"[0]

	start := 0
	for i := 0; i < len(d); i++ {
		if d[i] == sep {
			this.write(d[start:i])

			start = i+1
		}
	}

	//last one
	this.write(d[start:])
}

func (this *LogServer) write(b []byte) {

	buf := append([]byte(this.timeNow.Format("2006-01-02 15:04:05")), "|"...)

	buf = append(buf, b...)

	sysLog.Log("writing: %s", buf)

	this.log.Write(buf)
}

func (this *LogServer) getDataLog() {

}

func init() {

	cfgFile := flag.String("c", "config.json", "set configuration file")

	flag.Parse()

	config = loadConfig(*cfgFile)
}
