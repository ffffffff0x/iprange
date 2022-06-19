package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/projectdiscovery/mapcidr"
	"github.com/thinkeridea/go-extend/exnet"
	"net"
	"os"
	"strings"
)

var filename = flag.String("in", "", "输入文件名")
var cidrSlice []string
var iprangeSlice []string
var wrongSlice []string
var ipv4Left = ""
var ipv4Right = ""
var numSlice []uint

func main() {

	flag.Parse()
	if *filename == "" {
		*filename = "文件名不可为空!!!"
	} else {
		todo()
	}

}

func todo() {

	file, err := os.OpenFile(*filename, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("打开文件失败!", err)
		return
	}
	defer file.Close()

	buf1 := bufio.NewScanner(file)
	for buf1.Scan() {
		// cird
		check1 := strings.Contains(buf1.Text(), "/")

		// ip-range
		check2 := strings.Contains(buf1.Text(), "-")

		if check1 {
			//fmt.Println(buf1.Text(), "cird")
			cidrSlice = append(cidrSlice, buf1.Text())
		} else if check2 {
			//fmt.Println(buf1.Text(), "ip-range")
			iprangeSlice = append(iprangeSlice, buf1.Text())
		} else {
			// 非预期格式
			//fmt.Println(buf1.Text(), "wrong")
			wrongSlice = append(wrongSlice, buf1.Text())
		}
	}

	// 处理 ip-range 队列
	iprange()

	// 处理 cidr 队列
	cidr()

	// 处理非预期格式
	if len(wrongSlice) > 0 {
		fmt.Println("-----存在非预期格式-----")
		wrong()
	}

}

func iprange() {

	for _, v := range iprangeSlice {
		arr := strings.Split(v, "-")
		if len(arr) == 2 {
			ipv4Left = arr[0]
			ipv4Right = arr[1]
			// 左边地址
			addressRight := net.ParseIP(ipv4Left)
			// 右边地址
			addressRight2 := net.ParseIP(ipv4Right)
			if addressRight == nil {
				// 格式非预期
				wrongSlice = append(wrongSlice, v)
			} else if addressRight2 == nil {
				// 格式非预期
				wrongSlice = append(wrongSlice, v)
			} else {
				// 转10进制
				num1, _ := exnet.IPString2Long(ipv4Left)
				num2, _ := exnet.IPString2Long(ipv4Right)
				if num1 > num2 {
					// 格式非预期
					wrongSlice = append(wrongSlice, v)
				} else {
					// 循环自增
					for {
						// 添加进切片
						numSlice = append(numSlice, num1)
						if num1 == num2 {
							break
						}
						num1 += 1
					}
				}
			}
		} else {
			// 格式非预期
			wrongSlice = append(wrongSlice, v)
		}
	}

	for _, v := range numSlice {
		s, _ := exnet.Long2IPString(v)
		fmt.Println(s)
	}

}

func cidr() {

	for _, v := range cidrSlice {
		ips, _ := mapcidr.IPAddresses(v)
		for _, ip := range ips {
			fmt.Println(ip)
		}
	}

}

func wrong() {

	for _, v := range wrongSlice {
		fmt.Println(v)
	}

}
