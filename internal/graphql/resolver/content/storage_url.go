package content

import (
	"regexp"
	"strings"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/dictyBase/fp-go-loom/strutils"
)

// absoluteURLPattern matches one complete absolute http(s) URL. The character
// class stops at the characters that end a URL inside the stored Lexical JSON,
// so the whole URL including its path arrives in the replacer.
var absoluteURLPattern = regexp.MustCompile(`https?://[^"'\s)\\]+`)

const urlSchemeSeparator = "://"

// contentStorageState carries the caller-supplied S3 configuration and raw
// Lexical JSON content through the URL normalization pipeline. All derived
// fields are attached by lenses so every pipeline step receives a single
// state parameter.
type contentStorageState struct {
	storageURL    string
	bucket        string
	bucketPath    string
	rawContent    string
	segmentPrefix string
	replacer      func(string) string
}

var (
	storageURLLens = L.MakeLens(
		func(s contentStorageState) string {
			return s.storageURL
		},
		func(s contentStorageState, value string) contentStorageState {
			s.storageURL = value
			return s
		},
	)
	bucketLens = L.MakeLens(
		func(s contentStorageState) string {
			return s.bucket
		},
		func(s contentStorageState, value string) contentStorageState {
			s.bucket = value
			return s
		},
	)
	bucketPathLens = L.MakeLens(
		func(s contentStorageState) string {
			return s.bucketPath
		},
		func(s contentStorageState, value string) contentStorageState {
			s.bucketPath = value
			return s
		},
	)
	segmentPrefixLens = L.MakeLens(
		func(s contentStorageState) string {
			return s.segmentPrefix
		},
		func(s contentStorageState, value string) contentStorageState {
			s.segmentPrefix = value
			return s
		},
	)
	replacerLens = L.MakeLens(
		func(s contentStorageState) func(string) string {
			return s.replacer
		},
		func(
			s contentStorageState,
			value func(string) string,
		) contentStorageState {
			s.replacer = value
			return s
		},
	)

	trimTrailingSlash = F.Bind2nd(strings.TrimRight, "/")
	trimSlashes       = F.Bind2nd(strings.Trim, "/")

	hasStorageURL = F.Pipe1(
		S.IsNonEmpty,
		P.ContraMap(storageURLLens.Get),
	)
	hasBucket = F.Pipe1(
		S.IsNonEmpty,
		P.ContraMap(bucketLens.Get),
	)
	hasBucketPath = F.Pipe1(
		S.IsNonEmpty,
		P.ContraMap(bucketPathLens.Get),
	)
	isStorageConfigComplete = F.Pipe2(
		hasStorageURL,
		P.And(hasBucket),
		P.And(hasBucketPath),
	)
)

func buildSegmentPrefix(s contentStorageState) string {
	return strutils.JoinStrings("/", s.bucket, "/", s.bucketPath, "/")
}

// pathOfURL returns the path of an absolute URL. The search runs on a copy
// with an appended "/" so a URL without a path yields an empty path instead of
// an out-of-range slice.
func pathOfURL(raw string) string {
	rest := raw[strings.Index(raw, urlSchemeSeparator)+len(urlSchemeSeparator):]
	separatorIndex := F.Bind2nd(strings.Index, "/")
	pathStart := F.Pipe1(rest, F.Flow2(S.Append("/"), separatorIndex))
	return rest[pathStart:]
}

// buildURLReplacer rewrites the authority only when the path starts with
// /<bucket>/<bucketPath>/, and keeps every other URL byte-identical.
func buildURLReplacer(s contentStorageState) func(string) string {
	hasSegmentPath := F.Bind2nd(strings.HasPrefix, s.segmentPrefix)
	toStorageURL := F.Flow2(
		O.FromPredicate(hasSegmentPath),
		O.Map(S.Prepend(s.storageURL)),
	)
	return func(raw string) string {
		return F.Pipe2(
			pathOfURL(raw),
			toStorageURL,
			O.GetOrElse(F.Constant(raw)),
		)
	}
}

func rewriteStorageURLsInState(s contentStorageState) string {
	return absoluteURLPattern.ReplaceAllStringFunc(s.rawContent, s.replacer)
}
