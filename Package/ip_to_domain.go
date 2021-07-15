package ip_to_domain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MobileInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float32 `json:lat`
	Lon         float32 `json:lon`
	Timezone    string  `json:"tmiezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func ip_enter(ip_get string) string { //正则获取输入ip
	r, _ := regexp.Compile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}")
	if r.MatchString(ip_get) {
		return r.FindString(ip_get)
	}
	return "error"
}

func http_get(url string) string {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*2) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 2)) //设置发送接受数据超时
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 2,
		},
	}
	reqest, err := http.NewRequest("GET", url, nil) //请求
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.89 Safari/537.36")
	if err != nil {
		return "err"
	}
	resp, err := client.Do(reqest)
	if err != nil {
		return "err"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "err"
	}
	return string(body)
}

func ip_api_com(ip_get string) (string, string, string) { //从ip_api_com获取ip归属
	url := "http://ip-api.com/json/" + ip_get + "?lang=zh-CN"
	if http_get(url) == "err" {
		return "", "", ""
	}
	jsonStr := http_get(url) //解析api传回的内容
	var mobile MobileInfo
	err := json.Unmarshal([]byte(jsonStr), &mobile) //解析josn
	if err != nil {
		return "", "", ""
	}
	w := mobile.Country + mobile.RegionName + mobile.City
	return string(w), string(mobile.Isp), string(mobile.Org) //
}

func site_ip138_com(ip_get string, c chan []string) []string { //从site.ip138.com获取数据
	slist := []string{}
	url := "https://site.ip138.com/" + string(ip_get) + "/"
	if http_get(url) == "err" {
		c <- slist
		return slist
	}
	r := regexp.MustCompile("</span><a href=\"(.*?)\" target=\"_blank\">(.*?)</a></li>") //正则匹配特征
	urls := r.FindAllStringSubmatch(http_get(url), -1)
	for _, param := range urls {
		if !strings.Contains(param[2], ".") {
			continue
		}
		slist = append(slist, param[2])
	}
	c <- slist
	return slist
}

func dns_aizhan_com(ip_get string, c chan []string) []string { //从dns.aizhan.com获取数据
	slist := []string{}
	url := "https://dns.aizhan.com/" + string(ip_get) + "/"
	if http_get(url) == "err" {
		c <- slist
		return slist
	}
	r := regexp.MustCompile("<a href=\"(.*?)\" rel=\"nofollow\" target=\"_blank\">(.*?)</a>")
	urls := r.FindAllStringSubmatch(http_get(url), -1)
	for _, param := range urls {
		if !strings.Contains(param[2], ".") {
			continue
		}
		slist = append(slist, param[2])
	}
	c <- slist
	return slist
}

func www_chaxunle_cn(ip_get string, c chan []string) []string { //从www.chaxunle.cn获取数据
	slist := []string{}
	url := "https://www.chaxunle.cn/ip/" + string(ip_get) + ".html"
	if http_get(url) == "err" {
		c <- slist
		return slist
	}
	r := regexp.MustCompile("<a class=\"tip\" href=\"(.*?)\" target=\"_blank\" title=\"(.*?)\">(.*?)</a>")
	urls := r.FindAllStringSubmatch(http_get(url), -1)

	for _, param := range urls {
		if !strings.Contains(param[2], ".") {
			continue
		}
		slist = append(slist, param[2])
	}
	c <- slist
	return slist
}

func ip_yqie_com(ip_get string, c chan []string) []string { //从ip.yqie.com获取数据
	slist := []string{}
	url := "http://ip.yqie.com/iptodomain.aspx?ip=" + string(ip_get)
	if http_get(url) == "err" {
		c <- slist
		return slist
	}
	r := regexp.MustCompile("<td width=\"90%\" class=\"blue t_l\" style=\"text-align: center\">(.*?)</td>")
	urls := r.FindAllStringSubmatch(http_get(url), -1)
	for _, param := range urls {
		if param[1] == "域名" {
			continue
		}
		if !strings.Contains(param[1], ".") {
			continue
		}
		slist = append(slist, param[1])
	}
	c <- slist
	return slist
}

func RemoveRepeatedElement(arr []string) (newArr []string) { //结果整合去重
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

func if_survive(url_rr []string, port int) []string { //判断域名是否存活，可指定端口
	var wg sync.WaitGroup
	var mutex sync.Mutex
	params := make([]string, 0)
	for _, param := range url_rr {
		wg.Add(1)
		go func(param string, port int) {
			defer wg.Done()
			tcp_ip := fmt.Sprintf("%s:%d", param, port)
			conn, err := net.DialTimeout("tcp", tcp_ip, time.Second)
			if err != nil {
			} else {
				conn.Close()
				mutex.Lock()
				params = append(params, param)
				mutex.Unlock()
			}
		}(param, port)
	}
	wg.Wait()
	return params
}

func to_all(ip string) []string { //多线程调用获取数据函数
	c := make(chan []string)
	//c := make(chan []string)
	go dns_aizhan_com(ip, c)
	go site_ip138_com(ip, c)
	go www_chaxunle_cn(ip, c)
	go ip_yqie_com(ip, c)
	//time.Sleep(2)
	l1, l2, l3, l4 := <-c, <-c, <-c, <-c //对应删除则减少数据接口
	//单一接口出错，不影响程序运行
	urls := append(append(l3, (l1)...), append(l2, (l4)...)...)
	return RemoveRepeatedElement(urls)
	/*
		若不需要4个数据接口，则只需要注释 go 函数
		并且修改字符串拼接
	*/
}

func parameter_processing(args ...interface{}) []string { //参数处理
	var ip string = ""         //ip为默认空
	var number int = 30        //默认输出30条
	var if_s bool = false      //默认不对域名存活检测
	var port int = 80          //默认探测80端口
	for i, arg := range args { //对输入参数进行处理
		if i == 0 { //第一位参数只能为ip
			if ip_enter(string(fmt.Sprintf("%s", arg))) == "error" {
				fmt.Println("No IP input detected, please enter")
				os.Exit(0)
			} else {
				ip = string(fmt.Sprintf("%s", arg))
				continue
			}
		}
		switch arg.(type) { //arg.(type) 获取参数的类型
		case string: //是否判断域名存活
			str := string(fmt.Sprintf("%s", arg))
			if str == "ture" { //是否存活检测
				if_s = true
			}
			if str[0:1] == "n" { //截取参数第一位为n+数量
				number, _ = strconv.Atoi(str[1:])
			}
			if str[0:1] == "p" { //截取参数第一位为p+端口
				port, _ = strconv.Atoi(str[1:])
			}
		default:
			fmt.Println("Input error, input parameter can only be string")
			os.Exit(0)
		}
	}
	urls := to_all(ip) //数据整合获取
	if if_s {
		urls = if_survive(urls, port)
	}

	if len(urls) != number { //限制数据输出数量
		if number > len(urls) {
			return urls
		}
		if number < len(urls) {
			urls1 := urls[:number]
			return urls1
		}
	}
	return urls //返回整理好的域名
}

func values_put(ip string, S1 string, S2 string, S3 string, urls []string) string { //输出处理
	type Stu struct {
		Ip        string   `json:ip`
		Location  string   `json:location`
		Isp       string   `json:icp`
		Org       string   `json:org`
		Url_lists []string `json:urls`
	}
	stu := Stu{
		Ip:        ip,
		Location:  S1,
		Isp:       S2,
		Org:       S3,
		Url_lists: urls,
	}
	jsonStu, _ := json.Marshal(stu)
	return string(jsonStu)
}

func Ip_domain(args ...string) string { //主函数
	if len(args) < 1 {
		return "Enter at least one parameter IP"
		os.Exit(0)
	}
	if len(args) > 4 {
		return "Enter up to four parameters"
		os.Exit(0)
	}
	ip := ip_enter(string(args[0]))
	if ip == "error" {
		return "IP input error"
		os.Exit(0)
	}
	if len(args) == 1 {
		S1, S2, S3 := ip_api_com(ip)
		urls := parameter_processing(ip)
		return values_put(ip, S1, S2, S3, urls)
	}
	if len(args) == 2 {
		S1, S2, S3 := ip_api_com(ip)
		urls := parameter_processing(ip, args[1])
		return values_put(ip, S1, S2, S3, urls)

	}
	if len(args) == 3 {
		S1, S2, S3 := ip_api_com(ip)
		urls := parameter_processing(ip, args[1], args[2])
		return values_put(ip, S1, S2, S3, urls)

	}
	if len(args) == 4 {
		S1, S2, S3 := ip_api_com(ip)
		urls := parameter_processing(ip, args[1], args[2], args[3])
		return values_put(ip, S1, S2, S3, urls)
	}
	return "ip no find domain"
}
