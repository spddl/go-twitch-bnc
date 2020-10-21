package main

import (
	"log"
	"os"
	"path"
)

func saveLogFile(logs []byte, channel, logDir string) {
	logfile, err := os.OpenFile(path.Join(logDir, channel+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	_, err = logfile.Write(append([]byte{13, 10}, logs...))
	if err != nil {
		log.Println(err)
	}

	logfile.Close()
}
