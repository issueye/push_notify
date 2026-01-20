package services

import (
	"backend/internal/models"
)

type PushNotifyJob struct {
	Repo     *models.Repo
	Target   *models.Target
	Payload  *UnifiedPushPayload
	Template *models.Template
	Provider WebhookProvider
}

type PushNotifyQueue struct {
	jobs    chan PushNotifyJob
	handler func(PushNotifyJob)
}

func NewPushNotifyQueue(bufferSize int, workerCount int, handler func(PushNotifyJob)) *PushNotifyQueue {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	if workerCount <= 0 {
		workerCount = 5
	}

	q := &PushNotifyQueue{
		jobs:    make(chan PushNotifyJob, bufferSize),
		handler: handler,
	}

	for i := 0; i < workerCount; i++ {
		go func() {
			for job := range q.jobs {
				q.handler(job)
			}
		}()
	}

	return q
}

func (q *PushNotifyQueue) Enqueue(job PushNotifyJob) bool {
	select {
	case q.jobs <- job:
		return true
	default:
		return false
	}
}
