package main

import (
	"math/rand"
	"os"
	"tc-server/pkg/rtc"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/liuhailove/tc-base-go/protocol/logger"
)

// baseFlags --config=config.yaml 从yaml加载配置文件
var baseFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:  "bind",
		Usage: "IP监听的地址，如果这个flag使用多次，则可以绑定多个地址",
	},
	&cli.StringFlag{
		Name:  "config",
		Usage: "TCLive配置文件路径",
	},
	&cli.StringFlag{
		Name:    "config-body",
		Usage:   "YAML文件中的tclive配置，典型的是在容器中通过环境变量传入",
		EnvVars: []string{"TCLIVE_CONFIG"},
	},
	&cli.StringFlag{
		Name:  "key-file",
		Usage: "包含API keys/secrets的文件路径",
	},
	&cli.StringFlag{
		Name:    "keys",
		Usage:   "api keys(key: secret)",
		EnvVars: []string{"TCLIVE_KEYS"},
	},
	&cli.StringFlag{
		Name:    "region",
		Usage:   "当前node所在地区. 被区域意识的节点选择器使用",
		EnvVars: []string{"TCLIVE_REGION"},
	},
	&cli.StringFlag{
		Name:    "node-ip",
		Usage:   "当前节点的IP地址，用于通告给客户端。默认自动确定",
		EnvVars: []string{"NODE_IP"},
	},
	&cli.IntFlag{
		Name:    "udp-port",
		Usage:   "用于 WebRTC 流量的单个 UDP 端口",
		EnvVars: []string{"UDP_PORT"},
	},
	&cli.StringFlag{
		Name:    "redis-host",
		Usage:   "redis server host （incl. port)",
		EnvVars: []string{"REDIS_HOST"},
	},
	&cli.StringFlag{
		Name:    "redis-password",
		Usage:   "redis密码",
		EnvVars: []string{"REDIS_PASSWORD"},
	},
	&cli.StringFlag{
		Name:    "turn-cert",
		Usage:   "tls cert file for TURN SERVER",
		EnvVars: []string{"TCLIVE_TURN_CERT"},
	},
	&cli.StringFlag{
		Name:    "turn-key",
		Usage:   "tls key file for TURN SERVER",
		EnvVars: []string{"TCLIVE_TURN_KEY"},
	},
	// 调试flags
	&cli.StringFlag{
		Name:  "memprofile",
		Usage: "将内存的profile写入`文件`",
	},
	&cli.BoolFlag{
		Name:  "dev",
		Usage: "sets log-level to debug, console formatter, and /debug/pprof. insecure for production",
	},
	&cli.BoolFlag{
		Name:   "disable-strict-config",
		Usage:  "禁用严格的配置解析",
		Hidden: true,
	},
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	defer func() {
		if rtc.Recover(logger.GetLogger()) != nil {
			os.Exit(1)
		}
	}()

	generatedFlags, err := config.GenerateCLIFlags(baseFlags, true)
}
