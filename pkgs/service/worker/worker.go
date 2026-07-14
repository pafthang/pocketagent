package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/service/consumer"
)

// Worker is a base NATS JetStream worker lifecycle.
type Worker struct {
	Log           *slog.Logger
	Name          string
	NC            *nats.Conn
	JS            nats.JetStreamContext
	Consumer      *consumer.Consumer
	EventRecorder func(models.TaskEvent)
	healthServer  *echo.Echo
	healthAddr    string
}

// New connects to NATS and prepares a consumer.
func New(name, natsURL, logLevel string) (*Worker, error) {
	if logLevel != "" {
		_ = os.Setenv("LOG_LEVEL", logLevel)
	}

	nc, js, err := natsclient.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &Worker{
		Log:      common.NewSlogLogger(name),
		Name:     name,
		NC:       nc,
		JS:       js,
		Consumer: consumer.New(name, js),
	}, nil
}

// Subscribe delegates to the embedded consumer.
func (w *Worker) Subscribe(subject string, handler func(ctx context.Context, msg *nats.Msg)) (*nats.Subscription, error) {
	return w.Consumer.Subscribe(subject, handler)
}

// Publish sends a message to a JetStream subject.
func (w *Worker) Publish(subject string, data []byte) error {
	_, err := w.JS.Publish(subject, data)
	return err
}

// PublishMsg publishes a NATS message with optional correlation ID header.
func (w *Worker) PublishMsg(ctx context.Context, subject string, data []byte) error {
	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  make(nats.Header),
	}
	common.InjectContextHeaders(ctx, msg.Header)
	_, err := w.JS.PublishMsg(msg)
	return err
}

// PublishJSON marshals payload and publishes with correlation ID header.
func (w *Worker) PublishJSON(ctx context.Context, subject string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	return w.PublishMsg(ctx, subject, data)
}

// PublishEvent emits a task progress event on agents.events.{taskID}.
func (w *Worker) PublishEvent(ctx context.Context, taskID string, event models.TaskEvent) error {
	event.TaskID = taskID
	if w.EventRecorder != nil && event.SpaceID != "" {
		w.EventRecorder(event)
	}
	return w.PublishJSON(ctx, natsclient.EventSubject(taskID), event)
}

// Run blocks until shutdown signal and drains the NATS connection.
func (w *Worker) Run() error {
	defer w.NC.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	w.Log.Info("service started")

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	w.shutdownHealth(3 * time.Second)
	_ = common.ShutdownTelemetry(shutdownCtx)
	cancel()

	if err := w.NC.Drain(); err != nil {
		w.Log.Warn("nats drain failed", "error", err)
	}
	time.Sleep(10 * time.Second)

	w.Log.Info("service stopped")
	return nil
}