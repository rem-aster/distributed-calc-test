package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io"
	"strconv"
	"sync"
)

var IDCounter int32
var TaskIDCounter int32

type Task struct {
	ID int `json:"id"`
	Arg1 float64 `json:"arg1"`
	Arg2 float64 `json:"arg2"`
	Operation string `json:"operation"`
	OperationTime int    `json:"operation_time"`
}

type TaskResult struct {
	ID int `json:"id"`
	Result float64 `json:"result"`
}

const (
	OP_ADD = "add"
	OP_SUB = "subtract"
	OP_MULT = "multiply"
	OP_DIV = "divide"

	STATUS_CALC = "calculating"
	STATUS_READY = "ready"
	STATUS_ERROR = "error"
)

type RawExpression struct {
	Expression string `json:"expression"`
}

type ExpressionAcceptedID struct {
	ID int `json:"id"`
}

type ExpressionResult struct {
	ID int `json:"id"`
	Status string `json:"status"`
	Result float64 `json:"result"`
}
type SingleExpression struct {
	Exp ExpressionResult `json:"expression"`
}

type AllExpressions struct {
	Expressions map[int]ExpressionResult `json:"expressions"`
}

type TaskQueue struct {
	items []Task
	mu    sync.Mutex
}

var taskQueue TaskQueue
var resultChannels map[int]chan float64
var expressions AllExpressions
var cache map[string]string
var rcMu sync.Mutex
var cacheMu sync.Mutex
var expMu sync.Mutex
var OpStringsMap map[string]string
var OpTimeMap map[string]int

func (q *TaskQueue) Enqueue(item Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *TaskQueue) Dequeue() (Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var a Task
		return a, fmt.Errorf("queue empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *TaskQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items) == 0
}

func processRawExpression(w http.ResponseWriter, r *http.Request) {
	var rawExp RawExpression
	err := json.NewDecoder(r.Body).Decode(&rawExp)
	if err != nil {
		http.Error(w, "Failed to parse expression json", http.StatusInternalServerError)
		return
	}
	exp, err := toPostfix(rawExp.Expression)
	if err != nil {
		http.Error(w, "Expression is uncorrect", http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusCreated)
	id := getNextID()
	json.NewEncoder(w).Encode(ExpressionAcceptedID{ID: id})
	expMu.Lock()
	expressions.Expressions[id] = ExpressionResult{ID: id, Status: STATUS_CALC, Result: 0}
	expMu.Unlock()
	go processExpression(id, exp)
}

func processExpression(id int, exp []string) {
    expr := exp
    for {
        fmt.Println("Processing expression:", expr)
        if len(expr) == 1 {
            if isNumber(expr[0]) {
                res, _ := strconv.ParseFloat(expr[0], 64)
                expMu.Lock()
                expressions.Expressions[id] = ExpressionResult{ID: id, Status: STATUS_READY, Result: res}
                expMu.Unlock()
                fmt.Println("Expression result ready:", expressions.Expressions[id])  // Отладочное сообщение
                return
            } else {
                expMu.Lock()
                expressions.Expressions[id] = ExpressionResult{ID: id, Status: STATUS_ERROR, Result: 0}
                expMu.Unlock()
                fmt.Println("Error in expression:", expr)  // Отладочное сообщение
                return
            }
        } else if len(expr) > 2 {
            awaited := make(map[string]chan float64)
            for _, v := range findTriplets(expr) {
                cacheMu.Lock()
				if v[1] == "0" && v[2] == "/" {
					expr = []string{}
					fmt.Println("ERR div by 0")
					break
				}
                vc, ok := cache[fmt.Sprintf("%s %s %s", v[0], v[2], v[1])]
                vc2, ok2 := cache[fmt.Sprintf("%s %s %s", v[1], v[2], v[0])]
                if ok {
                    expr = replaceFirstSequence(expr, v, vc)
                } else if ok2 {
                    expr = replaceFirstSequence(expr, v, vc2)
                } else {
                    tID := getNextTaskID()
                    a1, _ := strconv.ParseFloat(v[0], 64)
                    a2, _ := strconv.ParseFloat(v[1], 64)
                    taskQueue.Enqueue(Task{ID: tID, Arg1: a1, Arg2: a2, Operation: OpStringsMap[v[2]], OperationTime: OpTimeMap[v[2]]})
                    ch := make(chan float64)
                    awaited[fmt.Sprintf("%s %s %s", v[0], v[2], v[1])] = ch
                    rcMu.Lock()
                    resultChannels[tID] = ch
                    rcMu.Unlock()
                }
                cacheMu.Unlock()
            }
            for k, v := range awaited {
                res := fmt.Sprintf("%f", <-v)
                close(v)
                cacheMu.Lock()
                cache[k] = res
                cacheMu.Unlock()
            }
        } else {
            expMu.Lock()
            expressions.Expressions[id] = ExpressionResult{ID: id, Status: STATUS_ERROR, Result: 0}
            expMu.Unlock()
            fmt.Println("Error in expression:", expr)  // Отладочное сообщение
            return
        }
    }
}

func getExpressions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type AllExpJSON struct {
		ExpRes []ExpressionResult `json:"expressions"`
	}
	var a AllExpJSON = AllExpJSON{
		ExpRes: make([]ExpressionResult, 0),
	}
	for _, v := range expressions.Expressions {
		a.ExpRes = append(a.ExpRes, v)
	}
	json.NewEncoder(w).Encode(a)
}

func getExpressionByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Path[len("/api/v1/expressions/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	expMu.Lock()
	v, ok := expressions.Expressions[id]
	expMu.Unlock()
	if ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SingleExpression{Exp: v})
		return
	}

	http.NotFound(w, r)
}

func GetPostAgent(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        t, err := taskQueue.Dequeue()
        if err != nil {
            http.NotFound(w, r)
        } else {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(t)
            fmt.Println("Sent task to agent:", t)  // Отладочное сообщение
        }
    case "POST":
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "error reading body", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()
        var res TaskResult
        err = json.Unmarshal(body, &res)
        if err != nil {
            http.Error(w, "error unmarshaling JSON", http.StatusBadRequest)
            return
        }
        rcMu.Lock()
        if ch, ok := resultChannels[res.ID]; ok {
            ch <- res.Result
            rcMu.Unlock()
            w.WriteHeader(http.StatusOK)
            fmt.Println("Received result from agent:", res)  // Отладочное сообщение
        } else {
            rcMu.Unlock()
            http.Error(w, "invalid task ID", http.StatusBadRequest)
        }
    default:
        http.Error(w, "wrong method", http.StatusMethodNotAllowed)
    }
}


func main() {
	T_ADD, err1 := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
    T_SUB, err2 := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
    T_MULT, err3 := strconv.Atoi(os.Getenv("TIME_MULTIPLICATION_MS"))
    T_DIV, err4 := strconv.Atoi(os.Getenv("TIME_DIVISION_MS"))
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		fmt.Println("ERR os env not set")
	}
	OpTimeMap = map[string]int{
		"+": T_ADD,
		"-": T_SUB,
		"*": T_MULT,
		"/": T_DIV,
	}
	OpStringsMap = map[string]string{
		"+": OP_ADD,
		"-": OP_SUB,
		"*": OP_MULT,
		"/": OP_DIV,
	}
	taskQueue.items = make([]Task, 0)
	expressions.Expressions = make(map[int]ExpressionResult)
	cache = make(map[string]string)
	resultChannels = make(map[int]chan float64)
	http.HandleFunc("/api/v1/calculate", processRawExpression)
	http.HandleFunc("/api/v1/expressions", getExpressions)
	http.HandleFunc("/api/v1/expressions/", getExpressionByID)
	http.HandleFunc("/internal/task", GetPostAgent)
	http.ListenAndServe(":8080", nil)
}