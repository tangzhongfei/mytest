package tools

import (
	"fmt"
	"strings"
	"os/exec"
	"unicode"
	"bytes"

	"gitlab.qiyi.domain/yunwei/ipmi-api/models"
	"github.com/astaxie/beego"
)

type Info struct{
	Code string
	Msg string
}


type IpmiController struct{
	beego.Controller
}

// CheckDracConn 检查Drac连通性
func CheckDracConn(hostname string) *Info{
	beego.Error("Ping Begin")
	cmd := fmt.Sprintf("ping -c 4 %s | grep 'packet loss' | awk '{print $4}'", hostname)
	stdout, err := ExecCommand(cmd)
	if err != nil {
		return &Info{
			Code: "-1",
			Msg:  err.Error(),
		}
	}
	rs := strings.FieldsFunc(stdout, unicode.IsSpace)
	fmt.Printf("%v %v\n",rs,len(rs))
	if len(rs) == 0 || rs[0] != "4" {
		return &Info{
			Code: "-1",
			Msg:  "ping unreachable.",
		}
	}
	beego.Error("Ping End")
	return nil
}

// CheckDracLoginInfo 检查Drac用户名密码
func CheckDracLoginInfo(hostname string, host models.Drac) (int, error){
	cmd := fmt.Sprintf("ipmitool -I lanplus -H %s -U %s -P '%s' user summary", hostname, host.User, host.Password)
	beego.Informational(cmd)
	output, err := ExecCommand(cmd)
	if err != nil {
		return -1, err
	}
	if strings.Contains(output, "Error") {
		return -1, nil
	}
	return 1, nil
}


// ExecCommand 运行本地命令
func ExecCommand1(commands string) (string, error) {
	out, err := exec.Command("bash", "-c", commands).Output()
	fmt.Printf("stdut: %s\n",out)
	return string(out), err
}


// ExecCommand 运行本地命令
func ExecCommand(commands string) (string, error) {
	cmd := exec.Command("bash", "-c", commands)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		outStr, _ := string(stdout.Bytes()), string(stderr.Bytes())
		return outStr, err
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:%s err:%s\n", outStr, errStr)
	return outStr, nil
}