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
	AUDIO_FAST              = "fast.wav"
	AUDIO_NORMAL            = "normal.wav"
	AUDIO_SLOW              = "slow.wav"

	CMD_PLAY_CORRECT = "play_correct"
	CMD_PLAY_WRONG   = "play_wrong"
	CMD_PLAY_START   = "play_start"
	CMD_PLAY_END     = "play_end"
	CMD_PLAY_FAST    = "play_fast"
	CMD_PLAY_NORMAL  = "play_normal"
	CMD_PLAY_SLOW    = "play_slow"

	CMD_PLAY_SOUND_PREFIX      = "play_sound_"
	CMD_PLAY_GAME_SOUND_PREFIX = "play_game_"
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
		CMD_PLAY_FAST:    path.Join(AUDIO_GENERAL_DIRECTORY, AUDIO_FAST),
		CMD_PLAY_NORMAL:  path.Join(AUDIO_GENERAL_DIRECTORY, AUDIO_NORMAL),
		CMD_PLAY_SLOW:    path.Join(AUDIO_GENERAL_DIRECTORY, AUDIO_SLOW),
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
			var command = strings.TrimSpace(string(buf[:n]))
			fmt.Printf("command: %s\n", command)
			if strings.HasPrefix(command, CMD_PLAY_SOUND_PREFIX) { // sounds 1-8
				file = path.Join(AUDIO_SOUNDS_DIRECTORY, strings.Trim(command, CMD_PLAY_SOUND_PREFIX)+".wav")
				fmt.Printf("file: %s\n", file)
			} else if strings.HasPrefix(command, CMD_PLAY_GAME_SOUND_PREFIX) { // game sounds 1-4
				file = path.Join(AUDIO_GENERAL_DIRECTORY, strings.Trim(command, CMD_PLAY_GAME_SOUND_PREFIX)+".wav")
				fmt.Printf("file: %s\n", file)
			} else { // general sounds
				file = AUDIO_GENERAL_MAPPING[command]
				fmt.Printf("file: %s\n", file)
				if file == "" { // unknown command received
					break
				}
			}
			_, err = playwav.FromFile(file) // actually play the sound, don't minding potential errors caused by a non existing file, etc.
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}
}
