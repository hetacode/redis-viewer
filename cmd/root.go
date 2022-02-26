package cmd

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/SaltFishPr/redis-viewer/internal/config"
	"github.com/SaltFishPr/redis-viewer/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "redis-viewer",
	Short: "view redis data in terminal.",
	Long:  `Redis Viewer is a tool to view redis data in terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig(cfgFile)
		cfg := config.GetConfig()

		var rdb *redis.Client
		switch cfg.Mode {
		case "sentinel":
			rdb = redis.NewFailoverClient(
				&redis.FailoverOptions{
					MasterName:    cfg.MasterName,
					SentinelAddrs: cfg.SentinelAddrs,
					Password:      cfg.Password,
				},
			)
		default:
			rdb = redis.NewClient(
				&redis.Options{
					Addr:     cfg.Addr,
					Password: cfg.Password,
					DB:       cfg.DB,
				},
			)
		}

		_, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal("connect to redis failed: ", err)
		}

		p := tea.NewProgram(tui.New(rdb), tea.WithAltScreen(), tea.WithMouseCellMotion())
		if err := p.Start(); err != nil {
			log.Fatal("start failed: ", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	configInfo := "config file (default is $HOME/.config/redis-viewer/redis-viewer.yaml)"
	if runtime.GOOS == "windows" {
		configInfo = "config file (default is $HOME/.redis-viewer.yaml)"
	}
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", configInfo)
}
