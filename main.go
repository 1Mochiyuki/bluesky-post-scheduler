package main

import (
	"github.com/1Mochiyuki/gosky/config"
	"github.com/1Mochiyuki/gosky/config/logger"
	"github.com/1Mochiyuki/gosky/db"
	"github.com/1Mochiyuki/gosky/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

var l = logger.Get()

func main() {
	configErr := config.InitConfig()
	if configErr != nil {

		log.Error("config err", configErr)
		return
	}

	if dbErr := db.InitDB(); dbErr != nil {
		l.Error().Err(dbErr).Msg("could not initalize database")
		panic(dbErr)
	}

	p := tea.NewProgram(ui.NewAppEntry())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
