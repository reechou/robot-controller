package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-controller/config"
	"github.com/reechou/robot-controller/ext"
	"github.com/reechou/robot-controller/models"
)

type Logic struct {
	sync.Mutex

	robotExt *ext.RobotExt
	
	cfg *config.Config
}

func NewLogic(cfg *config.Config) *Logic {
	l := &Logic{
		cfg: cfg,
	}
	l.robotExt = ext.NewRobotExt(cfg)
	
	l.init()

	return l
}

func (self *Logic) init() {
	models.InitDB(self.cfg)
	
	http.HandleFunc("/robot/receive_msg", self.RobotReceiveMsg)
	
	http.HandleFunc("/manager/create_robot_machine", self.CreateRobotMachine)
	http.HandleFunc("/manager/login_robot", self.LoginRobot)
	http.HandleFunc("/manager/get_robot", self.GetRobot)
	http.HandleFunc("/manager/get_robots_from_type", self.GetRobotListFromType)
	http.HandleFunc("/manager/get_login_robots_from_type", self.GetLoginRobotListFromType)
}

func (self *Logic) Run() {
	defer holmes.Start(holmes.LogFilePath("./log"),
		holmes.EveryDay,
		holmes.AlsoStdout,
		holmes.DebugLevel).Stop()

	if self.cfg.Debug {
		EnableDebug()
	}

	holmes.Info("server starting on[%s]..", self.cfg.Host)
	holmes.Infoln(http.ListenAndServe(self.cfg.Host, nil))
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func WriteBytes(w http.ResponseWriter, code int, v []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(v)
}

func EnableDebug() {

}
