package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/DavidGamba/go-getoptions"
)

var (
	//Trace : For debug level logging
	Trace *log.Logger
	//Info : For info level logging
	Info *log.Logger
	//Error : For error level logging
	Error *log.Logger
)

// Init : Logging initialization
func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	outputverbosePtr := ioutil.Discard
	outputdebugPtr := ioutil.Discard

	opt := getoptions.New()

	opt.Bool("help", false)
	listenaddrPtr := opt.String("listenport", "0.0.0.0:2055")
	destportPtr := opt.Int("destport", 2055)
	destipPtr := opt.String("destip", "")
	bindportPtr := opt.Int("bindport", 0)
	logverbosePtr := opt.Bool("verbose", false)
	logdebugPtr := opt.Bool("debug", false)

	remaining, err := opt.Parse(os.Args[1:])

	helpText := `Usage example for ./udpreflector
    --listenport 0.0.0.0:2055
    --destip X.X.X.X
    --destport 2055
	--bindport 2000
    --verbose
    --debug
    
    `

	if err != nil {
		fmt.Println(err)
		fmt.Println(helpText)
		os.Exit(1)
	}

	if opt.Called("help") {
		fmt.Println(helpText)
		os.Exit(0)
	}

	if *logdebugPtr {
		outputdebugPtr = os.Stdout
	}

	if *logverbosePtr {
		outputverbosePtr = os.Stdout
	}

	Init(outputdebugPtr, outputverbosePtr, os.Stderr)

	Info.Println("Starting up ...")
	Info.Println("Listen Address", *listenaddrPtr)
	Info.Println("Outbound Port", *bindportPtr)
	Info.Println("Destination Address", *destipPtr)
	Info.Println("Destination Port", *destportPtr)

	if remaining != nil {
		Info.Println("Remaining: ", remaining)
	}

	// listen for incoming udp packets
	inbound, err := net.ListenPacket("udp", *listenaddrPtr)
	if err != nil {
		log.Fatal(err)
	}

	localaddr := net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: *bindportPtr}
	remoteaddr := net.UDPAddr{IP: net.ParseIP(*destipPtr), Port: *destportPtr}
	outbound, err := net.DialUDP("udp", &localaddr, &remoteaddr)

	if err != nil {
		Error.Println(err)
	}

	defer inbound.Close()
	defer outbound.Close()

	buffer := make([]byte, 16384)

	for {
		//simple read

		receivedInt, srcAddr, err := inbound.ReadFrom(buffer)

		if err != nil {
			log.Fatal(err)
		}
		outbound.Write(buffer[0:receivedInt])
		Trace.Println("SRC:", srcAddr, "Received:", buffer)
	}

}
