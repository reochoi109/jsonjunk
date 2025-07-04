```go


func main(){
    pq := &scheduler.PriorityQueue{}
	ctx, cancel := context.WithCancel(context.Background())
	scheduler.StartScheduler(ctx, pq)

	now := time.Now()

	scheduler.Register(
		&scheduler.Task{
			Value:     "Task A",
			ExecuteAt: now.Add(1 * time.Second),
			Interval:  5 * time.Second,
			Action: func() {
				fmt.Println("Task A executed at", time.Now())
			},
		},
	)

	scheduler.Register(
		&scheduler.Task{
			Value:     "Task B",
			ExecuteAt: now.Add(2 * time.Second),
			Interval:  10 * time.Second,
			Action: func() {
				fmt.Println("Task B executed at", time.Now())
			},
		},
	)

}


```