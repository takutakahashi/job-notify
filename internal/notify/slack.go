package notify

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
	batchv1 "k8s.io/api/batch/v1"
)

// SlackNotifier implements the Notifier interface for Slack
type SlackNotifier struct {
	client     *slack.Client
	channel    string
}

func NewSlackNotifier() *SlackNotifier {
	return &SlackNotifier{}
}

func (s *SlackNotifier) Initialize(ctx context.Context, config map[string]string) error {
	token, ok := config["token"]
	if !ok {
		return fmt.Errorf("slack token is required")
	}
	
	channel, ok := config["channel"]
	if !ok {
		return fmt.Errorf("slack channel is required")
	}

	s.client = slack.New(token)
	s.channel = channel
	return nil
}

func (s *SlackNotifier) NotifyJobStatus(ctx context.Context, job *batchv1.Job) error {
	if job == nil {
		return fmt.Errorf("job is nil")
	}

	// Format the job status message
	var status string
	failed := false
	if job.Status.Succeeded > 0 {
		status = "succeeded"
	} else if job.Status.Failed > 0 {
		status = "failed"
		failed = true
	} else {
		status = "in progress"
	}

	message := fmt.Sprintf("Job '%s' in namespace '%s' %s", job.Name, job.Namespace, status)
	
	// Add completion time if available
	if job.Status.CompletionTime != nil {
		message += fmt.Sprintf(" at %s", job.Status.CompletionTime.Format("2006-01-02 15:04:05"))
	}

	// Create message attachments for more details
	attachment := slack.Attachment{
		Color: s.getStatusColor(failed),
		Fields: []slack.AttachmentField{
			{
				Title: "Status",
				Value: status,
				Short: true,
			},
			{
				Title: "Namespace",
				Value: job.Namespace,
				Short: true,
			},
		},
	}

	// Post the message to Slack
	_, _, err := s.client.PostMessageContext(ctx,
		s.channel,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		return fmt.Errorf("failed to send message to slack: %w", err)
	}

	return nil
}

func (s *SlackNotifier) getStatusColor(failed bool) string {
	if failed {
		return "danger"  // red for failed
	}
	return "good"       // green for success
}