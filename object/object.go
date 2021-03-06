package object

import (
	"fmt"
	"github.com/carsonip/monkey-interpreter/ast"
	"hash/fnv"
	"strings"
)

type Object interface {
	String() string
}

type Hashable interface {
	Hash() uint64
}

type Null struct {}

func (n Null) String() string {
	return ""
}

func (n Null) Hash() uint64 {
	return 0
}

var NULL = Null{}

type Integer struct {
	Value int
}

func (i Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i Integer) Hash() uint64 {
	return uint64(i.Value)
}

func NewInteger(value int) Integer {
	return Integer{Value: value}
}

type Boolean struct {
	Value bool
}

func (b Boolean) String() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}

func (b Boolean) Hash() uint64 {
	if b.Value {
		return 0
	} else {
		return 1
	}
}

func NewBoolean(value bool) Boolean {
	return Boolean{Value: value}
}

type Function struct {
	Params []string
	Body []ast.Node
}

func (f Function) String() string {
	return "fn"
}

func NewFunction(params []string, body []ast.Node) Function {
	return Function{Params: params, Body: body}
}

type BuiltinFunction struct {
	Fn func(args ...Object) Object
}

func (f BuiltinFunction) String() string {
	return "builtin"
}

type String struct {
	Value string
}

func (s String) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

func (s String) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return h.Sum64()
}

func NewString(value string) String {
	return String{Value: value}
}

type Array struct {
	Elements []Object
}

func (a Array) String() string {
	var strs []string
	for _, element := range a.Elements {
		strs = append(strs, element.String())
	}

	return fmt.Sprintf(`[%s]`, strings.Join(strs, ", "))
}

func NewArray(elements []Object) Array {
	return Array{Elements: elements}
}

type KV struct {
	Key Object
	Value Object
}

type Map struct {
	Elements map[uint64][]KV
}

func (m Map) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	first := true
	for _, pairs := range m.Elements {
		for _, kv := range pairs {
			if first {
				first = false
			} else {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s: %s", kv.Key.String(), kv.Value.String()))
		}
	}
	sb.WriteString("}")

	return sb.String()
}

func (m Map) Get(key Object) (Object, bool) {
	if hashable, ok := key.(Hashable); !ok {
		panic("key not hashable")
	} else if pairs, ok := m.Elements[hashable.Hash()]; !ok {
		return nil, false
	} else {
		for _, kv := range pairs {
			if kv.Key == key {
				return kv.Value, true
			}
		}
		return nil, false
	}
}

func NewMap(pairs [][2]Object) Map {
	m := Map{Elements: make(map[uint64][]KV)}
	for _, kv := range pairs {
		k := kv[0]
		v := kv[1]
		if hashable, ok := k.(Hashable); !ok {
			panic("key not hashable")
		} else {
			h := hashable.Hash()
			m.Elements[h] = append(m.Elements[h], KV{k, v})
		}
	}
	return m
}
