package conf

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type ConfigStruct struct {
	BotToken        string   `toml:"bot-token"`
	AdminUsername   []string `toml:"admin-username"`
	ApiKey          string   `toml:"api-key"`
	NewUserBalance  int      `toml:"new-user-balance"`
	SignGiftBalance int      `toml:"sign-gift-balance"`
	ImaginePrice    int      `toml:"imagine-price"`
}

var Conf ConfigStruct

func Setup() {
	v, err := os.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(v, &Conf)
}
