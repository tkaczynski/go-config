package main

import (
	"fmt"
	"os"

	"github.com/tkaczynski/go-config"
)

func main() {
	cfg, err := go_config.ConfigDefault()
	if err != nil {
		panic(err)
	}

	_ = os.Setenv("MYAPP_MYSETTINGBOOL", "true")
	_ = os.Setenv("MYAPP_MYSETTINGINT", "123")
	_ = os.Setenv("MYAPP_MYSETTINGSTR", "one-two-three")
	
	fmt.Println("MyApp.MySettingBool", cfg.MustBool("MyApp.MySettingBool"))
	fmt.Println("MyApp.MySettingInt", cfg.MustInt("MyApp.MySettingInt"))
	fmt.Println("MyApp.MySettingStr", cfg.MustString("MyApp.MySettingStr"))
}
