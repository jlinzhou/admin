package main

import (
	"admin/confs"
	"admin/models"
	"admin/pkg/logs"
)

func main(){
	//initWeb()

	if err:=confs.LoadConfig();err!=nil{
		logs.Error(err)
		return
	}
	if err:=models.InitMysql();err!=nil{
		logs.Error(err)
		return
	}
	go initWeb()
	select{}
}