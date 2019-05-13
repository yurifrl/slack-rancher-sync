package rancher

type Config struct {
	Endpoint string
	User     string
	Token    string
}

func (c *Config) Reconcile() {
}
