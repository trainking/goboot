package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/trainking/goboot/example/game_robot/robot"
)

func main() {
	for i := 1; i < 3; i++ {
		startRobot(strconv.Itoa(i), "123456")
	}

	time.Sleep(10 * time.Second)
}

func startRobot(account string, passowd string) {
	r := robot.New("kcp", "127.0.0.1:6001")
	if err := r.Login(account, passowd); err != nil {
		fmt.Printf("登录失败：%v\n", err)
		r.Quit()
		return
	}

	fmt.Printf("%s 登录成功\n", account)

	go r.Receive()

	if account != "1" {
		r.Say(1, "Hello")
	}

}
