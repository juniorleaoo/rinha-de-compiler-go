package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type File struct {
	Name       string `json:"name"`
	Expression Term   `json:"expression"`
}

type Int struct {
	Kind  string `json:"kind"`
	Value int32  `json:"value"`
}

type Str struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Bool struct {
	Kind  string `json:"kind"`
	Value bool   `json:"value"`
}

type Binary struct {
	Kind string `json:"kind"`
	Lhs  *Term  `json:"lhs"`
	Op   string `json:"op"`
	Rhs  *Term  `json:"rhs"`
}

type If struct {
	Kind      string `json:"kind"`
	Condition *Term  `json:"condition"`
	Then      *Term  `json:"then"`
	Otherwise *Term  `json:"otherwise"`
}

type Print struct {
	Kind  string `json:"kind"`
	Value Term   `json:"value"`
}

type Value interface{}

type Term struct {
	Kind      string `json:"kind"`
	Value     Value  `json:"value"`
	Op        string `json:"op,omitempty"`
	Lhs       *Term  `json:"lhs,omitempty"`
	Rhs       *Term  `json:"rhs,omitempty"`
	Condition *Term  `json:"condition"`
	Then      *Term  `json:"then"`
	Otherwise *Term  `json:"otherwise"`
}

func (t *Term) UnmarshalJSON(b []byte) error {
	var printValue Print
	var intValue Int
	var strValue Str
	var boolValue Bool
	var binaryValue Binary
	var ifValue If

	var term struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(b, &term); err == nil {
		t.Kind = term.Kind

		switch term.Kind {
		case "If":
			if err := json.Unmarshal(b, &ifValue); err == nil {
				t.Condition = ifValue.Condition
				t.Then = ifValue.Then
				t.Otherwise = ifValue.Otherwise
				return nil
			}

		case "Print":
			if err := json.Unmarshal(b, &printValue); err == nil {
				t.Value = printValue.Value
				return nil
			}

		case "Str":
			if err := json.Unmarshal(b, &strValue); err == nil {
				t.Value = strValue.Value
				return nil
			}

		case "Int":
			if err := json.Unmarshal(b, &intValue); err == nil {
				t.Value = intValue.Value
				return nil
			}

		case "Bool":
			if err := json.Unmarshal(b, &boolValue); err == nil {
				t.Value = boolValue.Value
				return nil
			}

		case "Binary":
			if err := json.Unmarshal(b, &binaryValue); err == nil {
				t.Op = binaryValue.Op
				t.Lhs = binaryValue.Lhs
				t.Rhs = binaryValue.Rhs
				return nil
			}
		}

		return nil
	}

	return fmt.Errorf("Deu ruim no parse")
}

func main() {
	var file File
	astFileBytes, _ := io.ReadAll(os.Stdin)
	err := json.Unmarshal(astFileBytes, &file)
	if err != nil {
		fmt.Printf("Não conseguiu deserializar o json %v", err)
	}
	evaluate(file.Expression)
}

func evaluate(term Term) Value {
	switch term.Kind {
	case "Int":
		return term.Value.(int32)
	case "Str":
		return term.Value.(string)
	case "Bool":
		return term.Value.(bool)
	case "Binary":
		lhs := evaluate(*term.Lhs)
		rhs := evaluate(*term.Rhs)
		switch term.Op {
		case "Add":
			lhsStringValue, lhsString := lhs.(string)
			lhsIntValue, lhsInt := lhs.(int32)
			rhsStringValue, rhsString := rhs.(string)
			rhsIntValue, rhsInt := rhs.(int32)

			if lhsInt && rhsInt {
				return lhsIntValue + rhsIntValue
			} else if lhsString && rhsString {
				return fmt.Sprintf("%v%v", lhsStringValue, rhsStringValue)
			} else if lhsInt && rhsString {
				return fmt.Sprintf("%d%v", lhsIntValue, rhsStringValue)
			} else if lhsString && rhsInt {
				return fmt.Sprintf("%v%d", lhsStringValue, rhsIntValue)
			}

		default:
			fmt.Print("BinaryOp não implementado")
		}

	case "Print":
		val := evaluate(term.Value.(Term))
		switch val.(type) {
		case string:
			fmt.Print(val)
		case int32:
			fmt.Printf("%d", val)
		case bool:
			fmt.Printf("%t", val)
		default:
			fmt.Print("Value não permitido")
		}

	case "If":
		val := evaluate(*term.Condition)
		switch val {
		case true:
			return evaluate(*term.Then)
		case false:
			return evaluate(*term.Otherwise)
		default:
			fmt.Print("Condição não permitida")
		}
	default:
		fmt.Print("Kind não implementado")
	}

	return nil
}
