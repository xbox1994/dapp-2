package main

import (
	"bytes"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
)

const url = "https://game-api.wanakafarm.com/api/v1/"

func main() {
	username := "282503846@qq.com"
	password := "11111111"
	t := login(username, password)
	if t == "" {
		panic("error account")
	}
	items, lands := userInfo(t)
	// 有几块地
	for _, land := range lands {
		for _, item := range items {
			// 收获，还需要加上判断哪些地可以收获，再调用
			harvestedItem(land.Int(), item.Int(), t)
			// 播种：1. 先看d.list_lands.cellList的田里有哪些地可以种，能种哪种菜
			//		2. 然后在filter_item里搜索可以种的对应的菜
			//		3. 执行播种操作
			// 包菜有一步挖田？？待解决
			// 浇水，还需要加上判断哪些地可以浇水，再调用
			water(item.Int(), t)
		}
	}
}

func login(username, password string) string {
	u := url + "user/connect_wallet"
	var jsonStr = `{
		"c": {},	
		"d": {
			"wallet_address": "` + username + `",
			"sign": "` + password + `"
		}
	}`
	r := post(u, jsonStr, "")
	t := gjson.Get(r, "d.access_token")
	if t.Str != "" {
		return t.Str
	}
	return ""
}

func userInfo(t string) (items, lands []gjson.Result) {
	u := url + "user/user_info"
	var jsonStr = `{
		"c": {},
		"d": {}
	}`
	r := post(u, jsonStr, t)
	return gjson.Get(r, "d.list_lands.0.list_items.#.id").Array(),
		gjson.Get(r, "d.list_lands.#.id").Array()
}

func water(itemId int64, t string) []gjson.Result {
	u := url + "action/watering"
	var jsonStr = `{
	"c": {},
	"d": {
		"itemId": ` + strconv.FormatInt(itemId, 10) + `
	}
}`
	fmt.Println("浇水：", itemId)
	r := post(u, jsonStr, t)
	s := gjson.Get(r, "d.list_lands.#.id")
	return s.Array()
}

func harvestedItem(landId, itemId int64, t string) {
	u := url + "inventory/harvested_item"
	var jsonStr = `{
	"c": {},
	"d": {
		"itemId": ` + strconv.FormatInt(itemId, 10) + `,
		"landId": ` + strconv.FormatInt(landId, 10) + `
	}
}`
	fmt.Println("收获：", itemId)
	post(u, jsonStr, t)
}

func filterItem(itemId float64, t string) []gjson.Result {
	u := url + "inventory/filter_item"
	itemIdS := strconv.FormatFloat(itemId, 'E', -1, 32)
	var jsonStr = `{
	"c": {},
	"d": {
		"type": "Tree"
	}
}`
	fmt.Println("树：", itemIdS)
	r := post(u, jsonStr, t)
	s := gjson.Get(r, "d.list_lands.#.id")
	return s.Array()
}

func growingItem(itemId float64, t string) []gjson.Result {
	u := url + "inventory/growing_item"
	itemIdS := strconv.FormatFloat(itemId, 'E', -1, 32)
	var jsonStr = `{
	"c": {},
	"d": {
		"itemId": ` + itemIdS + `,
		"landId": ` + itemIdS + `,
		"posX": ` + itemIdS + `,
		"posY": ` + itemIdS + `
	}
}`
	fmt.Println("播种：", itemIdS)
	r := post(u, jsonStr, t)
	s := gjson.Get(r, "d.list_lands.#.id")
	return s.Array()
}

func post(url, jsonStr, t string) string {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Authorization", t)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	s := string(body)
	fmt.Println("response Body:", s)
	return s
}
