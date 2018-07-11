package main

import (
	"log"
	"sync"
	"time"
)

var (
	wg   sync.WaitGroup
	num  int
	jobs chan *Job
	pool chan *Worker
)

type Job struct {
	id int
	do func() error
}

type Worker struct {
	id int
}

func NewWorker(i int) *Worker {
	return &Worker{
		id: i,
	}
}

func (self *Worker) run(job *Job) {
	job.do()
	log.Printf("worker %d finished job %d", self.id, job.id)
	wg.Done()
	pool <- self
}

func start() {
	for {
		job := <-jobs
		log.Printf("receive job %d", job.id)
		worker := <-pool
		log.Printf("get worker %d", worker.id)
		go worker.run(job)
		go log.Printf("jobs(%d), pool(%d)", len(jobs), len(pool))
	}
}

func init() {
	num = 3
	pool = make(chan *Worker, num)
	jobs = make(chan *Job, num*5)
	for i := 0; i < num; i++ {
		pool <- NewWorker(i)
	}
	go start()
}

func main() {
	start := time.Now()

	for i := 0; i < 3*10; i++ {
		id := i
		wg.Add(1)
		jobs <- &Job{
			id: id,
			do: func() error {
				time.Sleep(time.Second)
				return nil
			},
		}
		log.Printf("send job %d", id)
	}

	wg.Wait()
	d := time.Now().Sub(start)
	log.Println(d)
}
