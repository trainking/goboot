package main

import (
	"fmt"
	"time"

	"github.com/trainking/goboot/example/game_robot/robot"
)

func main() {
	// for i := 1; i < 3; i++ {
	// 	startRobot(strconv.Itoa(i), "123456")
	// }

	r1 := robot.New("tcp", "192.168.1.16:6001")
	startRobot(r1, "1", "123456")
	r2 := robot.New("tcp", "192.168.1.16:6001")
	startRobot(r2, "2", "123456")

	r2.Say(1, "2 say hello")

	time.Sleep(10 * time.Second)
}

func startRobot(r *robot.Robot, account string, passowd string) {
	if err := r.Login(account, passowd); err != nil {
		fmt.Printf("登录失败：%v\n", err)
		r.Quit()
		return
	}

	fmt.Printf("%s 登录成功\n", account)

	go r.Receive()
}
