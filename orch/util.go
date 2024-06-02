package main

import (
	"fmt"
	"unicode"
	"strconv"
	"sync/atomic"
)

func findSequenceIndex(slice []string, sequence []string) int {
	for i := 0; i <= len(slice)-len(sequence); i++ {
		match := true
		for j := 0; j < len(sequence); j++ {
			if slice[i+j] != sequence[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func replaceFirstSequence(slice []string, sequence []string, replacement string) []string {
	index := findSequenceIndex(slice, sequence)
	if index == -1 {
		return slice
	}
	newSlice := append(slice[:index], append([]string{replacement}, slice[index+len(sequence):]...)...)
	return newSlice
}


func getNextID() int {
    return int(atomic.AddInt32(&IDCounter, 1))
}

func getCurrentID() int {
    return int(atomic.LoadInt32(&IDCounter))
}

func getNextTaskID() int {
    return int(atomic.AddInt32(&TaskIDCounter, 1))
}

func getCurrentTaskID() int {
    return int(atomic.LoadInt32(&TaskIDCounter))
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}


func findTriplets(input []string) [][]string {
	var triplets [][]string
	for i := 0; i < len(input)-2; i++ {
		if isNumber(input[i]) && isNumber(input[i+1]) && isOperator(input[i+2]) {
			triplet := []string{input[i], input[i+1], input[i+2]}
			triplets = append(triplets, triplet)
		}
	}

	return triplets
}


func toPostfix(exp string) ([]string, error) {
	var output []string
	var stack []string
	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	greaterPrecedence := func(op1, op2 string) bool {
		return precedence[op1] >= precedence[op2]
	}
	i := 0
	for i < len(exp) {
		ch := rune(exp[i])

		if unicode.IsSpace(ch) {
			// Ignore spaces
			i++
			continue
		}

		if unicode.IsDigit(ch) || ch == '.' {
			start := i
			for i < len(exp) && (unicode.IsDigit(rune(exp[i])) || exp[i] == '.') {
				i++
			}
			output = append(output, exp[start:i])
			continue
		}

		if ch == '(' {
			stack = append(stack, string(ch))
			i++
			continue
		}

		if ch == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 || stack[len(stack)-1] != "(" {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
			i++
			continue
		}
		if isOperator(string(ch)) {
			for len(stack) > 0 && isOperator(stack[len(stack)-1]) && greaterPrecedence(stack[len(stack)-1], string(ch)) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, string(ch))
			i++
			continue
		}

		return nil, fmt.Errorf("invalid character: %c", ch)
	}
	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" || stack[len(stack)-1] == ")" {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}