package trollr

type AppConfig struct {
	TrollBin string `config:"troll-bin,required"`
	Mosmllib string `config:"mosmllib,required"`
	Listen string `config:"listen"`
}

