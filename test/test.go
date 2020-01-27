package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CxZMoE/ncov"
)

const (
	// SERVERADDR 服务器地址
	SERVERADDR = "0.0.0.0:1314"
)

func main() {
	log.Printf("\n>> 新型冠状病毒疫情实时监控系统 <<\n")
	serverStartd := make(chan bool, 1)
	go func() {
		<-serverStartd
		for {
			vm := ncov.GetVirusStatus()
			allAreaStatus := ncov.GetAllAreaStatus()
			fmt.Printf("%s\n", vm.GetString())
			simple := ">> \n"
			for _, area := range allAreaStatus.Areas {
				fmt.Printf("[%s]%s已经确诊感染%.0f人，疑患病%.0f人，目前已死亡%.0f人、治愈%.0f人。\n", vm.RecentTime, area.Area, area.ConfirmCount, area.SuspectCount, area.DeadCount, area.HealCount)
				simple += fmt.Sprintf("[%s]%s已经确诊感染%.0f人，疑患病%.0f人，目前已死亡%.0f人、治愈%.0f人。\n", vm.RecentTime, area.Area, area.ConfirmCount, area.SuspectCount, area.DeadCount, area.HealCount)
			}
			log.Printf(">> DUMP 模式已开启，每两个小时半刷新一次数据并保存文档\n")
			ncov.Dump(string(ncov.GetAllAreaStatus().GetAllAreaStatusJSON()))
			ncov.DumpSimple(simple)
			time.Sleep(time.Hour*2 + time.Minute*30)
		}
	}()

	go func() {
		api := &ncov.API{}
		server := &http.Server{
			Addr:         SERVERADDR,
			Handler:      api,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		fmt.Println(">> 服务器API已开启")
		fmt.Println(">> 地址:", SERVERADDR)
		serverStartd <- true
		server.ListenAndServe()
	}()

	for {
		input, _, _ := bufio.NewReader(os.Stdin).ReadLine()
		if string(input) == ":q" {
			break
		}
	}
}
