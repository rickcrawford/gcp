package commands

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

func init() {
	// set config defaults
	viper.SetDefault("garbage-collect", false)
	viper.SetConfigType("yml")
	viper.SetEnvPrefix("ta")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// flags
	RootCmd.PersistentFlags().StringP("conf-dir", "", "./conf", "Configuration Directory")
	RootCmd.PersistentFlags().StringP("database", "", "bolt.db", "BoltDB file")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug mode")
	RootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet mode. Do not display banner messages")

	// config
	viper.BindPFlag("conf-dir", RootCmd.PersistentFlags().Lookup("conf-dir"))
	viper.BindPFlag("database", RootCmd.PersistentFlags().Lookup("database"))
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("quiet", RootCmd.PersistentFlags().Lookup("quiet"))

	// local flags;
	RootCmd.Flags().StringVar(&configFile, "config", "", "/path/to/config.yml")

}

// https://medium.com/@skdomino/writing-better-clis-one-snake-at-a-time-d22e50e60056

// RootCmd is the main command to run the application
var RootCmd = &cobra.Command{
	Use:   "autocomplete.boltdb",
	Short: "Boltdb Engine",

	// parse the config if one is provided, or use the defaults. Set the backend
	// driver to be used
	PersistentPreRun: preRun,

	// run
	Run: run,
}

func preRun(ccmd *cobra.Command, args []string) {
	// if --config is passed, attempt to parse the config file
	if configFile != "" {
		// get the filepath
		abs, err := filepath.Abs(configFile)
		if err != nil {
			log.Fatal("Error reading filepath", err)
		}

		// get the config name
		base := filepath.Base(abs)

		// get the path
		path := filepath.Dir(abs)

		//
		viper.SetConfigName(strings.Split(base, ".")[0])
		viper.AddConfigPath(path)

		// Find and read the config file; Handle errors reading the config file
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Failed to read config file", err)
			os.Exit(1)
		}
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Println("Config file changed", e.Name)
			// Send a HUP signal to restart
			if p, err := os.FindProcess(os.Getpid()); err == nil {
				p.Signal(os.Interrupt)
			}
		})
	}
}

func run(cmd *cobra.Command, args []string) {
	// track signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	for {
		if exit := start(sig); exit {
			return
		}
	}

}
