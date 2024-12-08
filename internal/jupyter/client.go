package jupyter

import (
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

	// Quotes
	Quotes(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quotes, error)
}

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

// Token implements Client interface
func (c *client) Token(_ context.Context, address string) (*Token, error) {
	var (
		req   = fasthttp.AcquireRequest()
		res   = fasthttp.AcquireResponse()
		token *Token
		err   error
	)

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(c.config.Jupyter.Token.APIHost + "/" + address)

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
	_, _ = uri.WriteString(c.config.Jupyter.Quote.APIHost)

	// write params
	_, _ = uri.WriteString("?inputMint=" + url.QueryEscape(input.Address))
	_, _ = uri.WriteString("&outputMint=" + url.QueryEscape(output.Address))
	_, _ = uri.WriteString("&amount=" + url.QueryEscape(strconv.FormatInt(input.Amount(amount), 10)))

	req.SetRequestURI(uri.String())

	// do request
	if err = c.client.Do(req, res); err != nil {
		return nil, err
	}

	// check http status
	if st := res.StatusCode(); st != http.StatusOK {
		return nil, fmt.Errorf("received unexpected http status code, %d", st)
	}

	// parse json
	if err = json.NewDecoder(req.BodyStream()).Decode(&quote); err != nil {
		return nil, err
	}

	logger.Debug("taken quotas", zap.Any("quote", quote))

	return quote, nil
}

// Quotes implements Client interface
func (c *client) Quotes(ctx context.Context, input, output *instruments.Instrument, amount int64) (*Quotes, error) {
	var (
		quotes = new(Quotes)
		err    error
	)

	// take direct quotes
	if quotes.Direct, err = c.Quote(ctx, input, output, amount); err != nil {
		return nil, err
	}

	// calculate output amount for api
	amount = output.Amount(quotes.Direct.OutAmount)

	// take reverse quotes
	if quotes.Reverse, err = c.Quote(ctx, output, input, amount); err != nil {
		return nil, err
	}

	return quotes, nil
}
