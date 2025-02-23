package notify_test

import (
	"context"
	"testing"

	"github.com/takutakahashi/job-notify/internal/notify"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMockNotifier(t *testing.T) {
	ctx := context.Background()

	t.Run("successful initialization", func(t *testing.T) {
		notifier := notify.NewMockNotifier()
		config := map[string]string{
			"webhookURL": "https://hooks.slack.com/services/xxx",
			"channel":    "#test",
		}

		err := notifier.Initialize(ctx, config)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		for k, v := range config {
			if notifier.InitializedWithConfig[k] != v {
				t.Errorf("expected config[%s] = %s, got %s", k, v, notifier.InitializedWithConfig[k])
			}
		}
	})

	t.Run("failed initialization", func(t *testing.T) {
		notifier := notify.NewMockNotifier()
		notifier.ShouldFail = true

		err := notifier.Initialize(ctx, map[string]string{})
		if err != notify.ErrMockFailure {
			t.Errorf("expected ErrMockFailure, got %v", err)
		}
	})

	t.Run("successful notification", func(t *testing.T) {
		notifier := notify.NewMockNotifier()
		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "default",
			},
		}

		err := notifier.NotifyJobStatus(ctx, job)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(notifier.NotifiedJobs) != 1 {
			t.Errorf("expected 1 notified job, got %d", len(notifier.NotifiedJobs))
		}

		if notifier.NotifiedJobs[0].Name != job.Name {
			t.Errorf("expected job name %s, got %s", job.Name, notifier.NotifiedJobs[0].Name)
		}
	})

	t.Run("failed notification", func(t *testing.T) {
		notifier := notify.NewMockNotifier()
		notifier.ShouldFail = true

		err := notifier.NotifyJobStatus(ctx, &batchv1.Job{})
		if err != notify.ErrMockFailure {
			t.Errorf("expected ErrMockFailure, got %v", err)
		}
	})
}

func TestMockNotifierFactory(t *testing.T) {
	t.Run("create notifier", func(t *testing.T) {
		factory := notify.NewMockNotifierFactory()
		notifier, err := factory.Create("slack")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if notifier == nil {
			t.Error("expected non-nil notifier")
		}
	})
}