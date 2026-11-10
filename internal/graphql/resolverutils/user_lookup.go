package resolverutils

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	fperrors "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	T "github.com/IBM/fp-go/v2/tuple"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	"github.com/dictyBase/fp-go-loom/strutils"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	"github.com/dictyBase/graphql-server/internal/authentication"
	"github.com/dictyBase/graphql-server/internal/graphql/errorutils"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	fakeUserID                int64 = 2375
	depositorFallbackID       int64 = 0
	userResourceType                = "user"
	fakeUserFirstName               = "Dicty"
	fakeUserLastName                = "Stock Center"
	fakeUserEmail                   = "dictystocks@northwestern.edu"
	fakeUserOrg                     = "Northwestern University"
	defaultDepositorFirstName       = "External"
	defaultDepositorLastName        = "Depositor"
)

// UserLookupState carries the request context, user identifier or email,
// Logto client, logger, and intermediate lookup results.
type UserLookupState struct {
	ctx    context.Context
	query  string
	client authentication.LogtoClient
	logger *logrus.Entry
	resp   *authentication.UserResp
	userID int64
}

var (
	numRgx   = regexp.MustCompile("[0-9]+")
	alphaRgx = regexp.MustCompile("[A-Za-z]+")

	nilUser *user.User

	// LookupCtxLens accesses the request context on UserLookupState.
	LookupCtxLens = L.MakeLens(
		func(s UserLookupState) context.Context { return s.ctx },
		func(s UserLookupState, v context.Context) UserLookupState {
			s.ctx = v
			return s
		},
	)
	// LookupQueryLens accesses the user email or ID query on UserLookupState.
	LookupQueryLens = L.MakeLens(
		func(s UserLookupState) string { return s.query },
		func(s UserLookupState, v string) UserLookupState {
			s.query = v
			return s
		},
	)
	// LookupClientLens accesses the LogtoClient on UserLookupState.
	LookupClientLens = L.MakeLens(
		func(s UserLookupState) authentication.LogtoClient { return s.client },
		func(
			s UserLookupState,
			v authentication.LogtoClient,
		) UserLookupState {
			s.client = v
			return s
		},
	)
	// LookupLoggerLens accesses the logger on UserLookupState.
	LookupLoggerLens = L.MakeLens(
		func(s UserLookupState) *logrus.Entry { return s.logger },
		func(s UserLookupState, v *logrus.Entry) UserLookupState {
			s.logger = v
			return s
		},
	)
	lookupRespLens = L.MakeLens(
		func(s UserLookupState) *authentication.UserResp { return s.resp },
		func(
			s UserLookupState,
			v *authentication.UserResp,
		) UserLookupState {
			s.resp = v
			return s
		},
	)
	lookupUserIDLens = L.MakeLens(
		func(s UserLookupState) int64 { return s.userID },
		func(s UserLookupState, v int64) UserLookupState {
			s.userID = v
			return s
		},
	)
	userRespIDLens = L.MakeLensRef(
		func(r *authentication.UserResp) string { return r.ID },
		func(r *authentication.UserResp, v string) *authentication.UserResp {
			r.ID = v
			return r
		},
	)
	lookupRespIDLens = F.Pipe1(
		lookupRespLens,
		L.Compose[UserLookupState](userRespIDLens),
	)

	userDataLens = L.MakeLensRef(
		func(u *user.User) *user.UserData { return u.Data },
		func(u *user.User, v *user.UserData) *user.User {
			u.Data = v
			return u
		},
	)
	userDataTypeLens = L.MakeLensRef(
		func(d *user.UserData) string { return d.Type },
		func(d *user.UserData, v string) *user.UserData {
			d.Type = v
			return d
		},
	)
	userDataIDLens = L.MakeLensRef(
		func(d *user.UserData) int64 { return d.Id },
		func(d *user.UserData, v int64) *user.UserData {
			d.Id = v
			return d
		},
	)
	userDataAttrsLens = L.MakeLensRef(
		func(d *user.UserData) *user.UserAttributes { return d.Attributes },
		func(d *user.UserData, v *user.UserAttributes) *user.UserData {
			d.Attributes = v
			return d
		},
	)
	userTypeLens = F.Pipe1(
		userDataLens,
		L.Compose[*user.User](userDataTypeLens),
	)
	userIDLens = F.Pipe1(
		userDataLens,
		L.Compose[*user.User](userDataIDLens),
	)
	userAttrsLens = F.Pipe1(
		userDataLens,
		L.Compose[*user.User](userDataAttrsLens),
	)

	userAttrFirstNameLens = L.MakeLensRef(
		func(a *user.UserAttributes) string { return a.FirstName },
		func(a *user.UserAttributes, v string) *user.UserAttributes {
			a.FirstName = v
			return a
		},
	)
	userAttrLastNameLens = L.MakeLensRef(
		func(a *user.UserAttributes) string { return a.LastName },
		func(a *user.UserAttributes, v string) *user.UserAttributes {
			a.LastName = v
			return a
		},
	)
	userAttrEmailLens = L.MakeLensRef(
		func(a *user.UserAttributes) string { return a.Email },
		func(a *user.UserAttributes, v string) *user.UserAttributes {
			a.Email = v
			return a
		},
	)
	userAttrOrgLens = L.MakeLensRef(
		func(a *user.UserAttributes) string { return a.Organization },
		func(a *user.UserAttributes, v string) *user.UserAttributes {
			a.Organization = v
			return a
		},
	)
	userAttrFirstAddrLens = L.MakeLensRef(
		func(a *user.UserAttributes) string { return a.FirstAddress },
		func(a *user.UserAttributes, v string) *user.UserAttributes {
			a.FirstAddress = v
			return a
		},
	)
	userAttrCityLens = L.MakeLensRef(
		func(a *user.UserAttributes) string { return a.City },
		func(a *user.UserAttributes, v string) *user.UserAttributes {
			a.City = v
			return a
		},
	)
	userAttrStateLens = L.MakeLensRef(
		func(a *user.UserAttributes) string { return a.State },
		func(a *user.UserAttributes, v string) *user.UserAttributes {
			a.State = v
			return a
		},
	)
	userAttrZipLens = L.MakeLensRef(
		func(a *user.UserAttributes) string { return a.Zipcode },
		func(a *user.UserAttributes, v string) *user.UserAttributes {
			a.Zipcode = v
			return a
		},
	)
	userAttrCountryLens = L.MakeLensRef(
		func(a *user.UserAttributes) string { return a.Country },
		func(a *user.UserAttributes, v string) *user.UserAttributes {
			a.Country = v
			return a
		},
	)
	userAttrPhoneLens = L.MakeLensRef(
		func(a *user.UserAttributes) string { return a.Phone },
		func(a *user.UserAttributes, v string) *user.UserAttributes {
			a.Phone = v
			return a
		},
	)
	userAttrActiveLens = L.MakeLensRef(
		(*user.UserAttributes).GetIsActive,
		func(a *user.UserAttributes, v bool) *user.UserAttributes {
			a.IsActive = v
			return a
		},
	)

	userFirstNameLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrFirstNameLens),
	)
	userLastNameLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrLastNameLens),
	)
	userEmailLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrEmailLens),
	)
	userOrgLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrOrgLens),
	)
	userFirstAddrLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrFirstAddrLens),
	)
	userCityLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrCityLens),
	)
	userStateLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrStateLens),
	)
	userZipLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrZipLens),
	)
	userCountryLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrCountryLens),
	)
	userPhoneLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrPhoneLens),
	)
	userActiveLens = F.Pipe1(
		userAttrsLens,
		L.Compose[*user.User](userAttrActiveLens),
	)

	isAuthDisabledErr = F.Bind2nd(errors.Is, authentication.ErrAuthDisabled)
	isQueryEmpty      = F.Pipe1(
		S.IsEmpty,
		P.ContraMap(LookupQueryLens.Get),
	)

	recoverWithFakeUser = F.Pipe1(
		isAuthDisabledErr,
		P.Fold(
			IOE.Left[*user.User, error],
			F.Constant1[error](IOE.Of[error](FakeUser())),
		),
	)

	splitAtSign      = F.Bind2nd(strings.Split, "@")
	findAlphaTokens  = F.Bind2nd(alphaRgx.FindAllString, -1)
	extractEmailHead = F.Flow3(
		splitAtSign,
		A.Head[string],
		O.GetOrElse(F.Constant("")),
	)
	titleCaser      = cases.Title(language.Und)
	capitalizeToken = F.Flow2(
		S.ToLowerCase,
		titleCaser.String,
	)
	depositorTokensFromEmail = F.Flow3(
		extractEmailHead,
		findAlphaTokens,
		A.Map(capitalizeToken),
	)

	wrapEmailLookupErr = fperrors.OnError("unable to retreieve user")
	wrapNumericIDparse = fperrors.OnError(
		"unable to convert user id to integer",
	)
	wrapDirectIDParse = fperrors.OnError(
		"error in converting user id to int64",
	)

	parseBase10Int64 = E.Eitherize1(func(raw string) (int64, error) {
		return strconv.ParseInt(raw, 10, 64)
	})
	findDigits = F.Bind2nd(numRgx.FindAllString, -1)

	withNumericUserID = F.Flow2(
		E.Of[error, UserLookupState],
		E.Bind(lookupUserIDLens.Set, extractNumericUserID),
	)
	withDirectRespUserID = F.Flow2(
		E.Of[error, UserLookupState],
		E.Bind(
			lookupUserIDLens.Set,
			F.Flow3(
				lookupRespIDLens.Get,
				parseBase10Int64,
				E.MapLeft[int64](wrapDirectIDParse),
			),
		),
	)
	withDirectQueryUserID = F.Flow2(
		E.Of[error, UserLookupState],
		E.Bind(
			lookupUserIDLens.Set,
			F.Flow3(
				LookupQueryLens.Get,
				parseBase10Int64,
				E.MapLeft[int64](wrapDirectIDParse),
			),
		),
	)

	checkReadUserStep = F.Pipe1(
		isStateNoopAuthClient,
		P.Fold(
			runCheckReadUser,
			F.Constant1[UserLookupState](
				IOE.Left[UserLookupState](authentication.ErrAuthDisabled),
			),
		),
	)

	branchStockDepositor = F.Pipe1(
		isQueryEmpty,
		P.Fold(
			resolveDepositorWhenNonEmpty,
			F.Constant1[UserLookupState](IOE.Of[error](FakeUser())),
		),
	)
)

func emptyUserEntity() *user.User {
	return F.Pipe3(
		&user.User{},
		userDataLens.Set(&user.UserData{}),
		userAttrsLens.Set(&user.UserAttributes{}),
		userTypeLens.Set(userResourceType),
	)
}

// FakeUser returns the default Dicty Stock Center user representation.
func FakeUser() *user.User {
	return F.Pipe6(
		emptyUserEntity(),
		userIDLens.Set(fakeUserID),
		userFirstNameLens.Set(fakeUserFirstName),
		userLastNameLens.Set(fakeUserLastName),
		userEmailLens.Set(fakeUserEmail),
		userOrgLens.Set(fakeUserOrg),
		userActiveLens.Set(true),
	)
}

func buildRedactedDepositorFromTokens(tokens []string) *user.User {
	firstName := F.Pipe2(
		tokens,
		A.Head[string],
		O.GetOrElse(F.Constant(defaultDepositorFirstName)),
	)
	lastName := F.Pipe3(
		tokens,
		A.Tail[string],
		O.Map(strutils.JoinSpace),
		O.GetOrElse(F.Constant("")),
	)
	return F.Pipe4(
		emptyUserEntity(),
		userIDLens.Set(depositorFallbackID),
		userFirstNameLens.Set(firstName),
		userLastNameLens.Set(lastName),
		userActiveLens.Set(true),
	)
}

func defaultRedactedDepositor() *user.User {
	return F.Pipe4(
		emptyUserEntity(),
		userIDLens.Set(depositorFallbackID),
		userFirstNameLens.Set(defaultDepositorFirstName),
		userLastNameLens.Set(defaultDepositorLastName),
		userActiveLens.Set(true),
	)
}

// RedactedDepositorUser constructs a User from a depositor email without
// exposing the email address or returning the Dicty Stock Center fakeUser.
func RedactedDepositorUser(email string) *user.User {
	return F.Pipe3(
		email,
		depositorTokensFromEmail,
		O.FromPredicate(predarrays.IsNonEmpty[string]()),
		O.Fold(
			defaultRedactedDepositor,
			buildRedactedDepositorFromTokens,
		),
	)
}

func recoverWithRedactedDepositor(
	s UserLookupState,
	err error,
) IOE.IOEither[error, *user.User] {
	redacted := F.Pipe2(
		s,
		LookupQueryLens.Get,
		RedactedDepositorUser,
	)
	return F.Pipe2(
		err,
		O.FromPredicate(isAuthDisabledErr),
		O.Fold(
			func() IOE.IOEither[error, *user.User] {
				return IOE.Left[*user.User](err)
			},
			func(error) IOE.IOEither[error, *user.User] {
				return IOE.Of[error](redacted)
			},
		),
	)
}

func joinDigitParts(parts []string) string {
	return strutils.JoinStrings(parts...)
}

func extractNumericUserID(s UserLookupState) E.Either[error, int64] {
	rawID := lookupRespIDLens.Get(s)
	return F.Pipe6(
		rawID,
		findDigits,
		O.FromPredicate(predarrays.IsNonEmpty[string]()),
		E.FromOption[[]string](func() error {
			return fmt.Errorf("cannot convert user id to number %s", rawID)
		}),
		E.Map[error](joinDigitParts),
		E.Chain(parseBase10Int64),
		E.MapLeft[int64](wrapNumericIDparse),
	)
}

func fetchUserRespByEmail(
	s UserLookupState,
) IOE.IOEither[error, *authentication.UserResp] {
	client := LookupClientLens.Get(s)
	email := LookupQueryLens.Get(s)
	return F.Pipe1(
		IOE.TryCatchError(func() (*authentication.UserResp, error) {
			return client.UserWithEmail(email)
		}),
		IOE.MapLeft[*authentication.UserResp](wrapEmailLookupErr),
	)
}

func fetchUserRespByEmailRaw(
	s UserLookupState,
) IOE.IOEither[error, *authentication.UserResp] {
	client := LookupClientLens.Get(s)
	email := LookupQueryLens.Get(s)
	return IOE.TryCatchError(func() (*authentication.UserResp, error) {
		return client.UserWithEmail(email)
	})
}

func fetchUserRespByID(
	s UserLookupState,
) IOE.IOEither[error, *authentication.UserResp] {
	client := LookupClientLens.Get(s)
	id := LookupQueryLens.Get(s)
	return IOE.TryCatchError(func() (*authentication.UserResp, error) {
		return client.User(id)
	})
}

func isStateNoopAuthClient(s UserLookupState) bool {
	_, ok := LookupClientLens.Get(s).(authentication.NoopClient)
	return ok
}

func runCheckReadUser(s UserLookupState) IOE.IOEither[error, UserLookupState] {
	ctx := LookupCtxLens.Get(s)
	return IOE.TryCatchError(func() (UserLookupState, error) {
		return s, authentication.CheckReadUser(ctx)
	})
}

func buildUserFromState(s UserLookupState) *user.User {
	resp := lookupRespLens.Get(s)
	return F.Pipe12(
		emptyUserEntity(),
		userIDLens.Set(lookupUserIDLens.Get(s)),
		userFirstNameLens.Set(resp.Username),
		userLastNameLens.Set(resp.Name),
		userEmailLens.Set(resp.PrimaryEmail),
		userOrgLens.Set(resp.CustomData.Institution),
		userFirstAddrLens.Set(resp.CustomData.Address),
		userCityLens.Set(resp.CustomData.City),
		userStateLens.Set(resp.CustomData.State),
		userZipLens.Set(resp.CustomData.Zipcode),
		userCountryLens.Set(resp.CustomData.Country),
		userPhoneLens.Set(resp.PrimaryPhone),
		userActiveLens.Set(true),
	)
}

func recordUserLookupError(
	s UserLookupState,
	err error,
) IOE.IOEither[error, IO.Void] {
	ctx := LookupCtxLens.Get(s)
	logger := LookupLoggerLens.Get(s)
	return F.Pipe1(
		IO.FromImpure(func() {
			errorutils.AddGQLError(ctx, err)
			logger.Error(err)
		}),
		IOE.FromIO[error, IO.Void],
	)
}

func resolveDepositorWhenNonEmpty(
	s UserLookupState,
) IOE.IOEither[error, *user.User] {
	return F.Pipe6(
		s,
		IOE.Of[error, UserLookupState],
		IOE.Bind(lookupRespLens.Set, fetchUserRespByEmail),
		IOE.ChainEitherK(withNumericUserID),
		IOE.Map[error](buildUserFromState),
		IOE.OrElse(F.Bind1st(recoverWithRedactedDepositor, s)),
		IOE.TapLeft[*user.User](F.Bind1st(recordUserLookupError, s)),
	)
}

func onUserLookupError(err error) T.Tuple2[error, *user.User] {
	return T.MakeTuple2(err, nilUser)
}

func onUserLookupSuccess(u *user.User) T.Tuple2[error, *user.User] {
	return T.MakeTuple2[error](nil, u)
}

// UserByEmail resolves a User by email and returns FakeUser when
// authentication is disabled.
func UserByEmail(s UserLookupState) T.Tuple2[error, *user.User] {
	return F.Pipe8(
		s,
		IOE.Of[error, UserLookupState],
		IOE.Bind(lookupRespLens.Set, fetchUserRespByEmail),
		IOE.ChainEitherK(withNumericUserID),
		IOE.Map[error](buildUserFromState),
		IOE.OrElse(recoverWithFakeUser),
		IOE.TapLeft[*user.User](F.Bind1st(recordUserLookupError, s)),
		ioeutils.ToEither[error, *user.User],
		E.Fold(onUserLookupError, onUserLookupSuccess),
	)
}

// StockDepositor resolves a stock depositor. An empty depositor returns
// FakeUser. When authentication is disabled and the depositor email is
// non-empty, it returns a redacted depositor user without exposing the email.
func StockDepositor(s UserLookupState) T.Tuple2[error, *user.User] {
	return F.Pipe3(
		s,
		branchStockDepositor,
		ioeutils.ToEither[error, *user.User],
		E.Fold(onUserLookupError, onUserLookupSuccess),
	)
}

// UserByID resolves Query.user(id) and returns FakeUser when authentication is
// disabled.
func UserByID(s UserLookupState) T.Tuple2[error, *user.User] {
	return F.Pipe9(
		s,
		IOE.Of[error, UserLookupState],
		IOE.Chain(checkReadUserStep),
		IOE.Bind(lookupRespLens.Set, fetchUserRespByID),
		IOE.ChainEitherK(withDirectQueryUserID),
		IOE.Map[error](buildUserFromState),
		IOE.OrElse(recoverWithFakeUser),
		IOE.TapLeft[*user.User](F.Bind1st(recordUserLookupError, s)),
		ioeutils.ToEither[error, *user.User],
		E.Fold(onUserLookupError, onUserLookupSuccess),
	)
}

// QueryUserByEmail resolves Query.userByEmail(email) and returns FakeUser when
// authentication is disabled.
func QueryUserByEmail(s UserLookupState) T.Tuple2[error, *user.User] {
	return F.Pipe8(
		s,
		IOE.Of[error, UserLookupState],
		IOE.Bind(lookupRespLens.Set, fetchUserRespByEmailRaw),
		IOE.ChainEitherK(withDirectRespUserID),
		IOE.Map[error](buildUserFromState),
		IOE.OrElse(recoverWithFakeUser),
		IOE.TapLeft[*user.User](F.Bind1st(recordUserLookupError, s)),
		ioeutils.ToEither[error, *user.User],
		E.Fold(onUserLookupError, onUserLookupSuccess),
	)
}
