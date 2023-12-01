package query

import (
	"bytes"
)

func And(args ...Query) Query {
	return &BinaryQuery{
		Operation: AndOp,
		Args:      args,
	}
}

func Or(args ...Query) Query {
	return &BinaryQuery{
		Operation: OrOp,
		Args:      args,
	}
}

func Exact(field, value string) Query {
	return &ExactQuery{
		Field: field,
		Value: value,
	}
}

func Render(q Query) string {
	var buffer bytes.Buffer
	q.write(&buffer)
	return buffer.String()
}

type Query interface {
	write(*bytes.Buffer)
}

type BinaryOperation int

const (
	AndOp BinaryOperation = iota
	OrOp
)

type BinaryQuery struct {
	Operation BinaryOperation
	Args      []Query
}

func (e *BinaryQuery) write(dst *bytes.Buffer) {
	needParantheses := len(e.Args) > 1
	if needParantheses {
		_ = dst.WriteByte('(')
	}
	var operation string
	switch e.Operation {
	case AndOp:
		// we could specify here ` AND `, but older Anki doesn't understand it
		operation = " "
	case OrOp:
		operation = " OR "
	}
	for i, query := range e.Args {
		query.write(dst)
		if i != len(e.Args)-1 {
			_, _ = dst.WriteString(operation)
		}
	}
	if needParantheses {
		_ = dst.WriteByte(')')
	}
}

type ExactQuery struct {
	Field string
	Value string
}

func (e *ExactQuery) write(dst *bytes.Buffer) {
	_ = dst.WriteByte('"')
	if e.Field != "" {
		escape(dst, e.Field)
		_ = dst.WriteByte(':')
	}
	switch e.Field {
	// seems like anki has different encoding for different values,
	// searching in fields require more strict encoding
	case "deck", "note", "tag", "card":
		escape(dst, e.Value)
	default:
		escapeField(dst, e.Value)
	}
	_ = dst.WriteByte('"')
}

func escape(dst *bytes.Buffer, src string) {
	for i := 0; i < len(src); i++ {
		switch b := src[i]; b {
		case '"':
			_, _ = dst.WriteString(`\"`)
		case '*':
			_, _ = dst.WriteString(`\*`)
		case '_':
			_, _ = dst.WriteString(`\_`)
		case ':':
			_, _ = dst.WriteString(`\:`)
		case '\\':
			_, _ = dst.WriteString(`\\`)
		default:
			_ = dst.WriteByte(b)
		}
	}
}

func escapeField(dst *bytes.Buffer, src string) {
	for i := 0; i < len(src); i++ {
		switch b := src[i]; b {
		case '"':
			_, _ = dst.WriteString(`\"`)
		case '*':
			_, _ = dst.WriteString(`\*`)
		case '_':
			_, _ = dst.WriteString(`\_`)
		case ':':
			_, _ = dst.WriteString(`\:`)
		case '\\':
			_, _ = dst.WriteString(`\\`)
		case '&':
			_, _ = dst.WriteString(`&amp;`)
		case '<':
			_, _ = dst.WriteString(`&lt;`)
		case '>':
			_, _ = dst.WriteString(`&gt;`)
		default:
			_ = dst.WriteByte(b)
		}
	}
}
