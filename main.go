// file_trnas project main.go
package main

import (
	"PackFrame"
	"bufio"

	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	MAX_CONN_NUM = 100000
)

func main() {
	Server()
	//fmt.Println("create a server(s) or client(c)?")
	//reader := bufio.NewReader(os.Stdin)
	//input, _, _ := reader.ReadLine()
	//if string(input) == "s" {
	//	Server()
	//}
	//if string(input) == "c" {
	//	Client()
	//} else {
	//	fmt.Println("err arguments,entering again!.\r\n  alternaltive argument is server or client")
	//	os.Exit(0)
	//}
}
func Bar(vl int, width int) string {
	return fmt.Sprintf("%s%*c", strings.Repeat("█", vl/10), vl/10-width+1,
		([]rune(" ▏▎▍▌▋▋▊▉█"))[vl%10])
}

//func Show(s string) string {
//	enc := mahonia.NewEncoder("gbk") //中文转码有错误的函数。
//	return enc.ConvertString(s)
//}

///////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////

func Server() {
	//exit := make(chan bool)
	//ip := net.ParseIP("127.0.0.1")
	//addr := net.TCPAddr{ip, 8888, ""}
	go func() {

		logrus.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	///////////////////////////////////
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("running ...\n")

	//////////////////////////////////////////////////////

	//var cur_conn_num int = 0
		conn_chan := make(chan net.Conn)
	//	ch_conn_change := make(chan int)
	//
	//	go func() {
	//		for conn_change := range ch_conn_change {
	//			cur_conn_num += conn_change
	//		}
	//	}()

	//go func() {
	//	for _ = range time.Tick(10e9) {
	//		fmt.Printf("cur conn num: %d\n", cur_conn_num)
	//		//debug.FreeOSMemory()
	//		//runtime.GC()
	//	}
	//}()

	//////////////////////////////////////////////////
	//for i := 0; i < MAX_CONN_NUM; i++ {
	//	go func() {
	//		for conn := range conn_chan {
	//			ch_conn_change <- 1
	//			EchoFunc(conn)
	//			ch_conn_change <- -1
	//		}
	//	}()
	//}

	//////////////////////////////////////////////////
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				println("Error accept:", err.Error())
				return
			}

			conn_chan <- conn
		}
	}()

	for {
		conn:=<- conn_chan

			//fmt.Println(conn.RemoteAddr())
			go EchoFunc(conn)
			//defer conn.Close()

	}
}

// orgin server///////////////////////////////////
func EchoFunc(tcpcon net.Conn) {

	fmt.Println("Client connect", tcpcon.RemoteAddr())
	defer tcpcon.Close()
	//defer runtime.GC()
	//exit := make(chan bool)

	//recv file
	var pts_last uint16
	var fo *os.File
	var file_len int64 = 1024
	var ctr int = 0
	//	var t time.Time
	//	t = time.Now()
	var data []byte
	data = make([]byte, 16384)

	var pt PackFrame.PackTag
	var rec []byte
	tcpcon.SetDeadline(time.Unix(5,0),)
	for {
		n, err := tcpcon.Read(data)

		if err != io.EOF && err != nil {
			return
			//rr:=reflect.ValueOf(err).Elem().FieldByName("Err").Interface()
			//ff:=reflect.ValueOf(rr).Elem().FieldByName("Err")

			//switch  ff.Interface().(error) {
			//
			//	case net.ErrWriteToConnected:
			//		fallthrough
			//	case io.ErrUnexpectedEOF:
			//		fallthrough
			//
			//	case syscall.ENETRESET:
			//		fmt.Println("\nclient closed", )
			//		return
			//	case syscall.ECONNABORTED:
			//	case syscall.ECONNRESET:
			//	case syscall.Errno(10054):
			//		fmt.Println("\n10054 closed",tcpcon.RemoteAddr() )
			//		return
			//	default:
			//		//fmt.Printf("%v \n%+v \n%q\n %t\n\n", err, err, err, err)
			//		return
			//
			//}
		} //if err != io.EOF && err != nil

		if n == 0 {
			time.Sleep(time.Microsecond * 10)
		}

		for n > 0 {
			//fmt.Println(n, err, data[:n])
			//var ttt float32
			//tt = time.Now().Sub(t)
			//if tt != 0 {
			//	ttt = float32(time.Second) / float32(tt)
			//	fmt.Printf("%d \t/ %d \t  %02f %% %d %f kb\r", ctr, file_len/1024, 100*float32(ctr)/float32(file_len/1024), tt, ttt)
			//}
			//t = time.Now()
			//i := 100 * int64(ctr) / (file_len / 1024)
			//fmt.Printf("\f%s \t%d%% \r", Bar(int(i), 25), i)
			fmt.Println(tcpcon.RemoteAddr(),n)
			ctr++
			_L_next_head:
			if data[0] == 0x55 && data[1] == 0xAA {
				pt, rec, err = PackFrame.Depack(data[:12+int(data[2])+int(data[3])*256])
				if err != nil {

				}
				tmp := 13 + int(data[2]) + int(data[3])*256
				data = data[tmp-1:]
				n = n - tmp

			} else {
				//fmt.Println("rust data :", data[:n])
				//fmt.Printf("\r\nrust data :%s from: %s \r\n",data[:n] ,tcpcon.RemoteAddr())
				for i, chk := range data[:n] {
					//fmt.Println(i,chk)
					if chk == 0x55 && data[i+1] == 0xAA {
						data = data[i:]
						n = n - i
						fmt.Println("seek head next")
						goto _L_next_head
					}

				}
				//fmt.Println("junk all")
				n = 0
				data=nil
				continue
			}


			if pt.Pserial == pts_last {
				//what? the same pack?
			}
			pts_last = pt.Pserial

			switch pt.Pcmd {

			case fc_filehead:
				if pt.Ppara == fcp_fileName {
					fmt.Println(getCurrentDirectory())
					fo, err = os.Create(getCurrentDirectory() + "//rec//" + string(rec))
					//						fo, err = os.Create( "e://"  + string(rec))
					defer fo.Close()
				} else if pt.Ppara == fcp_fileSize {
					bb := bytes.NewBuffer(rec)
					binary.Read(bb, binary.LittleEndian, &file_len)
					//PackFrame.ByteToType(rec,file_len)
				} else if pt.Ppara == fcp_fileEOF {

					//exit <- true
					break
				} else {
					panic("how did reach the empty filehead?? ")
				}
			case fc_filebody:
				//write to the file
				jn, err := fo.Write(rec)
				fmt.Printf("%d \t/ %d \r", ctr, file_len/1024)
				jn = jn
				err = err
			default:

			}
		}

	}

	//<-exit
}

//frame cmd type list
const (
	fc_filehead = 0x10
	fc_filebody = 0x11
)

// fc file paras
const (
	fcp_fileName = 0x01
	fcp_fileEOF  = 0x02
	fcp_fileSize = 0x03
)

func Client() {
	//
	//open file
	fmt.Println("send ur file to the destination", "input ur filename:")
	reader := bufio.NewReader(os.Stdin)
	input, _, _ := reader.ReadLine()
	fmt.Println(string(input))
	fi, err := os.Open(string(input))
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fiinfo, err := fi.Stat()
	fmt.Println("the size of file is ", fiinfo.Size(), "bytes") //fiinfo.Size() return int64 type

	//to online
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("connect server fail！", err.Error())
		return
	}
	defer conn.Close()

	//send file name
	ready, err := PackFrame.Dopack([]byte(fiinfo.Name()),
		fc_filehead, fcp_fileName)

	fmt.Printf("%s", ready)
	_, err = conn.Write(ready)
	if err != nil {
		fmt.Println("conn.Write", err.Error())
	}
	//time.Sleep(time.Microsecond * 5)

	//send file size
	ready, err = PackFrame.Dopack(PackFrame.TypeToByte(fiinfo.Size()),
		fc_filehead, fcp_fileSize)

	_, err = conn.Write(ready)
	//_, err = conn.Write([]byte(string(fiinfo.Size())))
	if err != nil {
		fmt.Println("conn.Write", err.Error())
	}

	var ctr uint32 = 0
	for {
		//fmt.Println("No ", ctr, "of total ", fiinfo.Size()/1024)
		buff := make([]byte, 1024)
		n, err := fi.Read(buff)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			ready, err = PackFrame.Dopack(buff[:n], fc_filehead, fcp_fileEOF)
			_, err = conn.Write(ready)
			//conn.Write([]byte("filerecvend"))
			fmt.Println("filerecvend")
			break
		}
		//fmt.Printf("cnt%d ", ctr, buff[:n])
		ready, err = PackFrame.Dopack(buff[:n], fc_filebody, ctr)
		_, err = conn.Write(ready)
		//		_, err = conn.Write(buff)
		if err != nil {
			fmt.Println(err.Error())
		}
		ctr++
		//time.Sleep(time.Microsecond * 10)
	}
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
