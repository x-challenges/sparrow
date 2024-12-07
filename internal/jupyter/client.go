package jupyter

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// Client
type Client interface{}

// Client interface implementation
type client struct {
	logger *zap.Logger
	client *fasthttp.Client
	config *Config
}

var _ Client = (*client)(nil)

// NewClient
func newClient(logger *zap.Logger, fclient *fasthttp.Client, config *Config) (*client, error) {
	return &client{
		logger: logger,
		client: fclient,
		config: config,
	}, nil
}
