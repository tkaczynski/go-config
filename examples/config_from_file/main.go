package main

import (
	"fmt"
	
	"github.com/tkaczynski/go-config"
)

func main() {
	cfg, err := go_config.ConfigFromFile("my_settings.yaml")
	if err != nil {
		panic(err)
	}
	
	fmt.Println("MyApp.MySettingBool", cfg.MustBool("MyApp.MySettingBool"))
	fmt.Println("MyApp.MySettingInt", cfg.MustInt("MyApp.MySettingInt"))
	fmt.Println("MyApp.MySettingStr", cfg.MustString("MyApp.MySettingStr"))
	
	fmt.Println("MyApp.MySettingBool as str", cfg.MustString("MyApp.MySettingBool"))
	fmt.Println("MyApp.MySettingInt as str", cfg.MustString("MyApp.MySettingInt"))
}
