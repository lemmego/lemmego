package commands

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/lemmego/api/app"
	"github.com/lemmego/queue"
	"github.com/lemmego/tasker"
	"github.com/spf13/cobra"
)

type DummyJob struct {
	Message string `json:"message"`
	Delay   int    `json:"delay"`
}

func (j *DummyJob) Handle(ctx context.Context) error {
	slog.Info("processing dummy job", "message", j.Message)
	if j.Delay > 0 {
		time.Sleep(time.Duration(j.Delay) * time.Millisecond)
	}
	return nil
}

func (j *DummyJob) Tags() []string { return []string{"dummy"} }

type FailingDummyJob struct {
	Message string `json:"message"`
	Delay   int    `json:"delay"`
}

func (j *FailingDummyJob) Handle(ctx context.Context) error {
	return fmt.Errorf("simulated failure: %s", j.Message)
}

func (j *FailingDummyJob) Failed(ctx context.Context, payload []byte, err error) error {
	slog.Warn("dummy job failed handler called", "error", err)
	return nil
}

func (j *FailingDummyJob) Tags() []string { return []string{"dummy", "failing"} }

func init() {
	tasker.RegisterJob("*commands.DummyJob", func() tasker.Job { return &DummyJob{} })
	tasker.RegisterJob("*commands.FailingDummyJob", func() tasker.Job { return &FailingDummyJob{} })
}

func dispatchRun(cmd *cobra.Command, args []string) error {
	count, _ := cmd.Flags().GetInt("count")
	queueName, _ := cmd.Flags().GetString("queue")
	failRate, _ := cmd.Flags().GetInt("fail-rate")
	delayMs, _ := cmd.Flags().GetInt("delay")

	mgr := tasker.Global()
	if mgr == nil {
		return fmt.Errorf("tasker is not configured — add queue.QueueServiceProvider to bootstrap")
	}

	ctx := context.Background()
	for i := 0; i < count; i++ {
		msg := fmt.Sprintf("dummy job #%d (queue=%s)", i+1, queueName)
		var job tasker.Job
		if rand.Intn(100) < failRate {
			job = &FailingDummyJob{Message: msg, Delay: delayMs}
		} else {
			job = &DummyJob{Message: msg, Delay: delayMs}
		}
		_, err := queue.Dispatch(ctx, job, tasker.OnQueue(tasker.QueueName(queueName)))
		if err != nil {
			return fmt.Errorf("dispatch failed: %w", err)
		}
		fmt.Printf("dispatched job #%d\n", i+1)
	}
	fmt.Printf("done — %d jobs dispatched to %s\n", count, queueName)
	return nil
}

var DispatchCommand = func(a app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasker:dispatch",
		Short: "Dispatch dummy jobs to test the tasker dashboard",
		RunE:  dispatchRun,
	}
	cmd.Flags().IntP("count", "c", 5, "Number of jobs to dispatch")
	cmd.Flags().StringP("queue", "q", "default", "Queue to dispatch to")
	cmd.Flags().IntP("fail-rate", "f", 0, "Percentage chance of failure (0-100)")
	cmd.Flags().IntP("delay", "d", 0, "Simulated processing delay in ms")
	return cmd
}
