package record

import (
	"testing"

	F "github.com/ibm/fp-go/function"
	O "github.com/ibm/fp-go/option"
	"github.com/stretchr/testify/assert"
)

type MapType = map[string]int
type MapTypeString = map[string]string
type MapTypeO = map[string]O.Option[int]

func TestSimpleTraversalWithIndex(t *testing.T) {

	f := func(k string, n int) O.Option[int] {
		if k != "a" {
			return O.Some(n)
		}
		return O.None[int]()
	}

	tWithIndex := TraverseWithIndex(
		O.Of[MapType],
		O.Map[MapType, func(int) MapType],
		O.Ap[MapType, int],
		f)

	assert.Equal(t, O.None[MapType](), F.Pipe1(MapType{"a": 1, "b": 2}, tWithIndex))
	assert.Equal(t, O.Some(MapType{"b": 2}), F.Pipe1(MapType{"b": 2}, tWithIndex))
}

func TestSimpleTraversalNoIndex(t *testing.T) {

	f := func(k string) O.Option[string] {
		if k != "1" {
			return O.Some(k)
		}
		return O.None[string]()
	}

	tWithoutIndex := Traverse(
		O.Of[MapTypeString],
		O.Map[MapTypeString, func(string) MapTypeString],
		O.Ap[MapTypeString, string],
		f)

	assert.Equal(t, O.None[MapTypeString](), F.Pipe1(MapTypeString{"a": "1", "b": "2"}, tWithoutIndex))
	assert.Equal(t, O.Some(MapTypeString{"b": "2"}), F.Pipe1(MapTypeString{"b": "2"}, tWithoutIndex))
}

func TestSequence(t *testing.T) {
	// source map
	simpleMapO := MapTypeO{"a": O.Of(1), "b": O.Of(2)}
	// convert to an option of record

	s := Traverse(
		O.Of[MapType],
		O.Map[MapType, func(int) MapType],
		O.Ap[MapType, int],
		F.Identity[O.Option[int]],
	)

	assert.Equal(t, O.Of(MapType{"a": 1, "b": 2}), F.Pipe1(simpleMapO, s))

	s1 := Sequence(
		O.Of[MapType],
		O.Map[MapType, func(int) MapType],
		O.Ap[MapType, int],
		simpleMapO,
	)

	assert.Equal(t, O.Of(MapType{"a": 1, "b": 2}), s1)
}
