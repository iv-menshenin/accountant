package configstorage

type (
	configSrc interface {
	}
	ConfigStorage struct {
		source []configSrc
	}
)

func InitCmdConfig() {

}

func InitEnvConfig() {

}

func InitCustomConfig(config map[string]string) {

}

func LoadConfig() {

}
