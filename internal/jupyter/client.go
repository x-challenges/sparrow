package jupyter

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/x-challenges/raven/common/json"

	"sparrow/internal/instruments"
)

// Client
type Client interface {
	// Token
	Token(ctx context.Context, address string) (*Token, error)

	// Quote
	Quote(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quote, error)
}

// Client interface implementation
type client struct {
	logger *zap.Logger
	client *fasthttp.Client
	config *Config

	quoteHosts *Balancer
}

var _ Client = (*client)(nil)

// NewClient
func newClient(logger *zap.Logger, fclient *fasthttp.Client, config *Config) (*client, error) {
	return &client{
		logger:     logger,
		client:     fclient,
		config:     config,
		quoteHosts: NewBalancer(config.Jupyter.Quote.Hosts...),
	}, nil
}

// Token implements Client interface
func (c *client) Token(_ context.Context, address string) (*Token, error) {
	var (
		token *Token
		err   error
	)

	var req = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	var res = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(c.config.Jupyter.Token.Host + "/" + address)
	req.Header.Set("Connection", "keep-alive")

	// do request
	if err = c.client.Do(req, res); err != nil {
		return nil, err
	}

	// parse json
	if err = json.NewDecoder(req.BodyStream()).Decode(&token); err != nil {
		return nil, err
	}

	return token, nil
}

// Quote implements Client interface
func (c *client) Quote(_ context.Context, input, output *instruments.Instrument, amount int64) (*Quote, error) {
	var (
		quote *Quote
		err   error
	)

	// prepare logger
	var logger = c.logger.With(
		zap.Any("input", input),
		zap.Any("output", input),
		zap.Int64("amount", amount),
	)

	// acquire fasthttp request
	var req = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// acquire fasthttp response
	var res = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	// prepare strings builder for url
	var uri strings.Builder
	uri.Grow(256)

	// write host
	_, _ = uri.WriteString(c.quoteHosts.Next())

	// write params
	_, _ = uri.WriteString("?inputMint=" + url.QueryEscape(input.Address))
	_, _ = uri.WriteString("&outputMint=" + url.QueryEscape(output.Address))
	_, _ = uri.WriteString("&amount=" + url.QueryEscape(strconv.FormatInt(amount, 10)))

	req.SetRequestURI(uri.String())
	req.Header.SetMethod(fasthttp.MethodGet)

	req.Header.Set("Connection", "keep-alive")

	// do request
	if err = c.client.Do(req, res); err != nil {
		return nil, err
	}

	// check http status
	if st := res.StatusCode(); st != http.StatusOK {
		return nil, fmt.Errorf("received unexpected http status code, %d", st)
	}

	// parse json
	if err = json.NewDecoder(bytes.NewBuffer(res.Body())).Decode(&quote); err != nil {
		return nil, err
	}

	logger.Debug("taken quotas", zap.Any("quote", quote))

	return quote, nil
}
