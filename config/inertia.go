package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/validation"
)

const (
	DefaultInertiaContainerID         = "app"
	DefaultInertiaViteDevURL          = "http://localhost:5173/assets/dist"
	DefaultInertiaProtocolDebug       = false
	DefaultInertiaSSRRuntime          = "node"
	DefaultInertiaSSRBundle           = "assets/dist/ssr/ssr.js"
	DefaultInertiaSSRListen           = "http://127.0.0.1:13714"
	DefaultInertiaSSRURL              = "http://127.0.0.1:13714"
	DefaultInertiaSSRStartupTimeout   = 10 * time.Second
	DefaultInertiaSSRRequestTimeout   = 2 * time.Second
	DefaultInertiaSSRMaxResponseBytes = int64(2097152)
	DefaultInertiaSSRMinimumMajor     = 22
	DefaultInertiaSSRFailFast         = false
	DefaultInertiaEntryPoint          = "resources/js/app.tsx"
)

type Inertia struct {
	ContainerID         string
	ViteDevURL          string
	ProtocolDebug       bool
	SSRRuntime          string
	SSRBundle           string
	SSRListen           string
	SSRURL              string
	SSRStartupTimeout   time.Duration
	SSRRequestTimeout   time.Duration
	SSRMaxResponseBytes int64
	SSRMinimumMajor     int
	SSRFailFast         bool
	EntryPoint          string
}

func NewInertia() (Inertia, error) {
	env := newEnvironment()
	cfg := Inertia{
		ContainerID:   env.String("INERTIA_CONTAINER_ID", DefaultInertiaContainerID),
		ViteDevURL:    env.String("INERTIA_VITE_DEV_URL", DefaultInertiaViteDevURL),
		ProtocolDebug: env.Bool("INERTIA_PROTOCOL_DEBUG", DefaultInertiaProtocolDebug),
		SSRRuntime:    env.String("INERTIA_SSR_RUNTIME", DefaultInertiaSSRRuntime),
		SSRBundle:     env.String("INERTIA_SSR_BUNDLE", DefaultInertiaSSRBundle),
		SSRListen:     env.String("INERTIA_SSR_LISTEN", DefaultInertiaSSRListen),
		SSRURL:        env.String("INERTIA_SSR_URL", DefaultInertiaSSRURL),
		SSRStartupTimeout: env.Duration(
			"INERTIA_SSR_STARTUP_TIMEOUT",
			DefaultInertiaSSRStartupTimeout,
		),
		SSRRequestTimeout: env.Duration(
			"INERTIA_SSR_REQUEST_TIMEOUT",
			DefaultInertiaSSRRequestTimeout,
		),
		SSRMaxResponseBytes: env.Int64(
			"INERTIA_SSR_MAX_RESPONSE_BYTES",
			DefaultInertiaSSRMaxResponseBytes,
		),
		SSRMinimumMajor: env.Int("INERTIA_SSR_MINIMUM_MAJOR", DefaultInertiaSSRMinimumMajor),
		SSRFailFast:     env.Bool("INERTIA_SSR_FAIL_FAST", DefaultInertiaSSRFailFast),
		EntryPoint:      env.String("INERTIA_ENTRY_POINT", DefaultInertiaEntryPoint),
	}

	if err := errors.Join(env.Err(), cfg.validate()); err != nil {
		return Inertia{}, fmt.Errorf("config: inertia: %w", err)
	}

	return cfg, nil
}

func (c Inertia) validate() error {
	b := validation.NewBuilder()
	b.Required("ContainerID", c.ContainerID)
	b.Required("ViteDevURL", c.ViteDevURL)
	b.Required("EntryPoint", c.EntryPoint)
	b.Required("SSRRuntime", c.SSRRuntime)
	b.Required("SSRBundle", c.SSRBundle)
	b.Required("SSRListen", c.SSRListen)
	if err := b.Err(); err != nil {
		return err
	}

	if _, _, err := parseSSRListen(c.SSRListen); err != nil {
		return err
	}

	http := inertia.SSRClientConfig{
		URL:              c.SSRURL,
		Timeout:          c.SSRRequestTimeout,
		MaxResponseBytes: c.SSRMaxResponseBytes,
	}
	if err := http.Validate(); err != nil {
		return err
	}
	if c.SSRStartupTimeout <= 0 {
		return fmt.Errorf("SSRStartupTimeout must be positive")
	}
	if c.SSRMinimumMajor <= 0 {
		return fmt.Errorf("SSRMinimumMajor must be positive")
	}

	return nil
}

// SSRClientConfig is where cmd/app POSTs /render.
func (c Inertia) SSRClientConfig() inertia.SSRClientConfig {
	return inertia.SSRClientConfig{
		URL:              c.SSRURL,
		Timeout:          c.SSRRequestTimeout,
		MaxResponseBytes: c.SSRMaxResponseBytes,
	}
}

// SSRBindHost is the address Node listens on (INERTIA_SSR_HOST).
func (c Inertia) SSRBindHost() string {
	host, _, err := parseSSRListen(c.SSRListen)
	if err != nil {
		return ""
	}

	return host
}

// SSRBindPort is the port Node listens on (INERTIA_SSR_PORT).
func (c Inertia) SSRBindPort() string {
	_, port, err := parseSSRListen(c.SSRListen)
	if err != nil {
		return ""
	}

	return port
}

// SSRHealthURL is loopback on the listen port so cmd/ssr can probe the process
// it started, even when Node binds 0.0.0.0.
func (c Inertia) SSRHealthURL() string {
	host, port, err := parseSSRListen(c.SSRListen)
	if err != nil {
		return ""
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsUnspecified() {
		host = "127.0.0.1"
	}

	return "http://" + net.JoinHostPort(host, port)
}

func parseSSRListen(raw string) (string, string, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" || parsed.Scheme != "http" {
		return "", "", fmt.Errorf("SSRListen must be an HTTP URL")
	}

	host := parsed.Hostname()
	port := parsed.Port()
	if port == "" {
		return "", "", fmt.Errorf("SSRListen must include a port")
	}
	if !validSSRListenHost(host) {
		return "", "", fmt.Errorf(
			"SSRListen must bind an IP address or localhost, not a service hostname",
		)
	}

	return host, port, nil
}

func validSSRListenHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}

	return net.ParseIP(host) != nil
}
