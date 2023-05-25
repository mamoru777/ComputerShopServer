package ConfigServ

type Config struct {
	GRPCAddr string `env:"USER_GRPC_ADDR" envDefault:":1399"`
}
