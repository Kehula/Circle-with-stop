package main

import (
	"bufio"
	"fmt"
	"os"
)
var done = make(chan struct{})

type unreadChannel struct {
	ch chan string
}
var closingHandler = make(chan unreadChannel)
var chIn = make(chan string)
var chOut = make(chan string)
func main() {
	var done1 = make(chan string)
	fmt.Println("first send test message")
	
	go resend(0, chIn, chOut)
	go resend(1, chOut, chIn)
	go func() {
		scanner := bufio.NewReader(os.Stdin)
		if n, err := scanner.ReadByte(); n > 0 && err == nil {
			close(done)
			unreadChan := <- closingHandler
			//watching what channel we have to read value before close it
			if unreadChan.ch == chIn {
				<-chIn
			} else if unreadChan.ch == chOut {
				<-chOut
			}
			close(done1)
		}
	}()
	chIn <- "test message" //starting circle work
	<-done1
	<-closingHandler
	close(closingHandler)
	close(chOut)
	close(chIn)
}

func resend(num int, in chan string, out chan string) {
	loop:
		for {
			select {
			case <-done:
				break loop
			default:
				message := <-in
				if message != "" {
					fmt.Printf("%d: %s\n", num, message)
					out <- message
				} else {
					fmt.Println(num, "input closed")
				}
			}
		}
	fmt.Println(num, " is stopped!")
	closingHandler <- unreadChannel{in} //sending unread channel
}