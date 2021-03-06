package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	delayTime    = time.Tick(time.Millisecond * 500)
	jobScheduler = NewJobScheduler()
	jobPipeline  = NewJobPipeline()
)

type InitResult struct {
	City       string
	Kd         string
	TotalPage  int
	TotalCount int
}

type LoopResult struct {
	Success int
	Error   int
	Empty   int
	Errors  []string
}

func InitJobs(city string, pn int, kd string) ([]InitResult, error) {
	var (
		jobs       []Result
		totalPage  int
		totalCount int
		results    []InitResult

		err error
	)

	jobs, totalPage, totalCount, err = GetJobs(city, pn, kd)
	if err != nil {
		return nil, err
	}

	results = append(results, InitResult{
		City:       city,
		Kd:         kd,
		TotalPage:  totalPage,
		TotalCount: totalCount,
	})

	for i := 2; i <= totalPage; i++ {
		jobScheduler.Append(city, i, kd)
	}

	jobPipeline.Append(ToPipelineJobs(jobs))

	return results, nil
}

func LoopJobs() LoopResult {
	var (
		result LoopResult
		output = jobScheduler.Count()

		params = make(chan []Result)
	)

	for i := 0; i < output; i++ {
		<-delayTime
		go func() {
			if jobParam := jobScheduler.Pop(); jobParam != nil {
				jobs, _, _, err := GetJobs(jobParam.City, jobParam.Pn, jobParam.Kd)
				if err != nil {
					result.Error++
					result.Errors = append(result.Errors, err.Error())
				} else {
					params <- jobs
				}
			} else {
				result.Empty++
			}
		}()
	}

L:
	for {
		select {
		case p := <-params:
			result.Success++
			jobPipeline.Append(ToPipelineJobs(p))
		default:
			if (result.Success + result.Error + result.Empty) >= output {
				log.Printf("Break...")
				break L
			}
		}
	}

	return result
}

func GetJobs(city string, pn int, kd string) ([]Result, int, int, error) {
	totalPage := 0
	jobService := NewJobService(city)
	result, err := jobService.GetJobs(pn, kd)
	if err != nil {
		return nil, 0, 0, err
	}

	log.Printf("GetJobs Code: %d, GetJobs City: %s, Pn: %d, Kd: %s", result.Code, city, pn, kd)

	if result.Code == 0 && result.Success == true {
		content := result.Content
		if content.PositionResult.TotalCount > 0 && content.PageSize > 0 {
			totalPage = CalculateTotalPage(float64(content.PositionResult.TotalCount), float64(content.PageSize))
		}
	} else {
		return nil, 0, 0, errors.New(fmt.Sprintf("GetJobs City: %s, Pn: %d, Kd: %s, Result: %v", city, pn, kd, result))
	}

	return result.Content.PositionResult.Result, totalPage, result.Content.PositionResult.TotalCount, nil
}
