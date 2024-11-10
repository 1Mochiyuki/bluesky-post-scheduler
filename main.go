package main

import (
	"github.com/1Mochiyuki/gosky/ui/login"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func main() {
	configErr := initConfig()
	if configErr != nil {
		log.Error("config err", configErr)
		return
	}
	appPass := viper.GetString("gs_ap")
	log.Print("envs", appPass)
	p := tea.NewProgram(login.InitLoginScreenModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
