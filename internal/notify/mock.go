package notify

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
)

// MockNotifier implements the Notifier interface for testing
type MockNotifier struct {
	InitializedWithConfig map[string]string
	NotifiedJobs         []*batchv1.Job
	ShouldFail          bool
}

func NewMockNotifier() *MockNotifier {
	return &MockNotifier{
		InitializedWithConfig: make(map[string]string),
		NotifiedJobs:         make([]*batchv1.Job, 0),
	}
}

func (m *MockNotifier) Initialize(ctx context.Context, config map[string]string) error {
	if m.ShouldFail {
		return ErrMockFailure
	}
	m.InitializedWithConfig = config
	return nil
}

func (m *MockNotifier) NotifyJobStatus(ctx context.Context, job *batchv1.Job) error {
	if m.ShouldFail {
		return ErrMockFailure
	}
	m.NotifiedJobs = append(m.NotifiedJobs, job)
	return nil
}

// MockNotifierFactory implements the NotifierFactory interface for testing
type MockNotifierFactory struct {
	MockNotifier *MockNotifier
}

func NewMockNotifierFactory() *MockNotifierFactory {
	return &MockNotifierFactory{
		MockNotifier: NewMockNotifier(),
	}
}

func (f *MockNotifierFactory) Create(notifierType string) (Notifier, error) {
	return f.MockNotifier, nil
}