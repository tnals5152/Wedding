package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"wedding/api"
	"wedding/config"
	"wedding/utils"
)

var recoverCount uint16 = 2

func main() {
	fmt.Println("Wedding Project")
	err := godotenv.Load("./.env") //.env파일에 DB정보 및 각종 정보 저장하여 처음에 load
	utils.ErrorCheck(err)

	fpLog, err := os.OpenFile(os.Getenv("LogFile"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) //log file 생성
	utils.ErrorCheck(err)
	defer fpLog.Close()

	utils.Logger = log.New(fpLog, "error: ", log.Ldate|log.Ltime) //log file 출력 정보 세팅

	for {
		if recoverCount > 10 { //mainSetting이 10번 죽을 때까지 재실행
			break
		}
		mainSetting()
	}
}

func mainSetting() {
	defer func() { //recover
		if r := recover(); r != nil {
			utils.ErrorCheck(errors.New("main recover" + strconv.Itoa(int(recoverCount))))
			recoverCount++
		}
	}()

	config.ConnectDB() //DB연결
	r := gin.Default()

	r.LoadHTMLFiles("index.html")

	r.GET("/index", api.GetDefaultComment)
	r.GET("/admin-index", api.GetAdminComment)

	//api
	r.Static("/image", "./assets/image")
	http.Handle("/", http.FileServer(http.Dir("static")))
	// http.HandleFunc("/ws", socket.SocketHandler)

	r.GET("/ping", func(c *gin.Context) { //ping 테스트를 위한 코드
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/comment", api.SetComments)
	r.GET("/comment", api.GetComment)

	// go func() {
	// 	//웹소켓 포트 연결(고루틴 없이 작성 시 web_port 열리지 않아 go루틴으로 생성)
	// 	log.Fatal(http.ListenAndServe(":"+os.Getenv("WEBSOCKET_PORT"), nil))
	// }()

	r.Run(":" + os.Getenv("WEB_PORT"))
}
