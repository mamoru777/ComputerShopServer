package Config

type Config struct {
	Port int `env:"SERVER_PORT" envDefault:"13005"`

	PgPort   string `env:"PG_PORT" envDefault:"5432"`
	PgHost   string `env:"PG_HOST" envDefault:"localhost"`
	PgDBName string `env:"PG_DB_NAME" envDefault:"postgres"`
	PgUser   string `env:"PG_USER" envDefault:"mamoru1"`
	PgPwd    string `env:"PG_PWD" envDefault:""`
}
