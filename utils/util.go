package utils

import (
	"fmt"
	"log"
	"runtime"
)

var Logger *log.Logger

//에러 체크 후 에러 시 무조건 패닉 발생하는 함수 -> error발생 시 작동 X여서 다시 실행해야 될 때 사용
func IfErrorMakePanic(err error, errorInfo string) {
	if err != nil {
		log.Println(errorInfo)
		panic(err)
	}
}

func ErrorCheck(err error) {
	if err != nil {
		PrintErrFunc(err, 2)
	}
}

//error check시마다 print하는 것을 함수로 생성
func PrintErrFunc(err error, depth int) {
	if err == nil {
		return
	}
	funcName, file, line, ok := runtime.Caller(depth)
	logText := fmt.Sprintf("[%s] [%s] (%d line) - %s",
		file, runtime.FuncForPC(funcName).Name(), line, err.Error())
	if ok && Logger != nil {
		Logger.Printf(logText)
	} else if ok {
		log.Printf(logText)
	}
}
