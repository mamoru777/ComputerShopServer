package Config

type Config struct {
	Port int `env:"USER_GRPC_ADDR" envDefault:"13999"`

	PgPort   string `env:"PG_PORT" envDefault:"5432"`
	PgHost   string `env:"PG_HOST" envDefault:"5.3.65.108"`
	PgDBName string `env:"PG_DB_NAME" envDefault:"ComputerShop"`
	PgUser   string `env:"PG_USER" envDefault:"postgres"`
	PgPwd    string `env:"PG_PWD" envDefault:"159753"`
}
