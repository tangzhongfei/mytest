package controllers

import (
	"encoding/json"
	"strings"
	"fmt"
	"regexp"

	"github.com/astaxie/beego"
	"gitlab.qiyi.domain/yunwei/ipmi-api/models"
	"gitlab.qiyi.domain/yunwei/ipmi-api/tools"
	"github.com/astaxie/beego/httplib"
)

type SnController struct{
	beego.Controller
}

type ipConf struct{
	private_ip [10]string `json:private_ip`
	public_ip [10]string `json:public_ip`
}


// @Title 上传sn信息，获取其内外网连通性
// @Description 上传sn号， 返回内外网连通性测试结果
// @Param body body []models.Drac "post info"
// @Success 200 {object} []models.Resutl
// @router / [post]
func (this *SnController) Post(){
	snlist := []models.Sn{}
	length := 0
	results := [5000]models.SnResult{}
	url := "https://portal.qiyi.domain/?appkey=box&sn=@"
	if err := json.Unmarshal(this.Ctx.Input.RequestBody,&snlist); err!=nil{
		this.Data["json"] = []models.SnResult{
			models.SnResult{
				Sn: "json prase error",
				InternetPing: "false",
				IntranetPing: "false",
			},
		}
		this.ServeJSON()
	}
	 for count, sn := range snlist{
		ipconf := ipConf{}
		boxurl := strings.Replace(url, "@", sn.Sn, -1)
		req := httplib.Get(boxurl)
		reqstr, _ := req.String()
		public_length := 0
		private_length := 0
		public_result := ""
		private_result := ""

		fmt.Printf("%s\n---\n",reqstr)
		
		//多外网
		r := regexp.MustCompile(`public_ip":\[("[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}",){1,4}"[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}"`)
		params := r.FindStringSubmatch(reqstr)
		if params != nil{
			fmt.Printf("%s\n---\n",params[0])
			splitemp := strings.Split(params[0],"[")
			temp := strings.Split(splitemp[1],",")
			for index,str := range temp{
				ipconf.public_ip[index] = str
				public_length = index+1
			}
		}	else{
			//单外网
			r = regexp.MustCompile(`public_ip":\["[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}"`)
			params = r.FindStringSubmatch(reqstr)
			if params != nil{
				fmt.Printf("%s\n---\n",params[0])
				temp := strings.Split(params[0],"[")
				fmt.Printf("%s\n---\n",temp[1])
				ipconf.public_ip[0] = temp[1]
				public_length = 1
			}	else{
				fmt.Printf("no public ip---\n\n")
				public_length = 0
				public_result = "no public ip"
			}
		}
		//多内网
		r = regexp.MustCompile(`private_ip":\[("[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}",){1,4}"[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}"`)
		params = r.FindStringSubmatch(reqstr)
		if params != nil{
			fmt.Printf("%s\n---\n",params[0])
			splitemp := strings.Split(params[0],"[")
			temp := strings.Split(splitemp[1],",")
			for index,str := range temp{
				ipconf.private_ip[index] = str
				private_length = index+1
			}
		}	else{
			//单内网
			r = regexp.MustCompile(`private_ip":\["[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}"`)
			params = r.FindStringSubmatch(reqstr)
			if params != nil{
				fmt.Printf("%s\n---\n",params[0])
				temp := strings.Split(params[0],"[")
				fmt.Printf("%s\n---\n",temp[1])
				ipconf.private_ip[0] = temp[1]
				private_length = 1
			}	else{
			fmt.Printf("no private ip---\n\n")
			private_length = 0
			private_result = "no private_ip"
			}
		}


		fmt.Printf("%s %s\n",ipconf.private_ip,ipconf.public_ip)

		for i:=0;i<private_length;i++{
			private_info := tools.CheckDracConn(ipconf.private_ip[i])
			if private_info == nil {
				private_result = private_result + ipconf.private_ip[i] + " up  "
			}else{
				private_result = private_result + ipconf.private_ip[i] + " down  "
			}
		}
		for i:=0;i<public_length;i++{
			public_info := tools.CheckDracConn(ipconf.public_ip[i])
			if public_info == nil {
				public_result = public_result + ipconf.public_ip[i] + " up  "
			}else{
				public_result = public_result + ipconf.public_ip[i] + " down  "
			}
		}
		results[count] = models.SnResult{
			Sn: sn.Sn,
			InternetPing: public_result,
			IntranetPing: private_result,
		}
		length++
	
	 }
	 this.Data["json"] = results[0:length]
	 this.ServeJSON()
}