package main

import (
	"os"
	"bufio"
	"fmt"
	"github.com/badfortrains/spotcontrol"
	"strings"
	"strconv"
)

func chooseDevice(controller *spotcontrol.SpircController, reader *bufio.Reader) string{
	devices := controller.ListDevices()
	for i, d := range devices {
		fmt.Printf("%v %v %v \n", i, d.Name, d.Ident)
	}
	for {
		text, _ := reader.ReadString('\n')
		i, err := strconv.Atoi(strings.TrimSpace(text))
		if err == nil && i < len(devices) && i >= 0{
			return devices[i].Ident
		}
	}
}

func getDevice(controller *spotcontrol.SpircController, ident string, reader *bufio.Reader) string{
	if ident != "" {
		return ident
	} else {
		return chooseDevice(controller, reader)
	}
}

func main() {
	s := spotcontrol.Session{}
	s.StartConnection()
	s.Login()
	s.Run()

	username := os.Getenv("SPOT_USERNAME")
	sController := spotcontrol.SetupController(&s, username, "7288edd0fc3ffcbe93a0cf06e3568e28521687bc")
	
	go sController.Run()

	reader := bufio.NewReader(os.Stdin)
	var ident string
	for {
		fmt.Print("Enter a command: ")
		text, _ := reader.ReadString('\n')
		switch {
		case strings.TrimSpace(text) == "hello":
			sController.SendHello()
		case strings.TrimSpace(text) == "play":
			ident = getDevice(&sController, ident, reader)
			sController.SendPlay(ident)
		case strings.TrimSpace(text) == "pause":
			ident = getDevice(&sController, ident, reader)
			sController.SendPause(ident)
		case strings.TrimSpace(text) == "device":
			ident = chooseDevice(&sController, reader)
		}
	}

}