package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseSimple(t *testing.T) {
	buf := bytes.NewBufferString(`key: "value"`)
	l := NewLexer(buf)
	p := NewParser(l)
	v := p.Parse()
	if _, ok := v.(*ObjectValue); !ok {
		t.Errorf("expected result is object value. actual type: %T", v)
	}
}

func TestParseString(t *testing.T) {
	buf := bytes.NewBufferString(`"value"`)
	l := NewLexer(buf)
	p := NewParser(l)
	v := p.Parse()
	if _, ok := v.(*StringValue); !ok {
		t.Errorf("expected result is object value. actual type: %T", v)
	}
}

func TestParseNumber(t *testing.T) {
	buf := bytes.NewBufferString(`1234`)
	l := NewLexer(buf)
	p := NewParser(l)
	v := p.Parse()
	if _, ok := v.(*NumberValue); !ok {
		t.Errorf("expected result is object value. actual type: %T", v)
	}
}
func TestParseObject(t *testing.T) {
	text := `
a: "a"
b: "b"
c: "c"
`
	buf := bytes.NewBufferString(text)
	l := NewLexer(buf)
	p := NewParser(l)
	v := p.Parse()
	if v, ok := v.(*ObjectValue); !ok {
		t.Errorf("expected result is object value. actual type: %T", v)
	} else {
		assertPair(t, v.Pair[0], "a", &StringValue{Value: "a"})
		assertPair(t, v.Pair[1], "b", &StringValue{Value: "b"})
		assertPair(t, v.Pair[2], "c", &StringValue{Value: "c"})
	}
}

func TestParseIndent(t *testing.T) {
	text := `
a: 
  b: "b"
  c: "c"
  d: "d"
  e: "e"
`
	buf := bytes.NewBufferString(strings.TrimSpace(text))
	l := NewLexer(buf)
	p := NewParser(l)
	v := p.Parse()
	if v, ok := v.(*ObjectValue); !ok {
		t.Errorf("expected result is object value. actual type: %T, value:%+v", v, v)
	} else {
		assertPair(t, v.Pair[0], "a", &ObjectValue{Pair: []Pair{
			{
				Key: &Key{
					Identifier: &Identifier{Value: "b"},
				},
				Value: &StringValue{Value: "b"},
			},
			{
				Key: &Key{
					Identifier: &Identifier{Value: "c"},
				},
				Value: &StringValue{Value: "c"},
			},
			{
				Key: &Key{
					Identifier: &Identifier{Value: "d"},
				},
				Value: &StringValue{Value: "d"},
			},
			{
				Key: &Key{
					Identifier: &Identifier{Value: "e"},
				},
				Value: &StringValue{Value: "e"},
			},
		}})
	}
}

func TestParseIndent2(t *testing.T) {
	text := `a: 
  b: "b"
c: "c"
`
	buf := bytes.NewBufferString(strings.TrimSpace(text))
	l := NewLexer(buf)
	p := NewParser(l)
	v := p.Parse()
	if v, ok := v.(*ObjectValue); !ok {
		t.Errorf("expected result is object value. actual type: %T, value:%+v", v, v)
	} else {
		assertPair(t, v.Pair[0], "a", &ObjectValue{Pair: []Pair{
			{
				Key: &Key{
					Identifier: &Identifier{Value: "b"},
				},
				Value: &StringValue{Value: "b"},
			},
		}})
		assertPair(t, v.Pair[1], "c", &StringValue{Value: "c"})
	}
}

func assertPair(t *testing.T, pair Pair, expectedKey string, expectedValue Value) {
	if pair.Key.Identifier.Value != expectedKey {
		t.Errorf("unexpected key name. actual: %s, expected: %s", pair.Key.Identifier.Value, expectedKey)
	}
	assertValue(t, pair.Value, expectedValue)
}

func assertValue(t *testing.T, actualValue Value, expectedValue Value) {
	switch expected := expectedValue.(type) {
	case *StringValue:
		if actual, ok := actualValue.(*StringValue); !ok {
			t.Errorf("value is expected type. expected type: StringValue, actual type:%T, actual value: %v", actual, actual)
		} else if actual.Value != expected.Value {
			t.Errorf("value is expected value. expected value: %s, actual:%s", expected.Value, actual.Value)
		}
	case *ObjectValue:
		if actual, ok := actualValue.(*ObjectValue); !ok {
			t.Errorf("value is expected type. expected type: ObjectValue, actual type:%T, actual value: %v", actual, actual)
		} else {
			if len(actual.Pair) != len(expected.Pair) {
				t.Errorf("object keys does not have same length. expected:%d, actual:%d", len(expected.Pair), len(actual.Pair))
				return
			}
			for i, e := range expected.Pair {
				a := actual.Pair[i]
				assertPair(t, a, e.Key.Identifier.Value, e.Value)
			}
		}
	default:
		t.Error("Unsupported expected types")
	}

}
