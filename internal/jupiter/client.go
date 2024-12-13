package jupiter

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/x-challenges/raven/common/json"

	"sparrow/internal/jupiter/balancer"
)

// Client
type Client interface {
	// Tokens
	Tokens(ctx context.Context) (Tokens, error)

	// Prices
	Prices(ctx context.Context, input string, outputs ...string) (*Prices, error)

	// Quote
	Quote(
		ctx context.Context,
		input, output string,
		amount int64,
		options ...QuoteOptionFunc,
	) (*Quote, error)
}

// Client interface implementation
type client struct {
	logger *zap.Logger
	client *fasthttp.Client
	config *Config

	tokenHosts balancer.Balancer
	priceHosts balancer.Balancer
	quoteHosts balancer.Balancer
}

var _ Client = (*client)(nil)

// NewClient
func newClient(logger *zap.Logger, fclient *fasthttp.Client, config *Config) (*client, error) {
	return &client{
		logger:     logger,
		client:     fclient,
		config:     config,
		tokenHosts: balancer.NewBalancer(config.Jupiter.Token.Hosts...),
		priceHosts: balancer.NewBalancer(config.Jupiter.Price.Hosts...),
		quoteHosts: balancer.NewBalancer(config.Jupiter.Quote.Hosts...),
	}, nil
}

// Tokens implements Client interface
func (c *client) Tokens(_ context.Context) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	var req = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	var res = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	// prepare strings builder for url
	var uri strings.Builder
	uri.Grow(256)

	// write host
	_, _ = uri.WriteString(c.tokenHosts.Next())

	if tags := c.config.Jupiter.Token.Tags; len(tags) > 0 {
		_, _ = uri.WriteString("?tags=" + strings.Join(tags, ","))
	}

	req.SetRequestURI(uri.String())
	req.Header.SetMethod(fasthttp.MethodGet)

	req.Header.Set("Connection", "keep-alive")

	// do request
	if err = c.client.Do(req, res); err != nil {
		return nil, err
	}

	// check http status
	if st := res.StatusCode(); st != http.StatusOK {
		return nil, ErrUnexpectedStatusCode
	}

	// parse json
	if err = json.NewDecoder(bytes.NewBuffer(res.Body())).Decode(&tokens); err != nil {
		return nil, err
	}

	c.logger.Debug("taken tokens", zap.Int("count", len(tokens)))

	return tokens, nil
}

// Prices implement Client interface
func (c *client) Prices(_ context.Context, input string, outputs ...string) (*Prices, error) {
	var (
		prices *Prices
		err    error
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
	_, _ = uri.WriteString(c.priceHosts.Next())

	// write params
	_, _ = uri.WriteString("?vsToken=" + input)
	_, _ = uri.WriteString("&ids=" + strings.Join(outputs, ","))

	req.SetRequestURI(uri.String())
	req.Header.SetMethod(fasthttp.MethodGet)

	req.Header.Set("Connection", "keep-alive")

	// do request
	if err = c.client.Do(req, res); err != nil {
		return nil, err
	}

	// check http status
	if st := res.StatusCode(); st != http.StatusOK {
		return nil, ErrUnexpectedStatusCode
	}

	// parse json
	if err = json.NewDecoder(bytes.NewBuffer(res.Body())).Decode(&prices); err != nil {
		return nil, err
	}

	c.logger.Debug("taken prices", zap.Any("quote", prices))

	return prices, nil
}

// Quote implements Client interface
func (c *client) Quote(
	_ context.Context,
	input, output string,
	amount int64,
	opts ...QuoteOptionFunc,
) (*Quote, error) {
	var (
		options = NewQuoteOptions().Apply(opts...)
		quote   *Quote
		err     error
	)

	// prepare logger
	var logger = c.logger.With(
		zap.Any("input", input),
		zap.Any("output", input),
		zap.Int64("amount", amount),
		zap.Any("options", options),
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
	_, _ = uri.WriteString("?inputMint=" + url.QueryEscape(input))
	_, _ = uri.WriteString("&outputMint=" + url.QueryEscape(output))
	_, _ = uri.WriteString("&amount=" + url.QueryEscape(strconv.FormatInt(amount, 10)))

	// use onlyDirectRoutes
	if c.config.Jupiter.Quote.OnlyDirectRoutes {
		_, _ = uri.WriteString("&onlyDirectRoutes=true")
	}

	// use restrictIntermediateTokens
	if c.config.Jupiter.Quote.RestrictIntermediateTokens {
		_, _ = uri.WriteString("&restrictIntermediateTokens=true")
	}

	// swap mode
	switch options.SwapMode {
	case ExactIn:
		_, _ = uri.WriteString("&swapMode=ExactIn")
	case ExactOut:
		_, _ = uri.WriteString("&swapMode=ExactOut")
	}

	req.SetRequestURI(uri.String())
	req.Header.SetMethod(fasthttp.MethodGet)

	req.Header.Set("Connection", "keep-alive")

	// do request
	if err = c.client.Do(req, res); err != nil {
		return nil, err
	}

	// check http status
	switch st := res.StatusCode(); st {

	case http.StatusOK:
		break

	// cant found any route for quote
	case http.StatusBadRequest:
		return nil, ErrRouteNotFound

	// unexpected status code
	default:
		return nil, ErrUnexpectedStatusCode
	}

	// parse json
	if err = json.NewDecoder(bytes.NewBuffer(res.Body())).Decode(&quote); err != nil {
		return nil, err
	}

	logger.Debug("taken quote", zap.Any("quote", quote))

	return quote, nil
}
