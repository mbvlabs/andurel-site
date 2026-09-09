package ssr

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"andurel-site/config"

	"github.com/mbvlabs/andurel/pkg/inertia"
)

// Runtime starts Node using INERTIA_SSR_LISTEN (bind) and health-checks loopback.
// cmd/app talks to the renderer through INERTIA_SSR_URL, which is independent.
type Runtime struct {
	cfg      config.Inertia
	renderer *inertia.HTTPRenderer
	errors   chan error

	mu       sync.Mutex
	command  *exec.Cmd
	done     chan error
	stopping bool
}

func NewRuntime(cfg config.Inertia) (*Runtime, error) {
	renderer, err := inertia.NewHTTPRenderer(cfg.SSRHealthConfig())
	if err != nil {
		return nil, err
	}

	return &Runtime{
		cfg:      cfg,
		renderer: renderer,
		errors:   make(chan error, 1),
	}, nil
}

func (runtime *Runtime) Errors() <-chan error {
	if runtime == nil {
		return nil
	}

	return runtime.errors
}

func (runtime *Runtime) Start(ctx context.Context) error {
	if runtime == nil {
		return nil
	}

	if _, err := os.Stat(runtime.cfg.SSRBundle); err != nil {
		return fmt.Errorf("inertia SSR bundle %q: %w", runtime.cfg.SSRBundle, err)
	}

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

	command := exec.Command(executable, runtime.cfg.SSRBundle)
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
		if err := runtime.renderer.Health(startupCtx); err == nil {
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

	runtime.mu.Lock()
	command := runtime.command
	done := runtime.done
	if command == nil {
		runtime.mu.Unlock()
		return nil
	}

	runtime.stopping = true
	runtime.mu.Unlock()

	shutdownErr := runtime.renderer.Shutdown(ctx)
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

func normalizeWaitError(err error) error {
	if err == nil {
		return errors.New("process exited")
	}

	return err
}
