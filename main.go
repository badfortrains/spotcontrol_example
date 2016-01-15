package main

import (
	"bufio"
	"fmt"
	"github.com/badfortrains/spotcontrol"
	"strings"
	"strconv"
	"flag"
	"os"
)

func chooseDevice(controller *spotcontrol.SpircController, reader *bufio.Reader) string{
	devices := controller.ListDevices()
	fmt.Println("\n choose a device:")
	for i, d := range devices {
		fmt.Printf("%v) %v %v \n", i, d.Name, d.Ident)
	}
	
	for {
		fmt.Print("Enter device number: ")
		text, _ := reader.ReadString('\n')
		i, err := strconv.Atoi(strings.TrimSpace(text))
		if err == nil && i < len(devices) && i >= 0{
			return devices[i].Ident
		}
		fmt.Println("invalid device number")

	}
}

func getDevice(controller *spotcontrol.SpircController, ident string, reader *bufio.Reader) string{
	if ident != "" {
		return ident
	} else {
		return chooseDevice(controller, reader)
	}
}

func printHelp(){
	fmt.Println("\nAvailable commands:")
	fmt.Println("load <track1> [...more tracks]: load tracks by spotify base 62 id")
	fmt.Println("hello:                          ask devices to identify themselves")
	fmt.Println("play:                           play current track")
	fmt.Println("pause:                          pause playing track")
	fmt.Println("devices:                        list availbale devices")
	fmt.Println("help:                           show this list\n")
}


func main() {
	username := flag.String("username", "", "spotify username")
	password := flag.String("password", "", "spotify password")
	appkey := flag.String("appkey", "./spotify_appkey.key", "spotify appkey file path")
	flag.Parse()

	if *username == "" || *password == "" {
		fmt.Println("need to supply a username and password")
		fmt.Println("./spirccontroller --username SPOTIFY_USERNAME --password SPOTIFY_PASSWORD")
		return
	}

	s := spotcontrol.Session{}
	s.StartConnection()
	s.Login(*username, *password, *appkey)
	s.Run()


	//fmt.Println(convert62("3Vn9oCZbdI1EMO7jxdz2Rc 2nMW1mZmdIt5rZCsX1uh9J"))

	sController := spotcontrol.SetupController(&s, *username, "7288edd0fc3ffcbe93a0cf06e3568e28521687bc")
	
	go sController.Run()
	sController.SendHello()

	reader := bufio.NewReader(os.Stdin)
	var ident string
	printHelp()
	for {
		fmt.Print("Enter a command: ")
		text, _ := reader.ReadString('\n')
		cmds := strings.Split(strings.TrimSpace(text),  " ")

		switch {
		case cmds[0] == "load":
			ident = getDevice(&sController, ident, reader)
			sController.LoadTrack(ident, cmds[1:])
		case cmds[0] == "hello":
			sController.SendHello()
		case cmds[0] == "play":
			ident = getDevice(&sController, ident, reader)
			sController.SendPlay(ident)
		case cmds[0] == "pause":
			ident = getDevice(&sController, ident, reader)
			sController.SendPause(ident)
		case cmds[0] == "devices":
			ident = chooseDevice(&sController, reader)
		case cmds[0] == "help":
			printHelp()
		}
	}

}