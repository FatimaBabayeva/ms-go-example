package properties

import "github.com/alexflint/go-arg"

// RootPath is project root path
const RootPath = "/v1/go-example"

type args struct {
	LogLevel string `arg:"env:LOG_LEVEL"`
	Port     int    `arg:"env:PORT"`
	DbUrl    string `arg:"env:DB_URL"`
	DbUser   string `arg:"env:DB_USER"`
	DbPass   string `arg:"env:DB_PASS"`
}

// DbConnStr constructs connection string from env variables
func (p *args) DbConnStr() string {
	return "postgres://" + p.DbUser + ":" + p.DbPass + "@" + p.DbUrl
}

// Props is for storing environment properties
var Props args

// LoadConfig loads service configuration into environment
func LoadConfig() {
	arg.Parse(&Props)
}
