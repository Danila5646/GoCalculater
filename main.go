package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func calculate(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")

	get_precedence := func(operator rune) int {
		switch operator {
		case '+', '-':
			return 1
		case '*', '/':
			return 2
		}
		return 0
	}

	perform_operation := func(op1, op2 float64, operator rune) (float64, error) {
		switch operator {
		case '+':
			return op1 + op2, nil
		case '-':
			return op1 - op2, nil
		case '*':
			return op1 * op2, nil
		case '/':
			if op2 == 0 {
				return 0, errors.New("Деление на ноль!")
			}
			return op1 / op2, nil
		}
		return 0, errors.New("Неизвестный оператор!")
	}

	execute_operations := func(operands []float64, operators []rune) ([]float64, []rune, error) {
		if len(operands) < 2 || len(operators) == 0 {
			return operands, operators, errors.New("Некорректное выражение!")
		}
		op2 := operands[len(operands)-1]
		op1 := operands[len(operands)-2]
		operator := operators[len(operators)-1]
		operands = operands[:len(operands)-2]
		operators = operators[:len(operators)-1]
		result, err := perform_operation(op1, op2, operator)
		if err != nil {
			return operands, operators, err
		}
		return append(operands, result), operators, nil
	}

	var operands []float64
	var operators []rune
	for i := 0; i < len(expression); {
		char := expression[i]
		if (char >= '0' && char <= '9') || char == '.' {
			j := i
			for j < len(expression) && ((expression[j] >= '0' && expression[j] <= '9') || expression[j] == '.') {
				j++
			}
			num, err := strconv.ParseFloat(expression[i:j], 64)
			if err != nil {
				return 0, err
			}
			operands = append(operands, num)
			i = j
		} else if char == '(' {
			operators = append(operators, rune(char))
			i++
		} else if char == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				var err error
				operands, operators, err = execute_operations(operands, operators)
				if err != nil {
					return 0, err
				}
			}
			if len(operators) == 0 {
				return 0, errors.New("Некорректно расставлены скобки!")
			}
			operators = operators[:len(operators)-1]
			i++
		} else {
			for len(operators) > 0 && get_precedence(operators[len(operators)-1]) >= get_precedence(rune(char)) {
				var err error
				operands, operators, err = execute_operations(operands, operators)
				if err != nil {
					return 0, err
				}
			}
			operators = append(operators, rune(char))
			i++
		}
	}

	for len(operators) > 0 {
		var err error
		operands, operators, err = execute_operations(operands, operators)
		if err != nil {
			return 0, err
		}
	}

	if len(operands) != 1 {
		return 0, errors.New("Некорректное выражение!")
	}
	return operands[0], nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Метод запрещен!", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.Expression == "" {
		http.Error(w, "Некорректное выражение!", http.StatusUnprocessableEntity)
		return
	}
	result, err := calculate(request.Expression)
	if err != nil {
		http.Error(w, "Некорректное выражение!", http.StatusUnprocessableEntity)
		return
	}
	response := map[string]string{"result": fmt.Sprintf("%f", result)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/v1/calculate", handler)
	fmt.Println("Сервер запущен!")
	http.ListenAndServe(":8080", nil)
}
