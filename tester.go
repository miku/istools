package islint

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/miku/islint/assetutil"
	"github.com/miku/span"
	"github.com/miku/span/finc"
)

//go:generate stringer -type=Kind
type Kind uint16

const (
	KeyTooLong Kind = iota
	InvalidStartPage
	InvalidEndPage
	EndPageBeforeStartPage
	InvalidURL
	SuspiciousPageCount
	PublicationDateTooEarly
	PublicationDateTooLate
	InvalidCollection
	RepeatedSubtitle
	CurrencyInTitle
	ExcessivePunctuation
	NoPublisher
	ShortAuthorName
	EtAlAuthorName
	NAInAuthorName
	WhitespaceAuthor
	RepeatedSlash
	NoURL
	NonCanonicalISSN
	HTMLEntityInAuthorName
)

var (
	// EarliestDate is the earliest publication date we accept.
	EarliestDate = time.Date(1458, 1, 1, 0, 0, 0, 0, time.UTC)
	// LatestDate represents the latest publication date we accept.
	LatestDate = time.Now().AddDate(5, 0, 0)

	// AllowedCollections
	AllowedCollections = assetutil.MustLoadStringSet("assets/collections/collections.tsv",
		"assets/collections/crossref.tsv")

	// currencyPattern is a rather narrow pattern:
	// http://rubular.com/r/WjcnjhckZq, used by NoCurrencyInTitle
	currencyPattern = regexp.MustCompile(`[€$¥][+-]?[0-9]{1,3}(?:[0-9]*(?:[.,][0-9]{2})?|(?:,[0-9]{3})*(?:\.[0-9]{2})?|(?:\.[0-9]{3})*(?:,[0-9]{2})?)`)
	// suspiciousPatterns, used by NoExcessivePunctuation
	suspiciousPatterns = []string{"?????", "!!!!!", "....."}
	// issnPattern looks for a canonical ISSN
	issnPattern = regexp.MustCompile(`^[0-9]{4,4}-[0-9]{3,3}[0-9X]$`)
	// htmlEntityPattern looks for leftover entities: http://rubular.com/r/flzmBzpShX
	htmlEntityPattern = regexp.MustCompile(`&(?:[a-z\d]+|#\d+|#x[a-f\d]+);`)
)

// Issue contains information about a quality issue in an intermediate schema
// record.
type Issue struct {
	Kind    Kind
	Record  finc.IntermediateSchema
	Message string
}

// Error formats the error.
func (e Issue) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.Record.RecordID, e.Kind, e.Message)
}

// TSV returns a tab representation.
func (e Issue) TSV() string {
	return fmt.Sprintf("%s\t%s\t%s", e.Record.RecordID, e.Kind, e.Message)
}

// TestSuite is a group of tests.
type TestSuite []Tester

// Tester is a intermediate record checker.
type Tester interface {
	TestRecord(finc.IntermediateSchema) error
}

// TesterFunc makes a function satisfy an interface.
type TesterFunc func(finc.IntermediateSchema) error

// TestRecord delegates test to the given func.
func (f TesterFunc) TestRecord(is finc.IntermediateSchema) error {
	return f(is)
}

var DefaultTests = []Tester{
	TesterFunc(KeyLength),
	TesterFunc(PlausiblePageCount),
	TesterFunc(ValidURL),
	TesterFunc(PlausibleDate),
	TesterFunc(AllowedCollectionNames),
	TesterFunc(SubtitleRepetition),
	TesterFunc(NoCurrencyInTitle),
	TesterFunc(NoExcessivePunctuation),
	TesterFunc(HasPublisher),
	TesterFunc(FeasibleAuthor),
	TesterFunc(NoRepeatedSlash),
	TesterFunc(HasURL),
	TesterFunc(CanonicalISSN),
}

// KeyLength checks the length of the record id. memcachedb limits this to 250
// bytes.
func KeyLength(is finc.IntermediateSchema) error {
	if len(is.RecordID) > span.KeyLengthLimit {
		return Issue{Kind: KeyTooLong, Record: is}
	}
	return nil
}

// ValidURL checks, if a URL string is parseable.
func ValidURL(is finc.IntermediateSchema) error {
	for _, s := range is.URL {
		if _, err := url.Parse(s); err != nil {
			return Issue{Kind: InvalidURL, Record: is, Message: s}
		}
	}
	return nil
}

// PlausibleDate checks for suspicious dates, refs. #5686.
func PlausibleDate(is finc.IntermediateSchema) error {
	if is.Date.Before(EarliestDate) {
		return Issue{Kind: PublicationDateTooEarly, Record: is, Message: is.Date.String()}
	}
	if is.Date.After(LatestDate) {
		return Issue{Kind: PublicationDateTooLate, Record: is, Message: is.Date.String()}
	}
	return nil
}

// PlausiblePageCount checks, wether the start and end page look plausible.
func PlausiblePageCount(is finc.IntermediateSchema) error {
	const (
		maxPageDigits = 6
		maxPageCount  = 20000
	)
	if len(is.StartPage) > maxPageDigits {
		return Issue{Kind: InvalidStartPage, Record: is, Message: is.StartPage}
	}
	if len(is.EndPage) > maxPageDigits {
		return Issue{Kind: InvalidEndPage, Record: is, Message: is.EndPage}
	}
	if is.StartPage != "" && is.EndPage != "" {
		if s, err := strconv.Atoi(is.StartPage); err == nil {
			if e, err := strconv.Atoi(is.EndPage); err == nil {
				if e < s {
					return Issue{Kind: EndPageBeforeStartPage, Record: is, Message: fmt.Sprintf("%v-%v", s, e)}
				}
				if e-s > maxPageCount {
					return Issue{Kind: SuspiciousPageCount, Record: is, Message: fmt.Sprintf("%v-%v", s, e)}
				}
			} else {
				return Issue{Kind: InvalidEndPage, Record: is, Message: is.EndPage}
			}
		} else {
			return Issue{Kind: InvalidStartPage, Record: is, Message: is.StartPage}
		}
	}
	return nil
}

// AllowedCollectionNames checks for a fixed list of allowed collection names,
// stored under assets, refs. #6496.
func AllowedCollectionNames(is finc.IntermediateSchema) error {
	if !AllowedCollections.Contains(is.MegaCollection) {
		return Issue{Kind: InvalidCollection, Record: is, Message: is.MegaCollection}
	}
	return nil
}

// SubtitleRepetition, refs #6553.
func SubtitleRepetition(is finc.IntermediateSchema) error {
	if is.ArticleSubtitle != "" && strings.Contains(is.ArticleTitle, is.ArticleSubtitle) {
		return Issue{Kind: RepeatedSubtitle, Record: is,
			Message: fmt.Sprintf("TITLE: %s, SUBTITLE: %s", is.ArticleTitle, is.ArticleSubtitle)}
	}
	return nil
}

// NoCurrencyInTitle, e.g. http://goo.gl/HACBcW
// Cartier , Marie . Baby, You Are My Religion: Women, Gay Bars, and Theology
// Before Stonewall . Gender, Theology and Spirituality. Durham, UK: Acumen,
// 2013. xii+256 pp. $90.00 (cloth); $29.95 (paper).
func NoCurrencyInTitle(is finc.IntermediateSchema) error {
	if currencyPattern.MatchString(is.ArticleTitle) {
		return Issue{Kind: CurrencyInTitle, Record: is, Message: is.ArticleTitle}
	}
	return nil
}

// NoExcessivePuctuation should detect things like this title:
// CrossRef????????????? https://goo.gl/AD0V1o
func NoExcessivePunctuation(is finc.IntermediateSchema) error {
	for _, p := range suspiciousPatterns {
		if strings.Contains(is.ArticleTitle, p) {
			return Issue{Kind: ExcessivePunctuation, Record: is, Message: is.ArticleTitle}
		}
	}
	return nil
}

// HasPublisher tests, whether a publisher is given.
func HasPublisher(is finc.IntermediateSchema) error {
	switch len(is.Publishers) {
	case 0:
		return Issue{Kind: NoPublisher, Record: is}
	case 1:
		if is.Publishers[0] == "" {
			return Issue{Kind: NoPublisher, Record: is}
		}
	default:
		for _, p := range is.Publishers {
			if p == "" {
				return Issue{Kind: NoPublisher, Record: is, Message: "empty string as publisher name"}
			}
		}
	}
	return nil
}

// FeasibleAuthor checks for a few suspicious authors patterns, refs. #4892, #4940, #5895.
func FeasibleAuthor(is finc.IntermediateSchema) error {
	for _, author := range is.Authors {
		s := author.String()
		if len(s) < 5 {
			return Issue{Kind: ShortAuthorName, Record: is, Message: s}
		}
		lower := strings.ToLower(s)
		if strings.HasPrefix(lower, "et al") {
			return Issue{Kind: EtAlAuthorName, Record: is, Message: s}
		}
		if strings.Contains(lower, "&na;") {
			return Issue{Kind: NAInAuthorName, Record: is, Message: s}
		}
		if len(s) > 0 && strings.TrimSpace(s) == "" {
			return Issue{Kind: WhitespaceAuthor, Record: is, Message: "author contains whitespace only"}
		}
		if htmlEntityPattern.MatchString(s) {
			return Issue{Kind: HTMLEntityInAuthorName, Record: is, Message: s}
		}
	}
	return nil
}

// NoRepeatedSlash checks a DOI for repeated slashes, refs. #6312.
func NoRepeatedSlash(is finc.IntermediateSchema) error {
	if strings.Contains(is.DOI, "//") {
		return Issue{Kind: RepeatedSlash, Record: is, Message: is.DOI}
	}
	return nil
}

// HasURL checks for a value in URL. This is no URL validation.
func HasURL(is finc.IntermediateSchema) error {
	if len(is.URL) == 0 {
		return Issue{Kind: NoURL, Record: is}
	}
	return nil
}

// CanonicalISSN checks for the canonical ISSN format 1234-567X.
func CanonicalISSN(is finc.IntermediateSchema) error {
	for _, issn := range append(is.ISSN, is.EISSN...) {
		if !issnPattern.MatchString(issn) {
			return Issue{Kind: NonCanonicalISSN, Record: is, Message: issn}
		}
	}
	return nil
}
