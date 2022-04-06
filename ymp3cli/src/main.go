package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/paij0se/lmmp3"
	"github.com/paij0se/ymp3cli/src/cli"
	"github.com/paij0se/ymp3cli/src/cli/scdl/pkg/soundcloud"
	"github.com/paij0se/ymp3cli/src/client"
	"github.com/paij0se/ymp3cli/src/rpc"
	"github.com/paij0se/ymp3cli/src/server"
	"github.com/phayes/freeport"
)

var (
	version       = "0.6.1"
	help          bool
	update        bool
	showVersion   bool
	speed         string
	url           string
	urlSoundcloud string
	spotify       string
	port          int
	song          string
)

func startServer() (err error) {

	if port == 0 {
		port, err = freeport.GetFreePort()
		if err != nil {
			fmt.Println("Error getting free port")
			return
		}
	}

	go rpc.Rpc(port)
	go client.StartClient(fmt.Sprintf(":%d", port))
	server.StartServer(fmt.Sprintf(":%d", port))
	return
}
func init() {
	flag.BoolVar(&help, "help", false, "show help")
	flag.BoolVar(&help, "h", false, "show help")

	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")

	flag.StringVar(&url, "d", "", "download a song from youtube ")
	flag.StringVar(&url, "download", "", "download a song from youtube ")

	flag.StringVar(&speed, "s", "", "that allows changing the playback speed")
	flag.StringVar(&speed, "speed", "", "that allows changing the playback speed")

	flag.StringVar(&urlSoundcloud, "sd", "", "download a song from soundcloud ")
	flag.StringVar(&urlSoundcloud, "soundcloud", "", "download a song from soundcloud ")

	flag.StringVar(&spotify, "sp", "", "download a song from spotify ")
	flag.StringVar(&spotify, "spotify", "", "download a song from spotify ")

	flag.StringVar(&song, "p", "", "play a single song")
	flag.StringVar(&song, "play", "", "play a single song")

	flag.BoolVar(&update, "update", false, "update ymp3cli to the latest version")
	flag.BoolVar(&update, "u", false, "update ymp3cli to the latest version")

	flag.IntVar(&port, "port", 0, "port to run the server on")
	flag.Parse()
}

func main() {
	go cli.Stats()
	// create the folder if it doesn't exist
	if _, err := os.Stat("music"); os.IsNotExist(err) {
		os.Mkdir("music", 0777)

	}
	// This is going to install ffmpeg if is not installed(only on windows)
	lmmp3.DownloadFFmpeg()
	if url != "" {
		lmmp3.DownloadAndConvert(os.Args[2])
		if runtime.GOOS == "windows" {
			del := exec.Command(`cmd`, `/C`, "del", "*.mpeg")
			if del.Run() != nil {
				fmt.Println("Error deleting the mpeg files")
			}
		}
	} else if song != "" {
		cli.PlaySongCli(song)
	} else if speed != "" {
		go rpc.Speedrpc(os.Args[2])
		cli.Speedy(os.Args[2])
	} else if help {
		cli.HelpCommand()
	} else if urlSoundcloud != "" {
		soundcloud.ExtractSong(urlSoundcloud)
	} else if update {
		cli.Update()
	} else if spotify != "" {
		cmd := exec.Command("/bin/bash", "-c", "spotdl "+spotify)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	} else if showVersion {
		fmt.Println(version)
	} else {
		startServer()
	}
}
