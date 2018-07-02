package eslog

import (
	"net"
	"log"
	"fmt"
	"sync"
	"io"
	"github.com/sambios/goapl/eslog/telnet"
	"strings"
)

type TelnetCmdFunc func(args ...interface{})
type TelnetCmd struct {
	usage string
	cmdFunc TelnetCmdFunc
	m interface{}
}

// This log writer sends output to a file
type TelnetLogWriter struct {
	chanRecord chan *LogRecord
	listenConn net.Listener
	remoteConn net.Conn
	localPort int16
	isStopRun bool
	format string
	wg sync.WaitGroup
	myCmds map[string]*TelnetCmd
}


// Constructor
func NewTelnetLogWriter(port int16) *TelnetLogWriter {
	c := &TelnetLogWriter{
		format: "%T %D|%C|%L|(%S) %M",
		chanRecord:make(chan *LogRecord),
		localPort:port,
		myCmds:make(map[string]*TelnetCmd),
	}

	c.wg.Add(2)

	var err error
	c.listenConn, err = net.Listen("tcp", ":3000")
	if nil != err {
		log.Println(err)
		return nil
	}

	c.RegCommand("help", c, dbgHelp, "Print all commands")

	go c.routineAyncWrite()
	go c.acceptRoutine()

	return c
}

func (this *TelnetLogWriter)RegCommand(name string, mp interface{}, handler TelnetCmdFunc, usage string) {
	this.myCmds[name] = &TelnetCmd{usage:usage, cmdFunc:handler, m:mp}
}


func (this *TelnetLogWriter)routineLogCmd() {

	conn := this.remoteConn
	defer conn.Close()
	defer log.Printf("Connection from %s closed", conn.RemoteAddr())

	// Create telnet ReadWriter with no options.
	tn := telnet.NewReadWriter(conn)

	// Welcome banner.
	tn.Write([]byte("********************************\r\n"))
	tn.Write([]byte("   Welcome to debug console!\r\n"))
	tn.Write([]byte("********************************\r\n"))

	// Process input until connection is closed.
	buf := make([]byte, 1024)
	for {
		tn.Write([]byte("#"))
		n, err := tn.Read(buf)
		if err == io.EOF {
			return
		}

		//log.Printf("Read '%s' {% [1]x} n=%d", buf[:n], n)

		cmdline:= string(buf[:n])
		cmdline = strings.TrimRight(cmdline, "\r\n")


		args := strings.Split(cmdline, " ")
		if len(args) == 1 {
			args = strings.Split(cmdline, ",")
		}

		cmd := args[0]

		if myCmd, ok := this.myCmds[cmd]; ok {
			myCmd.cmdFunc(myCmd.m, args)
		}

		if strings.Compare(cmd, "bye") == 0 {
			fmt.Println("bye....")
             break
		}
	}

}

func (this *TelnetLogWriter)routineAyncWrite() {

	defer func () {
		this.wg.Done()
		fmt.Println("routineAyncWrite exit!")
	}()

	for rec := range this.chanRecord {
		// Check module status
		if this.isStopRun {
			break
		}

		var txt string
		if rec.logType == 1 {
			txt = fmt.Sprintf("%s", rec.message)
		}else{
			txt = FormatLogRecord(this.format, rec)
		}

		if nil == this.remoteConn {
			fmt.Printf(txt)
		}else{
			this.remoteConn.Write([]byte(txt))
		}
	}
}

func (this *TelnetLogWriter)acceptRoutine(){

	defer func () {
		this.wg.Done()
		fmt.Println("acceptRoutine exit!")
	}()


	for {
		conn, err := this.listenConn.Accept()
		if err != nil {
			log.Println("accept routine quit.error:", err)
			return
		}

		// Check module is still running
		if this.isStopRun {
			break
		}


		if nil != this.remoteConn {
			// only support one connection
			this.remoteConn.Close()
		}

		this.remoteConn = conn
		go this.routineLogCmd()
	}
}


//
// Implements LogWriter interface
//

func (this *TelnetLogWriter) Name() string {
	return "TelnetLogWriter"
}

// This will be called to log a LogRecord message.
func (this *TelnetLogWriter) LogWrite(rec *LogRecord) {
	this.chanRecord <- rec
}

// This should clean up anything lingering about the LogWriter, as it is called before
// the LogWriter is removed.  LogWrite should not be called after Close.
func (this *TelnetLogWriter) Close() {
	this.isStopRun = true

	close(this.chanRecord)

	if nil != this.listenConn {
		fmt.Println("Close listener!")
		this.listenConn.Close()
	}

	if nil != this.remoteConn {
		this.remoteConn.Close()
	}

	this.wg.Wait()
}


//
// Command
//

func dbgHelp(args ...interface{}) {
	c := args[0].(*TelnetLogWriter)

	for name, telnetCmd := range c.myCmds {
		outText:= fmt.Sprintf("%s:%s\n", name, telnetCmd.usage)
		c.remoteConn.Write([]byte(outText))
	}
}

