package main
//noinspection GoUnresolvedReference
import (
	"GOTcpService/test"
	"fmt"
	"net"
	"os"
	"time"
)

func handleConnection(conn net.Conn) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	t1 := time.Now() // get current time
	defer func() {
		cost := time.Since(t1)
		fmt.Println("cost=", cost)
	}()
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = test.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)

	}
}

func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			//log4go.Info(string(data))
			Log(string(data))
		}
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {

	netListen, err := net.Listen("tcp", ":9988")
	CheckError(err)

	defer netListen.Close()

	Log("Waiting for clients")

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		Log(conn.RemoteAddr().String(), " tcp connect success")

		go handleConnection(conn)
	}

}

