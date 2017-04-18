package controller

import (
	"fmt"
	"encoding/json"
	"net/http"
	
	"github.com/reechou/holmes"
	"github.com/reechou/robot-controller/proto"
	"github.com/reechou/robot-controller/models"
	"github.com/reechou/robot-controller/robot_proto"
)

func (self *Logic) CreateRobotMachine(w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		WriteJSON(w, http.StatusOK, rsp)
	}()
	
	if r.Method != "POST" {
		return
	}
	
	req := &models.Machine{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CreateRobotMachine json decode error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	err := models.CreateRobot(req)
	if err != nil {
		holmes.Error("models create machine error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *Logic) LoginRobot(w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		WriteJSON(w, http.StatusOK, rsp)
	}()
	
	if r.Method != "POST" {
		return
	}
	
	req := &robot_proto.StartWxReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("LoginRobot json decode error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	machineList, err := models.GetMachinesFromType(req.RobotType)
	if err != nil {
		holmes.Error("get machines error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	for _, v := range machineList {
		robotCount, err := models.GetRobotsFromTypeHostCount(req.RobotType, v.Host)
		if err != nil {
			holmes.Error("get robots count error: %v", err)
			rsp.Code = proto.RESPONSE_ERR
			continue
		}
		if v.RobotNum > robotCount {
			loginData, err := self.robotExt.LoginRobot(v.Host, req)
			if err != nil {
				holmes.Error("login robot error: %v", err)
				rsp.Code = proto.RESPONSE_ERR
				return
			}
			holmes.Info("robot login in host[%s] success", v.Host)
			rsp.Data = loginData
			return
		}
	}
}

func (self *Logic) GetRobot(w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		WriteJSON(w, http.StatusOK, rsp)
	}()
	
	if r.Method != "POST" {
		return
	}
	
	req := &proto.GetRobotReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobot json decode error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	robot := &models.Robot{
		RobotWx: req.RobotWx,
	}
	has, err := models.GetRobot(robot)
	if err != nil {
		holmes.Error("models get robot error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	if !has {
		rsp.Code = proto.RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("has none of this robot[%s]", req.RobotWx)
		return
	}
	
	rsp.Data = robot
}

func (self *Logic) GetRobotListFromType(w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		WriteJSON(w, http.StatusOK, rsp)
	}()
	
	if r.Method != "POST" {
		return
	}
	
	req := &proto.GetRobotListFromTypeReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetRobotListFromType json decode error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	list, err := models.GetRobotsFromType(req.RobotType)
	if err != nil {
		holmes.Error("models get robot list from type error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	rsp.Data = list
}

func (self *Logic) GetLoginRobotListFromType(w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		WriteJSON(w, http.StatusOK, rsp)
	}()
	
	if r.Method != "POST" {
		return
	}
	
	req := &proto.GetLoginRobotListFromTypeReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("GetLoginRobotListFromType json decode error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	hostList, err := models.GetRobotHostsFromType(req.RobotType)
	if err != nil {
		holmes.Error("models get robot list from type error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	var list []interface{}
	for _, v := range hostList {
		host := fmt.Sprintf("%s%s", v.Ip, v.OfPort)
		robotReq := &robot_proto.RobotGetLoginsReq{
			RobotType: req.RobotType,
		}
		result, err := self.robotExt.GetLoginRobotsFromType(host, robotReq)
		if err != nil {
			holmes.Error("get login robots from type error: %v", err)
			continue
		}
		robots := result.([]interface{})
		if robots == nil {
			holmes.Error("result error: %v", result)
			continue
		}
		list = append(list, robots)
	}
	rsp.Data = list
}
