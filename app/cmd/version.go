package cmd

import (
	"context"
	"fmt"

	"github.com/chunkburst/PreUSDT/app"
	"github.com/urfave/cli/v3"
)

var Version = &cli.Command{
	Name:  "version",
	Usage: "显示版本信息",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("preusdt 版本：" + app.Version)

		return nil
	},
}
