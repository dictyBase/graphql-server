package content

import (
	"regexp"
	"strings"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"
	"github.com/dictyBase/fp-go-loom/strutils"
)

// storageOriginPattern matches the scheme and authority of an absolute URL up
// to the S3 bucket segment that follows it. The character class excludes `/`
// so an external URL that happens to contain the bucket prefix deeper in its
// path is never rewritten.
const storageOriginPattern = `https?://[^/"'\s)\\]+`

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
	targetPrefix  string
	pattern       *regexp.Regexp
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
	targetPrefixLens = L.MakeLens(
		func(s contentStorageState) string {
			return s.targetPrefix
		},
		func(s contentStorageState, value string) contentStorageState {
			s.targetPrefix = value
			return s
		},
	)
	patternLens = L.MakeLens(
		func(s contentStorageState) *regexp.Regexp {
			return s.pattern
		},
		func(s contentStorageState, value *regexp.Regexp) contentStorageState {
			s.pattern = value
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

func buildTargetPrefix(s contentStorageState) string {
	return strutils.JoinStrings(s.storageURL, s.segmentPrefix)
}

func compileStorageURLPattern(s contentStorageState) O.Option[*regexp.Regexp] {
	return F.Pipe4(
		s.segmentPrefix,
		regexp.QuoteMeta,
		S.Prepend(storageOriginPattern),
		E.Eitherize1(regexp.Compile),
		E.ToOption[error, *regexp.Regexp],
	)
}

func buildPrefixReplacer(s contentStorageState) func(string) string {
	isTargetPrefix := predord.StrEq(s.targetPrefix)
	return F.Flow2(
		O.FromPredicate(isTargetPrefix),
		O.GetOrElse(F.Constant(s.targetPrefix)),
	)
}

func rewriteStorageURLsInState(s contentStorageState) string {
	return s.pattern.ReplaceAllStringFunc(s.rawContent, s.replacer)
}
