package internal

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"time"
)

var Config *viper.Viper

// TODO : Bind to cli flags ?

func init() {
	v := viper.New()

	// Default value
	v.SetDefault("inventory.ignoredDirs", []string{"template", "common_vars"})
	v.SetDefault("inventory.hostFiles", "*.ini")
	v.SetDefault("inventory.location", "inventories")

	// Location
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	//viper.AddConfigPath("$HOME/config.yml")
	v.AddConfigPath("./")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = v.WriteConfig()
			log.Println("could not found config file created default")
		} else {
			log.Fatal("error while reading config file", err.Error())
		}
	}

	Config = v

	// Watch update
	go func() {
		for {
			time.Sleep(time.Second * 5)
			v.WatchConfig()
			v.OnConfigChange(func(e fsnotify.Event) {
				log.Println("config file updated", e.Name)
			})
		}
	}()

}
