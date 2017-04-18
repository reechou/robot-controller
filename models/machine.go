package models

import (
	"fmt"
	"time"
	
	"github.com/reechou/holmes"
)

type Machine struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	Host      string `xorm:"not null default '' varchar(128) unique(machine_host)" json:"host"`
	RobotType int    `xorm:"not null default 0 int unique(machine_host)" json:"robotType"` // 0: just robot 1: robot group manager 2: robot wechat business
	RobotNum  int64  `xorm:"not null default 0 int" json:"robotNum"`
	CreatedAt int64  `xorm:"not null default 0 int" json:"createAt"`
	UpdatedAt int64  `xorm:"not null default 0 int" json:"updateAt"`
}

func CreateRobot(info *Machine) error {
	if info.Host == "" {
		return fmt.Errorf("machine host[%s] cannot be nil.", info.Host)
	}
	
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	
	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create machine error: %v", err)
		return err
	}
	holmes.Info("create machine[%v] success.", info)
	
	return nil
}

func GetMachinesFromType(robotType int) ([]Machine, error) {
	var list []Machine
	err := x.Where("robot_type = ?", robotType).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
