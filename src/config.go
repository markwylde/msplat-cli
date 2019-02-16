package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

// ReadConfig : Get config
func ReadConfig(config string) {
	viper.SetConfigType("yaml")

	configFileData, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Print("There was an error reading the config.yml file.\n\nPlease make sure you have specified either the --config flag, or the MSPLAT_CONFIG environment variable.\n\n")
		log.Fatal(err)
		os.Exit(1)
	}

	viper.SetDefault("paths.stacks", path.Join(filepath.Dir(config), "stacks"))

	viper.ReadConfig(bytes.NewBuffer(configFileData))
}
