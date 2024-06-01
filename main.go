package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	exmpl2()
}

/***
*	Second example 2 : Message communication with For Select
***/

type Message struct {
	from    string
	payload string
}

type Server struct {
	msgCh  chan Message
	quitCh chan struct{}
}

func (s *Server) StartAndListen() {
free:
	for {
		select {
		case msg := <-s.msgCh:
			fmt.Printf("Message from %s payload %s\n", msg.from, msg.payload)
		case <-s.quitCh:
			fmt.Println("Doing some cleaning up and updates, before shutdown")
			break free
		}
	}

	fmt.Println("server is down")
}

func sendMessage(user string, message string, msgCh chan Message) {
	msg := Message{
		from:    user,
		payload: message,
	}
	msgCh <- msg
}

func shutdownGracefully(qCh chan struct{}) {
	qCh <- struct{}{}
}

func exmpl2() {

	server := &Server{
		msgCh:  make(chan Message),
		quitCh: make(chan struct{}),
	}

	go server.StartAndListen()

	for i := 0; i < 5; i++ {
		go func(i int) {
			time.Sleep(2 * time.Second)
			sendMessage(fmt.Sprintf("User %d", i), fmt.Sprintf("Message %d", i), server.msgCh)

		}(i)
	}

	go func() {
		time.Sleep(5 * time.Second)
		shutdownGracefully(server.quitCh)
	}()

	select {}
}

/***
*	First example 1 : How ?
***/
func exmpl1() {
	now := time.Now()
	userID := 10
	respChan := make(chan string, 128)
	wg := &sync.WaitGroup{}

	go fetchUserData(userID, respChan, wg)
	go fetchUserRecommandations(userID, respChan, wg)
	go fetchUserLikes(userID, respChan, wg)

	wg.Add(3)
	wg.Wait()

	close(respChan)

	for r := range respChan {
		fmt.Println(r)
	}

	fmt.Println(time.Since(now))
}

func fetchUserData(userID int, respChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("... fetching user data")
	time.Sleep(8 * time.Millisecond)
	respChan <- fmt.Sprintf("user %d data", userID)
	fmt.Println("... end fetching user data")
}

func fetchUserRecommandations(userID int, respChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("... fetching user recommandations")
	time.Sleep(10 * time.Millisecond)
	respChan <- fmt.Sprintf("user %d recommandations", userID)
	fmt.Println("... end fetching user recommandations")
}

func fetchUserLikes(userID int, respChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("... fetching user likes")
	time.Sleep(3 * time.Millisecond)
	respChan <- fmt.Sprintf("user %d likes", userID)
	fmt.Println("... end fetching user likes")
}
