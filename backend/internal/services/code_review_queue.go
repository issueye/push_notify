package services

import (
	"strconv"
	"sync"
)

type CodeReviewJob struct {
	RepoID   uint
	PushID   uint
	CommitID string
	Branch   string
}

type CodeReviewQueue struct {
	jobs     chan CodeReviewJob
	handler  func(CodeReviewJob)
	mu       sync.Mutex
	inFlight map[string]struct{}
}

func NewCodeReviewQueue(bufferSize int, workerCount int, handler func(CodeReviewJob)) *CodeReviewQueue {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	if workerCount <= 0 {
		workerCount = 2
	}

	q := &CodeReviewQueue{
		jobs:     make(chan CodeReviewJob, bufferSize),
		handler:  handler,
		inFlight: make(map[string]struct{}),
	}

	for i := 0; i < workerCount; i++ {
		go func() {
			for job := range q.jobs {
				q.handler(job)
				q.finish(job)
			}
		}()
	}

	return q
}

func (q *CodeReviewQueue) Enqueue(job CodeReviewJob) bool {
	key := q.key(job)
	q.mu.Lock()
	if _, ok := q.inFlight[key]; ok {
		q.mu.Unlock()
		return false
	}
	q.inFlight[key] = struct{}{}
	q.mu.Unlock()

	select {
	case q.jobs <- job:
		return true
	default:
		q.finish(job)
		return false
	}
}

func (q *CodeReviewQueue) finish(job CodeReviewJob) {
	q.mu.Lock()
	delete(q.inFlight, q.key(job))
	q.mu.Unlock()
}

func (q *CodeReviewQueue) key(job CodeReviewJob) string {
	return fmtKey(job.RepoID, job.CommitID)
}

func fmtKey(repoID uint, commitID string) string {
	return strconv.FormatUint(uint64(repoID), 10) + ":" + commitID
}
