package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	OP_ADD = "add"
	OP_SUB = "subtract"
	OP_MULT = "multiply"
	OP_DIV = "divide"
)

type Task struct {
	ID            int `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
}

type TaskResult struct {
	ID int `json:"id"`
	Result float64 `json:"result"`
}

var (
	orchestratorURL string
	computingPower  int
	sem             chan struct{}
	wg              sync.WaitGroup
)

func init() {
	orchestratorURL = os.Getenv("ORCH_URL")
	if orchestratorURL == "" {
		fmt.Println("ORCH_URL environment variable is not set")
		os.Exit(1)
	}
	computingPowerEnv := os.Getenv("COMPUTING_POWER")
	var err error
	computingPower, err = strconv.Atoi(computingPowerEnv)
	if err != nil || computingPower <= 0 {
		fmt.Println("COMPUTING_POWER environment variable is not set or invalid")
		os.Exit(1)
	}
	sem = make(chan struct{}, computingPower)
}

func main() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			statusCode, err := processTask()
			if err != nil && statusCode != http.StatusNotFound {
				fmt.Printf("Error processing task (HTTP %d): %v\n", statusCode, err)
			}
			<-sem
		}()
	}
	wg.Wait()
}

func processTask() (int, error) {
	resp, err := http.Get(orchestratorURL)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to fetch task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return http.StatusNotFound, nil
	}
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return http.StatusUnprocessableEntity, fmt.Errorf("failed to decode task: %v", err)
	}
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	result, err := executeTask(task)
	if err != nil {
		return http.StatusUnprocessableEntity, fmt.Errorf("failed to execute task: %v", err)
	}
	statusCode, err := sendResult(task.ID, result)
	if err != nil {
		return statusCode, fmt.Errorf("failed to send result: %v", err)
	}
	return statusCode, nil
}

func executeTask(task Task) (float64, error) {
	arg1 := task.Arg1
	arg2 := task.Arg2
	var result float64
	switch task.Operation {
	case OP_ADD:
		result = arg1 + arg2
	case OP_SUB:
		result = arg1 - arg2
	case OP_MULT:
		result = arg1 * arg2
	case OP_DIV:
		if arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result = arg1 / arg2
	default:
		return 0, fmt.Errorf("unsupported operation: %s", task.Operation)
	}
	return result, nil
}

func sendResult(taskID int, result float64) (int, error) {
    resultData := TaskResult{
        ID: taskID,
        Result: result,
    }
    resultJSON, err := json.Marshal(resultData)
    if err != nil {
        return http.StatusInternalServerError, fmt.Errorf("failed to marshal result: %v", err)
    }
    resp, err := http.Post(orchestratorURL, "application/json", bytes.NewBuffer(resultJSON))
    if err != nil {
        return http.StatusInternalServerError, fmt.Errorf("failed to send result: %v", err)
    }
    defer resp.Body.Close()
    fmt.Println("Result sent:", resultData)  // Отладочное сообщение
    if resp.StatusCode != http.StatusOK {
        return resp.StatusCode, fmt.Errorf("unexpected status code when sending result: %d", resp.StatusCode)
    }
    return http.StatusOK, nil
}

