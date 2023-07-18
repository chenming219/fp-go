package http

import (
	"fmt"
	"io"
	"mime"
	H "net/http"
	"regexp"

	A "github.com/ibm/fp-go/array"
	E "github.com/ibm/fp-go/either"
	"github.com/ibm/fp-go/errors"
	F "github.com/ibm/fp-go/function"
	O "github.com/ibm/fp-go/option"
	R "github.com/ibm/fp-go/record/generic"
	T "github.com/ibm/fp-go/tuple"
)

type (
	ParsedMediaType = T.Tuple2[string, map[string]string]
)

var (
	// mime type to check if a media type matches
	reJsonMimeType = regexp.MustCompile(`application/(?:\w+\+)?json`)
	// ValidateResponse validates an HTTP response and returns an [E.Either] if the response is not a success
	ValidateResponse = E.FromPredicate(isValidStatus, StatusCodeError)
	// alidateJsonContentTypeString parses a content type a validates that it is valid JSON
	validateJsonContentTypeString = F.Flow2(
		ParseMediaType,
		E.ChainFirst(F.Flow2(
			T.First[string, map[string]string],
			E.FromPredicate(reJsonMimeType.MatchString, func(mimeType string) error {
				return fmt.Errorf("mimetype [%s] is not a valid JSON content type", mimeType)
			}),
		)),
	)
	// ValidateJsonResponse checks if an HTTP response is a valid JSON response
	ValidateJsonResponse = F.Flow2(
		E.Of[error, *H.Response],
		E.ChainFirst(F.Flow5(
			GetHeader,
			R.Lookup[H.Header](HeaderContentType),
			O.Chain(A.First[string]),
			E.FromOption[error, string](errors.OnNone("unable to access the [%s] header", HeaderContentType)),
			E.ChainFirst(validateJsonContentTypeString),
		)))
)

const (
	HeaderContentType = "Content-Type"
)

// ParseMediaType parses a media type into a tuple
func ParseMediaType(mediaType string) E.Either[error, ParsedMediaType] {
	return E.TryCatchError(func() (ParsedMediaType, error) {
		m, p, err := mime.ParseMediaType(mediaType)
		return T.MakeTuple2(m, p), err
	})
}

func GetHeader(resp *H.Response) H.Header {
	return resp.Header
}

func GetBody(resp *H.Response) io.ReadCloser {
	return resp.Body
}

func isValidStatus(resp *H.Response) bool {
	return resp.StatusCode >= H.StatusOK && resp.StatusCode < H.StatusMultipleChoices
}

func StatusCodeError(resp *H.Response) error {
	return fmt.Errorf("invalid status code [%d] when accessing URL [%s]", resp.StatusCode, resp.Request.URL)
}