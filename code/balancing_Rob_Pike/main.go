package main

import (
	"container/heap"
	"fmt"
	"math/rand/v2"
	"time"
)

/*
Rob Pike演讲中的样例
https://go.dev/talks/2012/waza.slide#54
*/

type Request[T any] struct {
	fn func() T
	c  chan T
}

type Respond[T any] struct {
	i      int
	result T
}

type Worker[T any] struct {
	load  int
	index int
	c     chan Request[T]
}

func NewWorker[T any]() *Worker[T] {
	return &Worker[T]{
		c: make(chan Request[T], 100),
	}
}

func (s *Worker[T]) loop(ch chan *Worker[T]) {
	for req := range s.c {
		ret := req.fn()
		req.c <- ret
		ch <- s
	}
}

type Heap[T any] []*Worker[T]

func (s *Heap[T]) Push(x any) {
	*s = append(*s, x.(*Worker[T]))
	(*s)[s.Len()-1].index = s.Len() - 1
}

func (s *Heap[T]) Pop() (v any) {
	*s, v = (*s)[:s.Len()-1], (*s)[s.Len()-1]
	return
}

func (s *Heap[T]) Len() int {
	return len(*s)
}

func (s *Heap[T]) Less(i, j int) bool {
	return (*s)[i].load < (*s)[j].load
}

func (s *Heap[T]) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	(*s)[i].index = i
	(*s)[j].index = j
}

type Balancing[T any] struct {
	heap *Heap[T]
	ch   chan *Worker[T]
}

func NewBalancing[T any](size int) *Balancing[T] {
	h := &Heap[T]{}
	ch := make(chan *Worker[T], size)

	for i := 0; i < size; i++ {
		work := NewWorker[T]()
		heap.Push(h, work)

		go work.loop(ch)
	}
	return &Balancing[T]{
		heap: h,
		ch:   ch,
	}
}

func (s *Balancing[T]) Run(in <-chan Request[T]) {
	for {
		select {
		case r := <-in:
			readyWork := s.dispatch()
			select {
			case readyWork.c <- r:
			case finishWork := <-s.ch:
				s.finish(finishWork)
			}
		case finishWork := <-s.ch:
			s.finish(finishWork)
		}
	}
}

func (s *Balancing[T]) dispatch() *Worker[T] {
	w := heap.Pop(s.heap).(*Worker[T])
	w.load++
	heap.Push(s.heap, w)
	return w
}

func (s *Balancing[T]) finish(w *Worker[T]) {
	heap.Remove(s.heap, w.index)
	w.load--
	heap.Push(s.heap, w)
}

func main() {
	ba := NewBalancing[Respond[int]](5)

	in := make(chan Request[Respond[int]], 100)
	out := make(chan Respond[int], 100)
	go ba.Run(in)

	go func() {
		for i := 0; i < 10000; i++ {
			j := i
			in <- Request[Respond[int]]{
				fn: func() Respond[int] {
					time.Sleep(time.Microsecond * time.Duration(rand.IntN(500)))
					return Respond[int]{i: j, result: j * j}
				},
				c: out,
			}

			time.Sleep(time.Microsecond * time.Duration(rand.IntN(50)))
		}
	}()

	go func() {
		for r := range out {
			fmt.Printf("i %d\tresult %d\n", r.i, r.result)
		}
	}()

	select {}
}
