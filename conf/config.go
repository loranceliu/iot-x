package config

type Config struct {
	Server Server `yaml:"Server"`
	Redis  Redis  `yaml:"Redis"`
}

type Server struct {
	Node      string `yaml:"Node"`
	IP        string `yaml:"IP"`
	Port      int    `yaml:"Port"`
	MachineId string `yaml:"MachineId"`
}

type Redis struct {
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
	Pwd  string `yaml:"Pwd"`
	DB   int    `yaml:"DB"`
}
