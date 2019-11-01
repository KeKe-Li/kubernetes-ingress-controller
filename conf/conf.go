package conf

// Config 项目配置类
type Config struct {
	Name    string `toml:"name"`
	Env     string `toml:"env"`
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	TLSPort int    `toml:"tls_port"`
	Log     struct {
		Level  string `toml:"level"`  // debug|info|warning|error|fatal
		Format string `toml:"format"` // text|json
	} `toml:"log"`

	PProf struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	} `toml:"pprof"`
}

var requires = func() []string {
	var requires = []string{
		"name",
		"env",
		"host",
		"port",
		"tls_port",

		"log.level",
		"log.format",

		"pprof.port",
	}
	return requires
}()
