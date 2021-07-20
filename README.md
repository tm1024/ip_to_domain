## 使用方法

### 1、目录结构

```TEXT
├── README.md
├── imgaes                     //使用截图      
│	└── windwos_11.png
│	└── linux_ubuntu.png
│	└── macos_15.7.png
├── Executable_file             //可执行文件
│	└── ip_to_domain_windows.exe
│	└── ip_to_domain
├── Source_code                //源码
│	└── ip_to_domain.go
├── Package                    //自定义包
│   └── ip_to_domain.go
备注： Source_code中的源码与Package中的不同
```



### 2、使用

####  包调用

```go
package main
import (
	"fmt"
	"ip_to_domain.go"
)
func main() {
	a := ip_to_domain.Ip_domain(ip)//仅域名反查
	//a := ip_to_domain.Ip_domain(ip,"n10")//最大输出10条域名数据
	//a := ip_to_domain.Ip_domain(ip,"true","p443")//确定进行存活检测,并指定检测443端口
	//a := ip_to_domain.Ip_domain(ip,true)//确定进行存活检测，支持正则
	//a := ip_to_domain.Ip_domain(ip,"true","p8080","n50")//确定进行存活检测，指定检测8080端口，最大输出50条域名数据
	fmt.Println(a)
}
/*
输出格式例：
{
    "Ip":"1.1.1.1",
    "Location":"",
    "Isp":"",
    "Org":"",
    "Url_lists":[
        "www.978bb.con",
        "6h58.con",
        "www.1314.con",
        "www.qunsf.con",
        "www.345.c",
        "www.tszy.com.c",
        "www25777.con",
        "kejidiy.c",
        "www.3344a.con",
        "www.80pipi.con",
        "www.0800encoder.com",
        "222mimi.con",
        "www.9ttt.con",
        "www.q2002.con",
        "www.191t.con",
        "ggwhcxzg.c",
        "www.zzg.gov.c",
        "wwwjigongxuexiao.con",
        "www.325dd.con",
        "www.055i.con"
    ]
}*/
```

#### 源码直接使用

```TEXT
使用方法:
	go run "ip_to_domain.go" ip
可选：
	是否检测域名存活
	go run "ip_to_domain.go" ip ture
	检测端口，默认80。例如检测443端口
	go run "ip_to_domain.go" ip ture p443
	输出数量，默认30，修改例如输出5
	go run "ip_to_domain.go" ip ture p443 n5
```

#### Mac/Windows/Linux可执行文件

```text

基本使用方法
	./ip_to_domain_macos ip  
	./ip_to_domain_linux ip     
	.\ip_to_domain_windows.exe ip  
可选：
	是否检测域名存活
	./ip_to_domain ip true
	检测端口，默认80。例如检测443端口
	./ip_to_domain ip p443
	输出数量，默认30条域名，修改例如输出5条
	./ip_to_domain ip n5


可选参数无顺序，但第一个必须为ip
```

