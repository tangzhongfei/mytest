package controllers

import (
	"fmt"
	"strings"
	"encoding/json"

	"gitlab.qiyi.domain/yunwei/ipmi-api/models"
	"gitlab.qiyi.domain/yunwei/ipmi-api/tools"
	"github.com/astaxie/beego"
)




type IpmiController struct{
	beego.Controller
}



// @Title 上传sn，user，password信息
// @Description 上传 Drac 返回连通性结果和验证用户密码结果
// @Param body body  []models.Drac "post info"
// @Success 200 {object} []models.Result
// @router / [post]
func (this *IpmiController) Post(){
	hosts := []models.Drac{}
	check := [1000]int{}

	length := 0
	results := [1000]models.Result{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody,&hosts); err!=nil{
		this.Data["json"] = []models.Result{
			models.Result{
				Sn: "json prase error",
				Pingresult: false,
				Loginresult: false,
			},
		}
		this.ServeJSON()
	}

	//初始化启动线程	
	maxSignal := make(chan int, 20)
	reSignal := make(chan int, 500)

	for index, host := range hosts{
		
		//测试连通性
		maxSignal <- 1
		fmt.Println("begin test!!!!!!")
		go testNet(host, &results[index], maxSignal, reSignal, index)
		length = index + 1
	}

	//这里做读取的检测，防止线程没跑完就返回结果，保证全部处理完返回结果
	for k:=0;; {
		if(k == length) {
			break
		}
		check[<-reSignal] = 1
		if(check[k] == 1){
			k++
			for{
				if(check[k] == 1){
					k++
				}else{
					break
				}
			}
		}
	}

	fmt.Printf("%d \n",length)
	this.Data["json"] = results[0:length]

	fmt.Printf("%v\n",results[0:length])
	this.ServeJSON()
}



func testNet(host models.Drac, res *models.Result, maxSignal chan int, reSignal chan int, index int){
	hosttpl := []string{"ilo%s.qiyi.drac", "idrac-%s.qiyi.drac", "hdm%s.qiyi.drac", "%s.qiyi.drac"}
	hostname := ""
		for _, tpl := range hosttpl{
			hosttest := strings.Replace(tpl, "%s", host.Sn, -1)
			info := tools.CheckDracConn(hosttest)
			fmt.Println("check Conn.")
			if info == nil {
				hostname = hosttest
				break
			}
		}
		if hostname == "" {
			*res = models.Result{
				Sn: host.Sn,
				Pingresult: false,
				Loginresult: false,
			}
			<-maxSignal
			reSignal <-index
			return 
		}

		//登录测试用户和密码
		response, err := tools.CheckDracLoginInfo(hostname, host)
		fmt.Println("check Drac.")
		if err!=nil || response == -1 {
			*res = models.Result{
				Sn: host.Sn,
				Pingresult: true,
				Loginresult: false,
			}
		}else{
			*res = models.Result{
				Sn: host.Sn,
				Pingresult: true,
				Loginresult: true,
			}
		}
		fmt.Printf("%s %t %t \n",res.Sn,res.Pingresult,res.Loginresult)
		<-maxSignal
		reSignal <-index
}