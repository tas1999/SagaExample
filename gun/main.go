package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	step1 := Step[string]{
		Name: "Step1",
		Run: func(data string) (string, error) {
			time.Sleep(time.Second * 1)
			if data == "start" {
				return "step1Finish", nil
			}
			return "Fail", nil
		},
	}
	step2 := Step[string]{
		Name: "Step2",
		Run: func(data string) (string, error) {
			time.Sleep(time.Second * 1)
			if data == "step1Finish" {
				return "step2Finish", nil
			}
			return "Fail", nil
		},
	}
	s := BaseSemulation{
		ThreadCount: 2,
	}
	g := Gun[string]{
		Scenario: []Step[string]{
			step1,
			step2,
		},
		Duration:   time.Second * 10,
		Simulation: s.Simulation,
	}
	g.Start("start")
}

type Step[T any] struct {
	Name string
	Run  func(data T) (T, error)
}

type Gun[T any] struct {
	Scenario   []Step[T]
	Duration   time.Duration
	Simulation func(scenario func() error)
}

func (t *Gun[T]) Start(data T) {
	cont, cancel := context.WithTimeout(context.Background(), t.Duration)
	defer cancel()
	t.Simulation(func() error {
		t.StartScenario(data, cont)
		return nil
	})
}
func (t *Gun[T]) StartScenario(data T, cont context.Context) {
	for {
		result := data
		var err error
		for _, s := range t.Scenario {
			select {
			case <-cont.Done():
				return
			default:
				fmt.Println("Step " + s.Name + "start")
				result, err = s.Run(result)
				if err != nil {
					fmt.Println("Step " + s.Name + "fail with error:" + err.Error())
					break
				}
				fmt.Println("Step " + s.Name + "finish")
			}
		}
	}
}

type BaseSemulation struct {
	ThreadCount int
}

func (s *BaseSemulation) Simulation(scenario func() error) {
	var wg sync.WaitGroup
	for ti := 0; ti < s.ThreadCount; ti++ {
		wg.Add(1)
		go func() {
			scenario()
			wg.Done()
		}()
	}
	wg.Wait()
}
