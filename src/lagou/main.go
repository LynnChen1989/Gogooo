package main

import (
	"log"
	"sync"
)

var (
	kds = []string{
		"golang",
	}
	citys = []string{
		"北京",
		"上海",
		"广州",
		"深圳",
		"杭州",
		"成都",
	}

	initResults = []InitResult{}
	loopResults = []LoopResult{}
	//jobPipeline = NewJobPipeline()

	wg sync.WaitGroup
)

func main() {
	for _, kd := range kds {
		for _, city := range citys {
			wg.Add(1)
			go func(city string, kd string) {
				defer wg.Done()
				initResult, err := InitJobs(city, 1, kd)
				if err != nil {
					log.Fatalln(err)
				}

				initResults = append(initResults, initResult...)
				loopResults = append(loopResults, LoopJobs())
			}(city, kd)
		}
	}

	wg.Wait()

	jobPipeline.Push()

	log.Printf("Init Results: %v", initResults)
	log.Printf("Loop Results: %v", loopResults)
}
