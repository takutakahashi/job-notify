package notify

import (
	"context"
	"testing"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSlackNotifier_Initialize(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]string
		wantErr bool
	}{
		{
			name: "valid config",
			config: map[string]string{
				"token":   "test-token",
				"channel": "test-channel",
			},
			wantErr: false,
		},
		{
			name: "missing token",
			config: map[string]string{
				"channel": "test-channel",
			},
			wantErr: true,
		},
		{
			name: "missing channel",
			config: map[string]string{
				"token": "test-token",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSlackNotifier()
			err := s.Initialize(context.Background(), tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlackNotifier_NotifyJobStatus(t *testing.T) {
	completionTime := metav1.Now()
	
	tests := []struct {
		name    string
		job     *batchv1.Job
		wantErr bool
	}{
		{
			name: "nil job",
			job:  nil,
			wantErr: true,
		},
		{
			name: "successful job",
			job: &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-job",
					Namespace: "default",
				},
				Status: batchv1.JobStatus{
					Succeeded:      1,
					CompletionTime: &completionTime,
				},
			},
			wantErr: false,
		},
		{
			name: "failed job",
			job: &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-job",
					Namespace: "default",
				},
				Status: batchv1.JobStatus{
					Failed:         1,
					CompletionTime: &completionTime,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSlackNotifier()
			s.Initialize(context.Background(), map[string]string{
				"token":   "test-token",
				"channel": "test-channel",
			})

			err := s.NotifyJobStatus(context.Background(), tt.job)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotifyJobStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}