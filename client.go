// Package elogik is a client for the e-Logik OMS e-commerce API v2.
package elogik

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// DefaultBaseURL is the EKAN Roissy instance of the e-Logik OMS API v2.
const DefaultBaseURL = "https://oms.ekan-roissy.fr/api/ecomm/v2"

// Client calls the API with HTTP Basic credentials. Use New to build one.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client

	user, password string

	mu        sync.Mutex
	callLimit string

	Articles             *ArticlesService
	Colis                *ColisService
	Commandes            *CommandesService
	CommandesFournisseur *CommandesFournisseurService
	Fournisseurs         *FournisseursService
	Packs                *PacksService
	Retours              *RetoursService
	Transporteurs        *TransporteursService
	Webhooks             *WebhooksService
}

// New returns a client for DefaultBaseURL using http.DefaultClient.
func New(user, password string) *Client {
	c := &Client{BaseURL: DefaultBaseURL, HTTPClient: http.DefaultClient, user: user, password: password}
	c.Articles = &ArticlesService{c}
	c.Colis = &ColisService{c}
	c.Commandes = &CommandesService{c}
	c.CommandesFournisseur = &CommandesFournisseurService{c}
	c.Fournisseurs = &FournisseursService{c}
	c.Packs = &PacksService{c}
	c.Retours = &RetoursService{c}
	c.Transporteurs = &TransporteursService{c}
	c.Webhooks = &WebhooksService{c}
	return c
}

// Error is returned for any non-2xx response.
//
// StatusCode 429 means the token bucket is empty: wait one second per
// token needed. A 400 with Message "Request limit exceeded" means a single
// call asked for more than the bucket size and must be split.
type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("elogik: HTTP %d %s: %s", e.StatusCode, e.Code, e.Message)
}

// CallLimit returns the token bucket state from the last
// X-OMS-Api-Call-Limit header seen. ok is false before any call.
func (c *Client) CallLimit() (remaining, size int, ok bool) {
	c.mu.Lock()
	v := c.callLimit
	c.mu.Unlock()
	_, err := fmt.Sscanf(v, "%d/%d", &remaining, &size)
	return remaining, size, err == nil
}

func (c *Client) do(ctx context.Context, method, path string, in, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.user, c.password)
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if v := resp.Header.Get("X-OMS-Api-Call-Limit"); v != "" {
		c.mu.Lock()
		c.callLimit = v
		c.mu.Unlock()
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		e := &Error{StatusCode: resp.StatusCode}
		if json.Unmarshal(b, e) != nil || (e.Code == "" && e.Message == "") {
			e.Message = strings.TrimSpace(string(b))
		}
		return e
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// call decodes the response into a new T, so each endpoint is one line.
func call[T any](ctx context.Context, c *Client, method, path string, in any) (*T, error) {
	var out T
	if err := c.do(ctx, method, path, in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
