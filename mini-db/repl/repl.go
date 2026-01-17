package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"fastabiz-mini-rdbms/mini-db/engine"
	"fastabiz-mini-rdbms/mini-db/storage"
)

type REPL struct {
	engine *engine.Engine
}

func New(engine *engine.Engine) *REPL {
	return &REPL{engine: engine}
}

func (r *REPL) Run() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Fastabiz Mini RDBMS")
	fmt.Println("Type 'exit' to quit")
	fmt.Println()

	for {
		fmt.Print("fastabiz> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}
		if input == "exit" {
			fmt.Println("bye 👋")
			return
		}

		r.handleInput(input)
	}
}

func (r *REPL) handleInput(input string) {
	tokens, err := engine.Tokenize(input)
	if err != nil {
		fmt.Println("token error:", err)
		return
	}

	parser := engine.NewParser(tokens)
	cmd, err := parser.Parse()
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	if err := r.execute(cmd); err != nil {
		fmt.Println("exec error:", err)
	}
}

func (r *REPL) execute(cmd any) error {
	fmt.Println("EXECUTE CALLED")

	switch c := cmd.(type) {

	case *engine.CreateTableCommand:
		if err := r.engine.CreateTable(*c); err != nil {
			return err
		}
		fmt.Println("OK")

	case *engine.InsertCommand:
		err := r.engine.Insert(*c)
		if err != nil {
			return err
		}
		fmt.Println("OK")

	case *engine.SelectCommand:
		rows, err := r.engine.Select(*c)
		if err != nil {
			return err
		}
		printRows(rows)

	case *engine.DeleteCommand:
		n, err := r.engine.Delete(c)
		if err != nil {
			return err
		}
		fmt.Printf("%d row(s) deleted\n", n)

	case *engine.UpdateCommand:
		n, err := r.engine.Update(*c)
		if err != nil {
			return err
		}
		fmt.Printf("%d row(s) updated\n", n)

	default:
		return fmt.Errorf("unknown command")
	}

	return nil
}

func printRows(rows []storage.Row) {
	if len(rows) == 0 {
		fmt.Println("(0 rows)")
		return
	}
	for _, row := range rows {
		fmt.Print("{ ")
		for k, v := range row {
			fmt.Printf("%s:%v ", k, v)
		}
		fmt.Println("}")
	}
}
