package main

import (
	"context"
	"os"

	"github.com/chunkburst/PreUSDT/app/cmd"
	"github.com/chunkburst/PreUSDT/app/conf"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func init() {
	// 不推荐引导小白参与修改各种配置文件
	_ = godotenv.Load()
}

func main() {
	c := &cli.Command{
		Name:  "preusdt",
		Usage: conf.Desc,
		Commands: []*cli.Command{
			cmd.Start,
			cmd.Version,
			cmd.Reset,
		},
	}
	if err := c.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
