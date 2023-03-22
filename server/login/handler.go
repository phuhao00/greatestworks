package main

import "net/http"

func RandomName() {
	//todo 随机名字
}

func Register() {
	preRegister()
	registerReward()
}

func preRegister() {
	//todo 预注册

}

func registerReward() {

}

func Login(w http.ResponseWriter, r *http.Request) {
	//todo check whitelist
	//todo return token

}

func ThirdPartyLogin() {

}

func GetGateWay() {

}

func GetWorldServers() {

}

func HealthyCheck() {

}
