package scheduler

import (
	"container/heap"
	"context"
	"time"
)

type Task struct {
	Value     string
	ExecuteAt time.Time                 // 다음 실행 시간 (우선순위 기준)
	Interval  time.Duration             // 반복 주기 (0이면 단발성)
	Index     int                       // 내부 heap 인덱스 (Fix나 Remove 시 사용 가능)
	Action    func(ctx context.Context) // 실행 함수
	Ctx       context.Context           // context
}
type PriorityQueue []*Task

var Scheduler *PriorityQueue

func Open(ctx context.Context) {
	Scheduler = &PriorityQueue{}
	heap.Init(Scheduler)
	run(ctx, Scheduler)
}

func Register(item *Task) {
	heap.Push(Scheduler, item)
}

func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].ExecuteAt.Before((*pq)[j].ExecuteAt)
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].Index = i
	(*pq)[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*Task)
	item.Index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	item.Index = -1
	return item
}
