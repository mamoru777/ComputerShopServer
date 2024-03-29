package Config

type Config struct {
	Port int `env:"USER_GRPC_ADDR" envDefault:"13999"`

	PgPort   string `env:"PG_PORT" envDefault:"5432"`
	PgHost   string `env:"PG_HOST" envDefault:"0.0.0.0"`
	PgDBName string `env:"PG_DB_NAME" envDefault:"ComputerShop"`
	PgUser   string `env:"PG_USER" envDefault:"postgres"`
	PgPwd    string `env:"PG_PWD" envDefault:"159753"`

	SmtpPort           string `env:"SMTP_PORT" envDefault:"587"`
	SmtpAdr            string `env:"SMTP_ADR" envDefault:"smtp.gmail.com"`
	SmtpSenderEmail    string `env:"SMTP_EMAIL" envDefault:"mamoru90000@gmail.com"`
	SmtpSenderPassword string `env:"SMTP_PASSWORD" envDefault:"txivhhptkxmnravz"`
}
