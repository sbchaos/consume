package example_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"

	"github.com/sbchaos/consume/comb"
	"github.com/sbchaos/consume/comb/list"
	"github.com/sbchaos/consume/comb/maps"
	b "github.com/sbchaos/consume/par"
	"github.com/sbchaos/consume/par/char"
	strings2 "github.com/sbchaos/consume/par/strings"
	ss "github.com/sbchaos/consume/stream/strings"
)

type Json interface {
	ToString() string
}

type JNull struct{}

func (j *JNull) ToString() string {
	return "null"
}

func ParseJNull() b.Parser[rune, Json] {
	return comb.FMap(func(_ string) Json {
		return &JNull{}
	}, strings2.String("null", strings2.EqualIgnoreCase))
}

type JBool struct {
	value bool
}

func (j *JBool) ToString() string {
	return fmt.Sprintf("%t", j.value)
}

func ParseJBool() b.Parser[rune, Json] {
	return comb.FMap(func(a string) Json {
		return &JBool{value: strings.EqualFold(a, "TRUE")}
	}, strings2.Choice([]string{"true", "false"}, strings2.EqualIgnoreCase))
}

type JString struct {
	value string
}

func (j *JString) ToString() string {
	return fmt.Sprintf("\"%s\"", j.value)
}

func ParseJString() b.Parser[rune, Json] {
	return comb.FMap(func(a string) Json {
		return &JString{value: a}
	}, QuotedString())
}

func QuotedString() b.Parser[rune, string] {
	return strings2.QuotedString(0, struct {
		Start rune
		End   rune
	}{Start: '"', End: '"'})
}

func Spaces() b.Parser[rune, string] {
	return comb.FMap(func(a rune) string {
		return ""
	}, char.WhiteSpaces())
}

type JNumber struct {
	value float64
}

func (j *JNumber) ToString() string {
	return fmt.Sprintf("%v", j.value)
}

func ParseJNumber() b.Parser[rune, Json] {
	return comb.FMap1(func(a string) (Json, error) {
		num, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return &JNumber{}, err
		}

		return &JNumber{value: num}, nil
	}, strings2.CustomString(func(a rune) bool {
		return unicode.IsDigit(a) || a == '.'
	}))
}

type JArray struct {
	values []string // simplify list to just strings
}

func (j *JArray) ToString() string {
	b := new(strings.Builder)

	b.WriteString("[")

	for index, value := range j.values {
		b.WriteString(value)

		if index < len(j.values)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("]")

	return b.String()
}

func ParseArray() b.Parser[rune, Json] {
	return comb.FMap(func(a []string) Json {
		return &JArray{values: a}
	}, list.List(QuotedString(), Spaces()))
}

type JObject struct {
	values map[string]string
}

func (j *JObject) ToString() string {
	sb := new(strings.Builder)

	sb.WriteString("{")
	index := 0

	for key, value := range j.values {
		ks := fmt.Sprintf("%q", key)

		sb.WriteString(ks)
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%q", value))

		index++

		if index < len(j.values) {
			sb.WriteString(",")
		}
	}

	sb.WriteString("}")

	return sb.String()
}

func ParseJObject() b.Parser[rune, Json] {
	return comb.FMap(func(a map[string]string) Json {
		return &JObject{values: a}
	}, maps.ObjectLiteral(QuotedString(), QuotedString(), Spaces()))
}

func JsonParser(t *testing.T) b.Parser[rune, Json] {
	boolP := b.Trace(t, "bool", ParseJBool())
	null := b.Trace(t, "null", ParseJNull())
	num := b.Trace(t, "number", ParseJNumber())
	str := b.Trace(t, "string", ParseJString())
	arr := b.Trace(t, "array", ParseArray())
	obj := b.Trace(t, "object", ParseJObject())

	return comb.Choice(
		null,
		boolP,
		num,
		str,
		arr,
		obj,
	)
}

func TestJsonParsing(t *testing.T) {
	t.Run("parses Null", func(t *testing.T) {
		input := ss.NewStringStream("null")

		j, err := b.Parse(input, JsonParser(t))
		assert.NoError(t, err)
		assert.Equal(t, "null", j.ToString())
	})

	t.Run("ParseJBool", func(t *testing.T) {
		t.Run("parses true value", func(t *testing.T) {
			input := ss.NewStringStream("true")

			j, err := b.Parse(input, JsonParser(t))
			assert.NoError(t, err)
			assert.Equal(t, "true", j.ToString())
		})
		t.Run("parses false value", func(t *testing.T) {
			input := ss.NewStringStream("false")

			j, err := b.Parse(input, JsonParser(t))
			assert.NoError(t, err)
			assert.Equal(t, "false", j.ToString())
		})
	})
	t.Run("ParseJNumber", func(t *testing.T) {
		t.Run("parses integer value", func(t *testing.T) {
			input := ss.NewStringStream("5")

			j, err := b.Parse(input, JsonParser(t))
			assert.NoError(t, err)
			assert.Equal(t, "5", j.ToString())
		})
		t.Run("parses double value", func(t *testing.T) {
			input := ss.NewStringStream("6.5")

			j, err := b.Parse(input, JsonParser(t))
			assert.NoError(t, err)
			assert.Equal(t, "6.5", j.ToString())
		})
	})
	t.Run("ParseJString parses string value", func(t *testing.T) {
		input := ss.NewStringStream(`"str"`)

		j, err := b.Parse(input, JsonParser(t))
		assert.NoError(t, err)
		assert.Equal(t, "\"str\"", j.ToString())
	})
	t.Run("parses array", func(t *testing.T) {
		input := ss.NewStringStream(`["abc", "def"]`)

		p1 := JsonParser(t)
		j, err := b.Parse(input, p1)
		assert.NoError(t, err)
		assert.Equal(t, "[abc,def]", j.ToString())
	})
	t.Run("parses object", func(t *testing.T) {
		input := ss.NewStringStream(`{"abc": "def"}`)

		p1 := JsonParser(t)
		j, err := b.Parse(input, p1)
		assert.NoError(t, err)
		assert.Equal(t, "{\"abc\":\"def\"}", j.ToString())
	})
}
