package main

import "fmt"

func main() {
	url, infoLog, errLog := InitSystem()

	fmt.Println("Start server", url)
	myServ := CreateServer(errLog)
	defer myServ.Close()

	listening(myServ, url, infoLog, errLog)
}
