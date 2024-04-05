package main

import (
	"calculator/stack"
	"fmt"
	"log"
	"os"
	"strconv"
)

const tracing = false

func main() {
	for {
		fmt.Print("> ")
		inputBytes := make([]byte, 64)
		n, _ := os.Stdin.Read(inputBytes)
		inputBytes = inputBytes[:n-2]

		if len(inputBytes) == 0 {
			continue
		}
		if inputBytes[0] == 4 {
			break
		}

		inputString := string(inputBytes)
		fmt.Println(Eval(inputString))
	}
}

func trace(a ...any) {
	if tracing {
		log.Println(a...)
	}
}

func Eval(s string) int {
	tokens := make([]any, 0, 64)
	lex(s, &tokens)

	tokenStack := stack.StackFromSlice(tokens)

	// for _, s := range tokens {
	// 	switch s := s.(type) {
	// 	case int:
	// 		fmt.Println(s)
	// 	case rune:
	// 		fmt.Println(string(s))
	// 	}
	// }

	parseQueue := stack.NewQueue[any]()

	expr(tokenStack, parseQueue)

	numStack := stack.NewStack[int]()

	for !parseQueue.IsEmpty() {
		switch x := parseQueue.Pop().(type) {
		case int:
			numStack.Push(x)
		case rune:
			switch x {
			case '+':
				a := numStack.Pop()
				b := numStack.Pop()
				numStack.Push(b + a)
			case '-':
				a := numStack.Pop()
				b := numStack.Pop()
				numStack.Push(b - a)
			case '/':
				a := numStack.Pop()
				b := numStack.Pop()
				numStack.Push(b / a)
			case '*':
				a := numStack.Pop()
				b := numStack.Pop()
				numStack.Push(b * a)
			case '~':
				numStack.Push(-numStack.Pop())
			}
		}
	}

	return numStack.Pop()
}

func match(tokens *stack.Stack[any], m any) {
	if tokens.IsEmpty() || tokens.Peek() != m {
		switch m := m.(type) {
		case int:
			switch n := tokens.Peek().(type) {
			case int:
				log.Fatal("Parse error in match ", m, ". Next symbol was ", tokens.Peek())
			case rune:
				log.Fatal("Parse error in match ", m, ". Next symbol was ", string(n))
			}
		case rune:
			switch n := tokens.Peek().(type) {
			case int:
				log.Fatal("Parse error in match ", string(m), ". Next symbol was ", n)
			case rune:
				log.Fatal("Parse error in match ", string(m), ". Next symbol was ", string(n))
			}
		}
	}
	tokens.Pop()
}

func expr(tokens *stack.Stack[any], queue *stack.Queue[any]) {
	trace("Entering expr")
	factor(tokens, queue)
exprLoop:
	for {
		if tokens.IsEmpty() {
			break
		}

		switch tokens.Peek() {
		case '+':
			match(tokens, '+')
			factor(tokens, queue)
			queue.Push('+')
			// exprRest(tokens, queue)
			continue
		case '-':
			match(tokens, '-')
			factor(tokens, queue)
			queue.Push('-')
			// exprRest(tokens, queue)
			continue
		default:
			break exprLoop
		}
	}
}

func factor(tokens *stack.Stack[any], queue *stack.Queue[any]) {
	trace("Entering factor")
	term(tokens, queue)
factorLoop:
	for {
		if tokens.IsEmpty() {
			break
		}

		switch tokens.Peek() {
		case '*':
			match(tokens, '*')
			term(tokens, queue)
			queue.Push('*')
			continue
		case '/':
			match(tokens, '/')
			term(tokens, queue)
			queue.Push('/')
			continue
		default:
			break factorLoop
		}
	}
}

func term(tokens *stack.Stack[any], queue *stack.Queue[any]) {
	trace("Entering term")
	if tokens.Peek() == '~' {
		match(tokens, '~')
		term(tokens, queue)
		queue.Push('~')
	} else {
		atom(tokens, queue)
	}
}

func atom(tokens *stack.Stack[any], queue *stack.Queue[any]) {
	trace("Entering atom")
	switch num := tokens.Peek().(type) {
	case int:
		queue.Push(num)
		match(tokens, num)
	case rune:
		match(tokens, '(')
		expr(tokens, queue)
		match(tokens, ')')
	default:
		log.Fatal("Parse error at term")
	}
}

func lex(s string, out *[]any) {
	readingNumber := false
	numberStart := 0

	finishNumber := func(finalLocation int) {
		num, err := strconv.ParseInt(s[numberStart:finalLocation], 10, 32)
		if err != nil {
			log.Fatal("Lex error at '", string(s[finalLocation-1]), "'!")
		}
		// tokenStack.Push(int(num))
		*out = append(*out, int(num))
		readingNumber = false
	}

	for i, r := range s {
		switch r {
		case '(', ')', '*', '/', '+', '-', '~':
			// tokenStack.Push(r)
			if readingNumber {
				finishNumber(i)
			}
			*out = append(*out, r)
			continue
		case ' ', '\n':
			if readingNumber {
				finishNumber(i)
			}
			continue
		default:
			if readingNumber {
				continue
			}
			readingNumber = true
			numberStart = i
			// num, err := strconv.ParseInt(strings.Split(s[i:], " ")[0], 10, 32)
			continue
		}
	}
	if readingNumber {
		finishNumber(len(s))
	}
}
