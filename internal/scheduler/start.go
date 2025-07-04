package scheduler

import (
	"container/heap"
	"context"
	"fmt"
	"time"
)

func StartScheduler(ctx context.Context, pq *PriorityQueue) {
	// ticker := time.NewTicker(100 * time.Millisecond)
	// defer ticker.Stop()
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("END")
				return
			case <-ticker.C:
				now := time.Now()
				for pq.Len() > 0 {
					task := (*pq)[0]
					if task.ExecuteAt.After(now) {
						break
					}

					heap.Pop(pq)
					go task.Action(ctx)

					if task.Interval > 0 {
						task.ExecuteAt = task.ExecuteAt.Add(task.Interval)
						heap.Push(pq, task)
					}
				}
			}
		}
	}()
}
