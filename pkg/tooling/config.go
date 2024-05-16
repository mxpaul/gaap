package tooling

type Config struct {
	HTTPListenAddress string `yaml:"http_listen_address,omitempty"`
	MetricsPath       string `yaml:"metrics_path,omitempty"`
	LogRequests       bool   `yaml:"log_requests,omitempty"`
}
