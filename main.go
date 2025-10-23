package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	_ "go.xoder.cn/shortener/internal/bootstrap"
	"go.xoder.cn/shortener/internal/shared"

	"github.com/spf13/viper"

	"go.xoder.cn/shortener/internal/routers"
)

//go:embed version.txt
var version string

var description = `
   _____ _                _
  / ____| |              | |
 | (___ | |__   ___  _ __| |_ ___ _ __   ___ _ __
  \___ \| '_ \ / _ \| '__| __/ _ \ '_ \ / _ \ '__|
  ____) | | | | (_) | |  | ||  __/ | | |  __/ |
 |_____/|_| |_|\___/|_|   \__\___|_| |_|\___|_|

 Shortener: %s
 Project: https://git.jetsung.com/idev/shortener-server
 Website: https://i.jetsung.com
 Author:  Jetsung Chan <i@jetsung.com>

`

func main() {
	// 定义命令行参数
	configFile := flag.String("config", "", "配置文件路径 (例如: -config /path/to/config.toml)")
	flag.Parse()

	// 如果指定了配置文件，则设置环境变量
	if *configFile != "" {
		os.Setenv("SHORTENER_CONFIG_FILE", *configFile)
		fmt.Printf("使用配置文件: %s\n\n", *configFile)
	}

	addr := viper.GetString("server.address")
	if addr == "" {
		addr = ":8080"
	}

	fmt.Printf(description, version)
	fmt.Printf(" api key: %s \n username: %s \n password: %s \n\n",
		shared.GlobalAPIKey,
		shared.GlobalUser.Username,
		shared.GlobalUser.Password,
	)

	r := routers.NewRouter()
	if err := r.Run(addr); err != nil {
		panic("run server failed: " + err.Error())
	}
}
