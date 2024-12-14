package jupiter

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/x-challenges/raven/common/json"
	"github.com/x-challenges/raven/modules/fasthttp"

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

// Request
func (c *client) getRequest(ctx context.Context, host, getParams string, resp any) error {
	var (
		started    = time.Now()
		statusCode int
		err        error
	)

	defer func() {
		requestLatency.Record(ctx, time.Since(started).Milliseconds(), metric.WithAttributes(
			attribute.String("host", host),
			attribute.Int("status_code", statusCode),
		))
	}()

	var req = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	var res = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(host + getParams)
	req.Header.Set("Connection", "keep-alive")

	// do request
	if err = c.client.Do(req, res); err != nil {
		return err
	}

	// remap errors
	switch statusCode = res.StatusCode(); statusCode {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		err = ErrNotFoundStatusCode
	default:
		err = ErrUnexpectedStatusCode
	}

	if err != nil {
		return err
	}

	// parse json
	if err = json.NewDecoder(bytes.NewBuffer(res.Body())).Decode(resp); err != nil {
		return err
	}

	return nil
}

// Tokens implements Client interface
func (c *client) Tokens(ctx context.Context) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	// prepare strings builder for url
	var uri strings.Builder
	uri.Grow(128)

	if tags := c.config.Jupiter.Token.Tags; len(tags) > 0 {
		_, _ = uri.WriteString("?tags=" + strings.Join(tags, ","))
	}

	// make request
	if err = c.getRequest(ctx, c.tokenHosts.Next(), uri.String(), &tokens); err != nil {
		return nil, err
	}

	c.logger.Debug("taken tokens", zap.Int("count", len(tokens)))

	return tokens, nil
}

// Prices implement Client interface
func (c *client) Prices(ctx context.Context, input string, outputs ...string) (*Prices, error) {
	var (
		prices *Prices
		err    error
	)

	// prepare strings builder for url
	var uri strings.Builder
	uri.Grow(128)

	// write params
	_, _ = uri.WriteString("?vsToken=" + input)
	_, _ = uri.WriteString("&ids=" + strings.Join(outputs, ","))

	// make request
	if err = c.getRequest(ctx, c.priceHosts.Next(), uri.String(), &prices); err != nil {
		return nil, err
	}

	c.logger.Debug("taken prices", zap.Any("quote", prices))

	return prices, nil
}

// Quote implements Client interface
func (c *client) Quote(
	ctx context.Context,
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

	// prepare strings builder for url
	var uri strings.Builder
	uri.Grow(256)

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

	// make request
	if err = c.getRequest(ctx, c.quoteHosts.Next(), uri.String(), &quote); err != nil {
		if errors.Is(err, ErrNotFoundStatusCode) {
			return nil, ErrRouteNotFound
		}
		return nil, err
	}

	logger.Debug("taken quote", zap.Any("quote", quote))

	return quote, nil
}
