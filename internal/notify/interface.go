package notify

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
)

// Notifier defines the interface for notification systems
type Notifier interface {
	// Initialize sets up the notifier with necessary configurations
	Initialize(ctx context.Context, config map[string]string) error
	// NotifyJobStatus sends a notification about job status
	NotifyJobStatus(ctx context.Context, job *batchv1.Job) error
}

// NotifierFactory creates a new instance of a notifier
type NotifierFactory interface {
	// Create returns a new Notifier instance
	Create(notifierType string) (Notifier, error)
}