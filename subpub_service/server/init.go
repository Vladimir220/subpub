package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitSystem() (url string, infoLog, errLog *log.Logger) {

	err := godotenv.Load("./server/.env")
	if err != nil {
		panic(err)
	}

	errLogFile, err := os.OpenFile("./server/Log/err.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	errLog = log.New(errLogFile, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)

	infoLogFile, err := os.OpenFile("./server/Log/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	infoLog = log.New(infoLogFile, "INFO:", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)

	url = os.Getenv("SERVER_URL")

	return
}
