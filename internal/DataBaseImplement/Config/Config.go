package Config

type Config struct {
	Port int `env:"USER_GRPC_ADDR" envDefault:"13999"`

	PgPort   string `env:"PG_PORT" envDefault:"5432"`
	PgHost   string `env:"PG_HOST" envDefault:"localhost"`
	PgDBName string `env:"PG_DB_NAME" envDefault:"postgres"`
	PgUser   string `env:"PG_USER" envDefault:"postgres"`
	PgPwd    string `env:"PG_PWD" envDefault:"159753"`
}
