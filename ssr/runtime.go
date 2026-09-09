package ssr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"andurel-site/config"
)

// Runtime starts Node from INERTIA_SSR_LISTEN. It never reads INERTIA_SSR_URL.
type Runtime struct {
	cfg        config.Inertia
	healthURL  string
	healthHTTP *http.Client
	errors     chan error

	mu            sync.Mutex
	command       *exec.Cmd
	done          chan error
	stopping      bool
	bundleCleanup func()
}

func NewRuntime(cfg config.Inertia) (*Runtime, error) {
	healthURL := cfg.SSRHealthURL()
	if healthURL == "" {
		return nil, fmt.Errorf("inertia SSR listen address is invalid")
	}

	return &Runtime{
		cfg:       cfg,
		healthURL: healthURL,
		healthHTTP: &http.Client{
			Timeout: cfg.SSRRequestTimeout,
		},
		errors: make(chan error, 1),
	}, nil
}

func (runtime *Runtime) Errors() <-chan error {
	if runtime == nil {
		return nil
	}

	return runtime.errors
}

func (runtime *Runtime) Start(ctx context.Context) (err error) {
	if runtime == nil {
		return nil
	}

	bundlePath, cleanup, err := resolveSSRBundle(runtime.cfg.SSRBundle)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil && cleanup != nil {
			runtime.mu.Lock()
			runtime.bundleCleanup = nil
			runtime.mu.Unlock()
			cleanup()
		}
	}()

	executable, err := exec.LookPath(runtime.cfg.SSRRuntime)
	if err != nil {
		return fmt.Errorf("inertia SSR runtime %q: %w", runtime.cfg.SSRRuntime, err)
	}

	output, err := exec.CommandContext(ctx, executable, "--version").Output()
	if err != nil {
		return fmt.Errorf("inspect SSR runtime version: %w", err)
	}

	var major int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(output)), "v%d.", &major); err != nil {
		return fmt.Errorf(
			"inspect SSR runtime version %q: %w",
			strings.TrimSpace(string(output)),
			err,
		)
	}

	if major < runtime.cfg.SSRMinimumMajor {
		return fmt.Errorf(
			"inertia SSR requires runtime major %d or newer (found %q)",
			runtime.cfg.SSRMinimumMajor,
			strings.TrimSpace(string(output)),
		)
	}

	runtime.mu.Lock()
	if runtime.command != nil {
		runtime.mu.Unlock()
		return fmt.Errorf("inertia SSR runtime is already started")
	}

	command := exec.Command(executable, bundlePath)
	command.Env = append(os.Environ(),
		"INERTIA_SSR_HOST="+runtime.cfg.SSRBindHost(),
		"INERTIA_SSR_PORT="+runtime.cfg.SSRBindPort(),
	)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	done := make(chan error, 1)
	if err := command.Start(); err != nil {
		runtime.mu.Unlock()
		return fmt.Errorf("start inertia SSR runtime: %w", err)
	}

	runtime.command = command
	runtime.done = done
	runtime.stopping = false
	runtime.bundleCleanup = cleanup
	runtime.mu.Unlock()

	go func() {
		waitErr := command.Wait()
		done <- waitErr
		close(done)

		runtime.mu.Lock()
		expected := runtime.stopping
		if runtime.command == command {
			runtime.command = nil
			runtime.done = nil
			runtime.stopping = false
		}
		runtime.mu.Unlock()

		if expected {
			return
		}

		err := fmt.Errorf("inertia SSR runtime stopped unexpectedly: %w", normalizeWaitError(waitErr))
		select {
		case runtime.errors <- err:
		default:
		}
	}()

	startupCtx, cancel := context.WithTimeout(ctx, runtime.cfg.SSRStartupTimeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		if err := runtime.health(startupCtx); err == nil {
			cleanup = nil
			return nil
		}

		select {
		case waitErr := <-done:
			return fmt.Errorf(
				"inertia SSR runtime exited during startup: %w",
				normalizeWaitError(waitErr),
			)
		case <-ticker.C:
		case <-startupCtx.Done():
			_ = command.Process.Kill()
			return fmt.Errorf("inertia SSR runtime health check: %w", startupCtx.Err())
		}
	}
}

func (runtime *Runtime) Stop(ctx context.Context) error {
	if runtime == nil {
		return nil
	}
	defer runtime.removeExtractedBundle()

	runtime.mu.Lock()
	command := runtime.command
	done := runtime.done
	if command == nil {
		runtime.mu.Unlock()
		return nil
	}

	runtime.stopping = true
	runtime.mu.Unlock()

	shutdownErr := runtime.shutdown(ctx)
	select {
	case waitErr := <-done:
		if waitErr != nil && shutdownErr == nil {
			shutdownErr = waitErr
		}

		return shutdownErr
	case <-ctx.Done():
		killErr := command.Process.Kill()
		return errors.Join(shutdownErr, ctx.Err(), killErr)
	}
}

func (runtime *Runtime) removeExtractedBundle() {
	runtime.mu.Lock()
	cleanup := runtime.bundleCleanup
	runtime.bundleCleanup = nil
	runtime.mu.Unlock()
	if cleanup != nil {
		cleanup()
	}
}

func (runtime *Runtime) health(ctx context.Context) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, runtime.healthURL+"/health", nil)
	if err != nil {
		return err
	}

	response, err := runtime.healthHTTP.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, 64<<10))
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("health status %d", response.StatusCode)
	}

	var health struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &health); err != nil {
		return err
	}
	if !strings.EqualFold(health.Status, "ok") {
		return fmt.Errorf("renderer is not healthy")
	}

	return nil
}

func (runtime *Runtime) shutdown(ctx context.Context) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		runtime.healthURL+"/shutdown",
		nil,
	)
	if err != nil {
		return err
	}

	response, err := runtime.healthHTTP.Do(request)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 64<<10))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("shutdown status %d", response.StatusCode)
	}

	return nil
}

func normalizeWaitError(err error) error {
	if err == nil {
		return errors.New("process exited")
	}

	return err
}
