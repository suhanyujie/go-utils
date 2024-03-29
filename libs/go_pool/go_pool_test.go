package go_pool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPoolRun1(t *testing.T) {
	t1 := time.Now()
	wg := sync.WaitGroup{}
	// 等待 3 个任务执行。
	wg.Add(3)
	// 大小为 4 的协程池
	pool := NewGoPool(4)
	pool.NewTask(func() {
		fmt.Printf("task1 is executing...\n")
		time.Sleep(time.Second * 3)
		fmt.Printf("task1 is finished.\n")
		wg.Done()
	})
	pool.NewTask(func() {
		fmt.Printf("task2 is executing...\n")
		time.Sleep(time.Second * 3)
		fmt.Printf("task2 is finished.\n")
		wg.Done()
	})
	pool.NewTask(func() {
		fmt.Printf("task3 is executing...\n")
		time.Sleep(time.Second * 3)
		fmt.Printf("task3 is finished.\n")
		wg.Done()
	})
	// 等待所有任务执行完毕
	wg.Wait()
	fmt.Printf("time consume: %.3f\n", time.Since(t1).Seconds())
}

func TestPoolRun2(t *testing.T) {
	t1 := time.Now()
	wg := sync.WaitGroup{}
	// 等待 3 个任务执行。
	wg.Add(3)
	pool := NewGoPool(2)
	for i := 0; i < 3; i++ {
		RunOne(pool, i, &wg)
		fmt.Printf("loop %d\n", i)
	}
	// 等待所有任务执行完毕
	wg.Wait()
	fmt.Printf("time consume: %.3f\n", time.Since(t1).Seconds())
}

func RunOne(pool *GoPool, val int, wg *sync.WaitGroup) {
	pool.NewTask(func() {
		fmt.Printf("task %d is executing...\n", val)
		time.Sleep(time.Second * 3)
		fmt.Printf("task %d is finished.\n", val)
		wg.Done()
	})
}
