package main

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/aerth/playwav"
	"go.bug.st/serial"
)

const (
	AUDIO_GENERAL_DIRECTORY = "audio/general/"
	AUDIO_SOUNDS_DIRECTORY  = "audio/sounds/"
	AUDIO_CORRECT           = "correct.wav"
	AUDIO_WRONG             = "wrong.wav"
	AUDIO_START             = "start.wav"
	AUDIO_END               = "end.wav"

	CMD_PLAY_CORRECT      = "play_correct\n"
	CMD_PLAY_WRONG        = "play_wrong\n"
	CMD_PLAY_START        = "play_start\n"
	CMD_PLAY_END          = "play_end\n"
	CMD_PLAY_SOUND_PREFIX = "play_sound_"
)

var (
	mode = &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	AUDIO_GENERAL_MAPPING = map[string]string{
		CMD_PLAY_CORRECT: path.Join(AUDIO_GENERAL_DIRECTORY, AUDIO_CORRECT),
		CMD_PLAY_WRONG:   path.Join(AUDIO_GENERAL_DIRECTORY, AUDIO_WRONG),
		CMD_PLAY_START:   path.Join(AUDIO_GENERAL_DIRECTORY, AUDIO_START),
		CMD_PLAY_END:     path.Join(AUDIO_GENERAL_DIRECTORY, AUDIO_END),
	}
)

func main() {
	fmt.Println("starting...")
	for { // use first serial port
		fmt.Println("looking for serial port...")
		time.Sleep(5 * time.Second) // check every 5 seconds

		ports, err := serial.GetPortsList()
		if err != nil {
			continue
		}
		if len(ports) != 1 {
			continue
		}

		// open port
		fmt.Println("opening serial port...")
		port, err := serial.Open(ports[0], mode)
		if err != nil {
			continue
		}
		// read from port
		fmt.Println("reading from serial port...")
		buf := make([]byte, 1024)
		for {
			n, err := port.Read(buf)
			if err != nil {
				break
			}
			var file string
			fmt.Println(string(buf[:n]))
			if strings.HasPrefix(string(buf[:n]), CMD_PLAY_SOUND_PREFIX) { // sounds 1-8
				file = path.Join(AUDIO_SOUNDS_DIRECTORY, strings.TrimPrefix(string(buf[:n]), CMD_PLAY_SOUND_PREFIX)+".wav")
			} else { // general sounds
				file = AUDIO_GENERAL_MAPPING[string(buf[:n])]
				if file == "" { // unknown command received
					break
				}
			}
			fmt.Println(file)
			playwav.FromFile(file) // actually play the sound, don't minding potential errors caused by a non existing file, etc.
		}
	}
}
