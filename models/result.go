package models

type Result struct{
	Sn string `json:"sn"`
	Pingresult bool `json:pingresult`
	Loginresult bool `json:loginresult`
}


type SnResult struct{
	Sn string `json:"sn"`
	InternetPing string `json:"internetPing"`
	IntranetPing string `json:"IntranetPing"`
}