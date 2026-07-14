package supervisor

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pafthang/pocketagent/internal/ctrl/catalog"
	"github.com/pafthang/pocketagent/internal/ctrl/port"
	"github.com/pafthang/pocketagent/pkgs/common"
)

const (
	defaultStartTimeout = 30 * time.Second
	defaultStopTimeout  = 10 * time.Second
)

// Supervisor manages child service processes for local development.
type Supervisor struct {
	root         string
	configDir    string
	startTimeout time.Duration
	stopTimeout  time.Duration
	services     []catalog.Service
	procs        map[string]*managedProc
	order        []string
}

// New creates a supervisor bound to the project root and configs directory.
func New(root, configDir string, ctrlCfg *common.CtrlConfig) *Supervisor {
	startTimeout := defaultStartTimeout
	stopTimeout := defaultStopTimeout
	if ctrlCfg != nil {
		if ctrlCfg.StartTimeoutSec > 0 {
			startTimeout = time.Duration(ctrlCfg.StartTimeoutSec) * time.Second
		}
		if ctrlCfg.StopTimeoutSec > 0 {
			stopTimeout = time.Duration(ctrlCfg.StopTimeoutSec) * time.Second
		}
	}

	return &Supervisor{
		root:         root,
		configDir:    configDir,
		startTimeout: startTimeout,
		stopTimeout:  stopTimeout,
		procs:        make(map[string]*managedProc),
	}
}

// Register adds services to the startup graph.
func (s *Supervisor) Register(services ...catalog.Service) {
	s.services = append(s.services, services...)
}

// Run starts services in dependency order and blocks until shutdown.
func (s *Supervisor) Run() error {
	order, err := topoSort(s.services)
	if err != nil {
		return err
	}
	s.order = order

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("ctrl: starting %d services from %s\n", len(s.order), s.root)

	if err := port.FreeAll(port.CollectManaged(s.services)); err != nil {
		fmt.Printf("ctrl: warning: stale port cleanup: %v\n", err)
	}

	for _, name := range s.order {
		def := s.defByName(name)
		if def == nil {
			return fmt.Errorf("unknown service %q", name)
		}

		for _, dep := range def.DependsOn {
			if _, ok := s.procs[dep]; !ok {
				return fmt.Errorf("service %q depends on %q which is not running", name, dep)
			}
		}

		for _, p := range port.ServicePorts(def) {
			if err := port.EnsureFree(p); err != nil {
				s.shutdownStarted()
				return fmt.Errorf("port %d for %s: %w", p, name, err)
			}
		}

		proc, err := s.start(def)
		if err != nil {
			s.shutdownStarted()
			return fmt.Errorf("start %s: %w", name, err)
		}
		s.procs[name] = proc

		if def.WaitPort > 0 {
			addr := fmt.Sprintf("127.0.0.1:%d", def.WaitPort)
			if err := waitForPort(ctx, addr, s.startTimeout, proc); err != nil {
				s.shutdownStarted()
				return fmt.Errorf("wait %s (%s): %w", name, addr, err)
			}
			fmt.Printf("ctrl: %s ready on %s\n", name, addr)
		} else if def.HealthPort > 0 {
			addr := fmt.Sprintf("http://127.0.0.1:%d/health", def.HealthPort)
			if err := waitForHealth(ctx, addr, s.startTimeout, proc); err != nil {
				s.shutdownStarted()
				return fmt.Errorf("wait %s (%s): %w", name, addr, err)
			}
			fmt.Printf("ctrl: %s ready on %s\n", name, addr)
		} else {
			fmt.Printf("ctrl: %s started (pid %d)\n", name, proc.cmd.Process.Pid)
		}
	}

	fmt.Println("ctrl: all services running — Ctrl+C to stop")

	done := make(chan error, 1)
	var waitWG sync.WaitGroup
	for name, proc := range s.procs {
		waitWG.Add(1)
		go func(n string, p *managedProc) {
			defer waitWG.Done()
			if err := p.cmd.Wait(); err != nil {
				select {
				case done <- fmt.Errorf("%s exited: %w", n, err):
				default:
				}
			}
		}(name, proc)
	}
	go func() {
		waitWG.Wait()
		select {
		case done <- nil:
		default:
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("\nctrl: shutting down...")
		s.shutdownStarted()
		return nil
	case err := <-done:
		s.shutdownStarted()
		return err
	}
}

func (s *Supervisor) defByName(name string) *catalog.Service {
	for i := range s.services {
		if s.services[i].Name == name {
			return &s.services[i]
		}
	}
	return nil
}