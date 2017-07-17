package main 

import (
	"github.com/rnpridgeon/zendb/zendesk"
	"fmt"
	"encoding/json"
	"os"
)

func main() {
	var zdConf zendesk.Config
	conf, _ := os.Open("./exclude/conf.json")
	json.NewDecoder(conf).Decode(&zdConf)

	zd, _ := zendesk.Open(&zdConf)
	fmt.Println(zd.Groups())
	fmt.Println(zd.Organizations())
	fmt.Println(zd.Users())
}

