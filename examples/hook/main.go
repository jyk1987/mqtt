package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/logrusorgru/aurora"
	"github.com/snple/mqtt"
	"github.com/snple/mqtt/listener"
	"github.com/snple/mqtt/packets"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	fmt.Println(aurora.Magenta("MQTT Server initializing..."), aurora.Cyan("TCP"))

	server := mqtt.New()

	server.SetHook(&MyHook{})

	tcp := listener.NewTCP("t1", ":1883", &mqtt.AuthAllow{})

	err := server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println(aurora.BgMagenta("  Started!  "))

	<-done
	fmt.Println(aurora.BgRed("  Caught Signal  "))

	server.Close()
	fmt.Println(aurora.BgGreen("  Finished  "))
}

type MyHook struct {
}

var _ mqtt.Hook = (*MyHook)(nil)

func (*MyHook) Accept(*mqtt.Server, *mqtt.Client) bool {
	return true
}

func (*MyHook) Remove(*mqtt.Server, *mqtt.Client, error) {}

func (*MyHook) Recv(*mqtt.Server, *mqtt.Client, *packets.Packet) bool {
	return true
}

func (*MyHook) Send(*mqtt.Server, *mqtt.Client, *packets.Packet) bool {
	return true
}

func (*MyHook) Emit(*mqtt.Server, *mqtt.Client, *packets.Packet) bool {
	return true
}

func (*MyHook) Push(*mqtt.Server, *mqtt.Client, *packets.Packet) bool {
	return true
}
