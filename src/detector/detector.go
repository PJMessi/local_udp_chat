package detector

// import (
// 	"fmt"
// 	"net"
// 	"time"
// )
//
// const SEARCH_ADDRESS = "0.0.0.0:34354"
//
// func NewFinder(clientsChan chan clientstracker.ClientChanVal) *Finder {
// 	return &Finder{
// 		clientsChan: clientsChan,
// 	}
// }
//
// type Finder struct {
// 	clientsChan chan clientstracker.ClientChanVal
// }
//
// func (f *Finder) Search() error {
// 	fmt.Println("\nSearching for clients, please wait...")
//
// 	addr, err := net.ResolveUDPAddr("udp", SEARCH_ADDRESS)
// 	if err != nil {
// 		return fmt.Errorf("err resolving search addr: %s", err)
// 	}
//
// 	listener, err := net.ListenUDP("udp", addr)
// 	if err != nil {
// 		return fmt.Errorf("err creating broadcast listener: %s", err)
// 	}
//
// 	startTime := time.Now()
//
// 	for {
// 		eclapsed := time.Since(startTime)
// 		if eclapsed > 5*time.Second {
// 			break
// 		}
//
// 		var msgBuf []byte = make([]byte, 1024)
// 		n, clientAddr, err := listener.ReadFromUDP(msgBuf)
// 		if err != nil {
// 			return fmt.Errorf("err reading broadcast msg: %s", err)
// 		}
//
// 		username := string(msgBuf[:n])
//
// 		f.clientsChan <- clientstracker.ClientChanVal{
// 			Username:   username,
// 			ClientAddr: clientAddr,
// 		}
// 	}
//
// 	return nil
// }
//
// // BROADCASTER
