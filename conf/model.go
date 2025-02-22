package conf

type MainConfig struct {
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
	} `mapstructure:"redis"`
	Cluster struct {
		IP      string `mapstructure:"ip"`
		Port    int    `mapstructure:"port"`
		SSLPort int    `mapstructure:"ssl_port"`
	}
}
