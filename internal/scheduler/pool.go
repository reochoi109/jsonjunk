package scheduler

import (
	"container/heap"
	"context"
	logger "jsonjunk/pkg/logging"
	"time"
)

func run(ctx context.Context, pq *PriorityQueue) {
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				logger.Log.Info("scheduler pool close")
				return
			case <-ticker.C:
				now := time.Now()
				for pq.Len() > 0 {
					task := (*pq)[0]
					if task.ExecuteAt.After(now) {
						break
					}

					heap.Pop(pq)
					go task.Action(task.Ctx)

					if task.Interval > 0 {
						task.ExecuteAt = task.ExecuteAt.Add(task.Interval)
						heap.Push(pq, task)
					}
				}
			}
		}
	}()
}
