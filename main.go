package main

import (
	"github.com/reechou/robot-controller/config"
	"github.com/reechou/robot-controller/controller"
)

func main() {
	controller.NewLogic(config.NewConfig()).Run()
}
