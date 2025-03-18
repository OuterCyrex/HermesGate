package conf

type MainConfig struct {
	DashBoard   DashBoardConfig   `mapstructure:"dashboard"`
	ProxyServer ProxyServerConfig `mapstructure:"proxy_server"`
}

type DashBoardConfig struct {
	Host  string `mapstructure:"host"`
	Port  int    `mapstructure:"port"`
	Mysql struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"mysql"`
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		Database int    `mapstructure:"database"`
	} `mapstructure:"redis"`
}

type ProxyServerConfig struct {
	Host      string `mapstructure:"host"`
	HttpPort  int    `mapstructure:"http_port"`
	HttpsPort int    `mapstructure:"https_port"`
}
