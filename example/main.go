package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/FDUTCH/ez-scoreboard"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	conf, err := readConfig(slog.Default())
	if err != nil {
		panic(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	index := 0

	sc := ExampleScoreBoard{
		Greet: "Hi!",
		Name: func(p *player.Player) string {
			return "Nick: " + p.Name()
		},
		Online: func(p *player.Player) string {
			return fmt.Sprint("your index: ", index)
		},
	}

	updater, err := scoreboard.BuildUpdater(sc, "name")
	if err != nil {
		panic(err)
	}

	srv.Listen()
	for p := range srv.Accept() {
		updater.Update(p)
		index++
	}
}

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}

type ExampleScoreBoard struct {
	Greet  scoreboard.StaticLine
	Name   scoreboard.DynamicLine
	Space1 scoreboard.Space
	Space2 scoreboard.Space
	Online scoreboard.DynamicLine
}
