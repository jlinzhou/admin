package confs

import (
	"admin/pkg/logs"
	"encoding/json"
	"fmt"
	"time"
	"github.com/go-ini/ini"
)

type ServerConf struct {
	RunMode      string
	HTTPPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout time.Duration
	TokenExpired int
}
type MysqlConf struct {
	User         string
	Password     string
	Host         string
	DataBase     string
	Debug        bool
	MaxOpenConns int
	MaxLifetime  int
	MaxIdleConns int
}
type LogConf struct {
	Filename string
	Level     string
	Maxsize   int
	Maxdays   int
}
type Conf struct{
	ServerConf
	MysqlConf
	LogConf
}
var (
	AdminConf Conf
	cfg    *ini.File
)

func loadServer() (err error) {
	sec, err := cfg.GetSection("server")
	if err != nil {
		err = fmt.Errorf("Fail to get section 'server': %v", err)
		return
	}
	//取setting值发生错误就取MustString
	runMode := sec.Key("mode").MustString("")
	httpPort := sec.Key("port").MustString("")
	readTimeout := time.Duration(sec.Key("readTimeout").MustInt()) * time.Second
	writeTimeout := time.Duration(sec.Key("writeTimeout").MustInt()) * time.Second
	idleTimeout:= time.Duration(sec.Key("idleTimeout").MustInt()) * time.Second
	tokenExpired := sec.Key("tokenExpired").MustInt()
	if len(runMode) == 0 {
		err = fmt.Errorf("init server failed,runMode  is null")
		return
	}
	AdminConf.ServerConf.RunMode = runMode
	AdminConf.ServerConf.HTTPPort= httpPort
	AdminConf.ServerConf.ReadTimeout = readTimeout
	AdminConf.ServerConf.WriteTimeout = writeTimeout
	AdminConf.ServerConf.IdleTimeout = idleTimeout
	AdminConf.ServerConf.TokenExpired = tokenExpired
	return
}

func loadmysql() (err error) {
	sec, err := cfg.GetSection("mysql")
	if err != nil {
		err = fmt.Errorf("Fail to get section 'static': %v", err)
		return
	}
	AdminConf.MysqlConf.User = sec.Key("user").MustString("")
	AdminConf.MysqlConf.Password = sec.Key("password").MustString("")
	AdminConf.MysqlConf.Host = sec.Key("host").MustString("")
	AdminConf.MysqlConf.DataBase = sec.Key("database").MustString("")
	AdminConf.MysqlConf.Debug = sec.Key("debug").MustBool()
	AdminConf.MysqlConf.MaxOpenConns = sec.Key("maxOpenConns").MustInt()
	AdminConf.MysqlConf.MaxLifetime = sec.Key("maxLifetime").MustInt()
	AdminConf.MysqlConf.MaxIdleConns = sec.Key("maxIdleConns").MustInt()

	return
}


func loadLog() (err error) {
	sec, err := cfg.GetSection("log")
	if err != nil {
		err = fmt.Errorf("Fail to get section 'log': %v", err)
		return
	}
	filename := sec.Key("filename").MustString("")
	level := sec.Key("level").MustString("")
	maxsize := sec.Key("maxsize").MustInt() //相当于500kb
	maxdays := sec.Key("maxdays").MustInt()
	if len(filename) == 0  || len(level) == 0 {
		err = fmt.Errorf("init log failed,log_filename or log_level or ...  is null")
		return
	}
	AdminConf.LogConf.Filename = filename
	AdminConf.LogConf.Level = level
	AdminConf.LogConf.Maxsize = maxsize
	AdminConf.LogConf.Maxdays = maxdays

	return
}

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}
func initLog() {
	config := make(map[string]interface{})
	config["filename"] = AdminConf.LogConf.Filename
	config["level"] = convertLogLevel(AdminConf.LogConf.Level)
	config["maxsize"] = AdminConf.LogConf.Maxsize
	config["maxdays"] = AdminConf.LogConf.Maxdays

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed,err:", err)
		return
	}

	_=logs.SetLogger(logs.AdapterConsole)
	_=logs.SetLogger(logs.AdapterFile, string(configStr))
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	return
}
func LoadConfig()(err error) {
	cfg, err = ini.Load("./confs/app.ini")
	if err != nil {
		return
	}
	err = loadServer()
	if err != nil {
		return
	}
	err = loadLog()
	if err != nil {
		return
	}
	err = loadmysql()
	if err != nil {
		return
	}
	initLog()
	return

}
