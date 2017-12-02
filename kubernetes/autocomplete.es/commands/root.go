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
	viper.SetConfigType("yml")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// flags
	RootCmd.PersistentFlags().StringP("conf-dir", "", "./conf", "Configuration Directory")
	RootCmd.PersistentFlags().StringP("elastic-url", "", "http://localhost:9200", "Elastic Search URL")
	RootCmd.PersistentFlags().StringP("elastic-login", "", "elastic", "Elastic Search Login")
	RootCmd.PersistentFlags().StringP("elastic-password", "", "changeme", "Elastic Search Password")
	RootCmd.PersistentFlags().StringP("elastic-index-name", "", "autocomplete", "Elastic Search Index Name")
	RootCmd.PersistentFlags().StringP("redis-url", "", "redis://localhost:6379", "Redis URL")
	RootCmd.PersistentFlags().StringP("project-id", "", "typeahead-183622", "Project ID")
	RootCmd.PersistentFlags().StringP("topic-name", "", "updates", "PubSub Topic Name")
	RootCmd.PersistentFlags().StringP("subscription-name", "", "autocomplete.es", "PubSub Subscription Name")
	RootCmd.PersistentFlags().IntP("http-port", "", 8080, "HTTP port")
	RootCmd.PersistentFlags().IntP("https-port", "", 8443, "HTTPS port")
	RootCmd.PersistentFlags().BoolP("debug", "", false, "Enable debugging")

	// config
	viper.BindPFlag("conf-dir", RootCmd.PersistentFlags().Lookup("conf-dir"))
	viper.BindPFlag("elastic-url", RootCmd.PersistentFlags().Lookup("elastic-url"))
	viper.BindPFlag("elastic-login", RootCmd.PersistentFlags().Lookup("elastic-login"))
	viper.BindPFlag("elastic-password", RootCmd.PersistentFlags().Lookup("elastic-password"))
	viper.BindPFlag("elastic-index-name", RootCmd.PersistentFlags().Lookup("elastic-index-name"))
	viper.BindPFlag("redis-url", RootCmd.PersistentFlags().Lookup("redis-url"))
	viper.BindPFlag("project-id", RootCmd.PersistentFlags().Lookup("project-id"))
	viper.BindPFlag("topic-name", RootCmd.PersistentFlags().Lookup("topic-name"))
	viper.BindPFlag("subscription-name", RootCmd.PersistentFlags().Lookup("subscription-name"))
	viper.BindPFlag("http-port", RootCmd.PersistentFlags().Lookup("http-port"))
	viper.BindPFlag("https-port", RootCmd.PersistentFlags().Lookup("https-port"))
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	// local flags;
	RootCmd.Flags().StringVar(&configFile, "config", "", "/path/to/config.yml")

}

// https://medium.com/@skdomino/writing-better-clis-one-snake-at-a-time-d22e50e60056

// RootCmd is the main command to run the application
var RootCmd = &cobra.Command{
	Use:   "esautocomplete",
	Short: "auto complete using elastic",
	Run:   run,

	// parse the config if one is provided, or use the defaults. Set the backend
	// driver to be used
	PersistentPreRun: preRun,
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

func preRun(ccmd *cobra.Command, args []string) {
	// if --config is passed, attempt to parse the config file
	if configFile != "" {
		// get the filepath
		abs, err := filepath.Abs(configFile)
		if err != nil {
			log.Fatal("Error reading filepath", err)
			os.Exit(1)
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
