// generated from grammar.peg
package dslparser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"
	"golang.org/x/text/unicode/norm"
)

var transcriptionMap = map[rune]string{
    'Ћ': "Ө",
    'Ђ': "i:",
    '\'': "ˈ",
    'э': "ɪ",
    'Ѓ': "ɑ:",
    '†': "ə",
    '‡': "æ",
    '…': "ʌ",
    'Њ': "ŋ",
    'ю': "ɜ:",
    '‚': "ɔ",
    'ѓ': "u:",
    '¶': "ʊ",
    '€': "ɔ",
    '‹': "dʒ",
    'Џ': "ʃ",
    ' ': "tʃ", // ʧ
    '‰': "ð",
    'Љ': "ʒ",
}

var g = &grammar {
	rules: []*rule{
{
	name: "All",
	pos: position{line: 40, col: 1, offset: 566},
	expr: &actionExpr{
	pos: position{line: 40, col: 7, offset: 574},
	run: (*parser).callonAll1,
	expr: &seqExpr{
	pos: position{line: 40, col: 7, offset: 574},
	exprs: []interface{}{
&labeledExpr{
	pos: position{line: 40, col: 7, offset: 574},
	label: "b",
	expr: &ruleRefExpr{
	pos: position{line: 40, col: 9, offset: 576},
	name: "Body",
},
},
&ruleRefExpr{
	pos: position{line: 40, col: 14, offset: 581},
	name: "EOF",
},
	},
},
},
},
{
	name: "Body",
	pos: position{line: 44, col: 1, offset: 617},
	expr: &actionExpr{
	pos: position{line: 44, col: 8, offset: 626},
	run: (*parser).callonBody1,
	expr: &labeledExpr{
	pos: position{line: 44, col: 8, offset: 626},
	label: "ps",
	expr: &oneOrMoreExpr{
	pos: position{line: 44, col: 11, offset: 629},
	expr: &ruleRefExpr{
	pos: position{line: 44, col: 11, offset: 629},
	name: "Part",
},
},
},
},
},
{
	name: "Part",
	pos: position{line: 52, col: 1, offset: 797},
	expr: &choiceExpr{
	pos: position{line: 52, col: 8, offset: 806},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 52, col: 8, offset: 806},
	name: "EscapeSequence",
},
&ruleRefExpr{
	pos: position{line: 52, col: 25, offset: 823},
	name: "Accent",
},
&ruleRefExpr{
	pos: position{line: 52, col: 34, offset: 832},
	name: "Transcription",
},
&ruleRefExpr{
	pos: position{line: 52, col: 50, offset: 848},
	name: "OpenTag",
},
&ruleRefExpr{
	pos: position{line: 52, col: 60, offset: 858},
	name: "CloseTag",
},
&ruleRefExpr{
	pos: position{line: 52, col: 71, offset: 869},
	name: "Newline",
},
&ruleRefExpr{
	pos: position{line: 52, col: 81, offset: 879},
	name: "Anything",
},
	},
},
},
{
	name: "EscapeSequence",
	pos: position{line: 54, col: 1, offset: 889},
	expr: &actionExpr{
	pos: position{line: 54, col: 18, offset: 908},
	run: (*parser).callonEscapeSequence1,
	expr: &seqExpr{
	pos: position{line: 54, col: 18, offset: 908},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 54, col: 18, offset: 908},
	val: "\\",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 54, col: 23, offset: 913},
	label: "ch",
	expr: &choiceExpr{
	pos: position{line: 54, col: 27, offset: 917},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 54, col: 27, offset: 917},
	val: "[",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 54, col: 33, offset: 923},
	val: "]",
	ignoreCase: false,
},
	},
},
},
	},
},
},
},
{
	name: "Accent",
	pos: position{line: 58, col: 1, offset: 969},
	expr: &actionExpr{
	pos: position{line: 58, col: 10, offset: 980},
	run: (*parser).callonAccent1,
	expr: &seqExpr{
	pos: position{line: 58, col: 10, offset: 980},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 58, col: 10, offset: 980},
	val: "[']",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 58, col: 16, offset: 986},
	label: "ch",
	expr: &anyMatcher{
	line: 58, col: 19, offset: 989,
},
},
&litMatcher{
	pos: position{line: 58, col: 21, offset: 991},
	val: "[/']",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Transcription",
	pos: position{line: 65, col: 1, offset: 1185},
	expr: &actionExpr{
	pos: position{line: 65, col: 17, offset: 1203},
	run: (*parser).callonTranscription1,
	expr: &seqExpr{
	pos: position{line: 65, col: 17, offset: 1203},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 65, col: 17, offset: 1203},
	val: "[t]",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 65, col: 23, offset: 1209},
	label: "ps",
	expr: &oneOrMoreExpr{
	pos: position{line: 65, col: 26, offset: 1212},
	expr: &choiceExpr{
	pos: position{line: 65, col: 28, offset: 1214},
	alternatives: []interface{}{
&ruleRefExpr{
	pos: position{line: 65, col: 28, offset: 1214},
	name: "EscapeSequence",
},
&charClassMatcher{
	pos: position{line: 65, col: 45, offset: 1231},
	val: "[^[]",
	chars: []rune{'[',},
	ignoreCase: false,
	inverted: true,
},
	},
},
},
},
&litMatcher{
	pos: position{line: 65, col: 53, offset: 1239},
	val: "[/t]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "OpenTag",
	pos: position{line: 90, col: 1, offset: 1968},
	expr: &actionExpr{
	pos: position{line: 90, col: 11, offset: 1980},
	run: (*parser).callonOpenTag1,
	expr: &seqExpr{
	pos: position{line: 90, col: 11, offset: 1980},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 90, col: 11, offset: 1980},
	val: "[",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 90, col: 15, offset: 1984},
	label: "tnitf",
	expr: &ruleRefExpr{
	pos: position{line: 90, col: 21, offset: 1990},
	name: "Tagname",
},
},
&labeledExpr{
	pos: position{line: 90, col: 29, offset: 1998},
	label: "targ",
	expr: &zeroOrOneExpr{
	pos: position{line: 90, col: 34, offset: 2003},
	expr: &ruleRefExpr{
	pos: position{line: 90, col: 34, offset: 2003},
	name: "Tagarg",
},
},
},
&litMatcher{
	pos: position{line: 90, col: 42, offset: 2011},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "CloseTag",
	pos: position{line: 129, col: 1, offset: 2860},
	expr: &actionExpr{
	pos: position{line: 129, col: 12, offset: 2873},
	run: (*parser).callonCloseTag1,
	expr: &seqExpr{
	pos: position{line: 129, col: 12, offset: 2873},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 129, col: 12, offset: 2873},
	val: "[/",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 129, col: 17, offset: 2878},
	label: "tnitf",
	expr: &ruleRefExpr{
	pos: position{line: 129, col: 23, offset: 2884},
	name: "Tagname",
},
},
&litMatcher{
	pos: position{line: 129, col: 31, offset: 2892},
	val: "]",
	ignoreCase: false,
},
	},
},
},
},
{
	name: "Tagname",
	pos: position{line: 161, col: 1, offset: 3494},
	expr: &actionExpr{
	pos: position{line: 161, col: 11, offset: 3506},
	run: (*parser).callonTagname1,
	expr: &labeledExpr{
	pos: position{line: 161, col: 11, offset: 3506},
	label: "n",
	expr: &choiceExpr{
	pos: position{line: 161, col: 14, offset: 3509},
	alternatives: []interface{}{
&litMatcher{
	pos: position{line: 161, col: 14, offset: 3509},
	val: "m1",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 21, offset: 3516},
	val: "m2",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 28, offset: 3523},
	val: "m3",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 35, offset: 3530},
	val: "m",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 41, offset: 3536},
	val: "b",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 47, offset: 3542},
	val: "com",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 55, offset: 3550},
	val: "c",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 61, offset: 3556},
	val: "'",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 67, offset: 3562},
	val: "sup",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 75, offset: 3570},
	val: "p",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 81, offset: 3576},
	val: "ex",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 88, offset: 3583},
	val: "lang",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 97, offset: 3592},
	val: "i",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 103, offset: 3598},
	val: "trn",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 111, offset: 3606},
	val: "*",
	ignoreCase: false,
},
&litMatcher{
	pos: position{line: 161, col: 117, offset: 3612},
	val: "ref",
	ignoreCase: false,
},
	},
},
},
},
},
{
	name: "Tagarg",
	pos: position{line: 165, col: 1, offset: 3659},
	expr: &actionExpr{
	pos: position{line: 165, col: 10, offset: 3670},
	run: (*parser).callonTagarg1,
	expr: &seqExpr{
	pos: position{line: 165, col: 10, offset: 3670},
	exprs: []interface{}{
&litMatcher{
	pos: position{line: 165, col: 10, offset: 3670},
	val: " ",
	ignoreCase: false,
},
&labeledExpr{
	pos: position{line: 165, col: 14, offset: 3674},
	label: "v",
	expr: &oneOrMoreExpr{
	pos: position{line: 165, col: 16, offset: 3676},
	expr: &charClassMatcher{
	pos: position{line: 165, col: 16, offset: 3676},
	val: "[a-z0-9=]",
	chars: []rune{'=',},
	ranges: []rune{'a','z','0','9',},
	ignoreCase: false,
	inverted: false,
},
},
},
	},
},
},
},
{
	name: "Newline",
	pos: position{line: 173, col: 1, offset: 3842},
	expr: &actionExpr{
	pos: position{line: 173, col: 11, offset: 3854},
	run: (*parser).callonNewline1,
	expr: &litMatcher{
	pos: position{line: 173, col: 11, offset: 3854},
	val: "\n",
	ignoreCase: false,
},
},
},
{
	name: "Anything",
	pos: position{line: 177, col: 1, offset: 3887},
	expr: &actionExpr{
	pos: position{line: 177, col: 12, offset: 3900},
	run: (*parser).callonAnything1,
	expr: &labeledExpr{
	pos: position{line: 177, col: 12, offset: 3900},
	label: "ch",
	expr: &anyMatcher{
	line: 177, col: 15, offset: 3903,
},
},
},
},
{
	name: "EOF",
	pos: position{line: 181, col: 1, offset: 3946},
	expr: &notExpr{
	pos: position{line: 181, col: 7, offset: 3954},
	expr: &anyMatcher{
	line: 181, col: 8, offset: 3955,
},
},
},
	},
}
func (c *current) onAll1(b interface{}) (interface{}, error) {
    return b.(string), nil
}

func (p *parser) callonAll1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAll1(stack["b"])
}

func (c *current) onBody1(ps interface{}) (interface{}, error) {
    var result strings.Builder
    for _, pitf := range ps.([]interface{}) {
        result.WriteString(pitf.(string))
    }
    return result.String(), nil
}

func (p *parser) callonBody1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBody1(stack["ps"])
}

func (c *current) onEscapeSequence1(ch interface{}) (interface{}, error) {
    return string(ch.([]byte)), nil
}

func (p *parser) callonEscapeSequence1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEscapeSequence1(stack["ch"])
}

func (c *current) onAccent1(ch interface{}) (interface{}, error) {
    buff := norm.NFD.Bytes(ch.([]byte))
    buff = norm.NFD.AppendString(buff, "\u0301")
    buff = norm.NFC.Bytes(buff)
    return fmt.Sprintf(`<v-accent>%s</v-accent>`, buff), nil
}

func (p *parser) callonAccent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAccent1(stack["ch"])
}

func (c *current) onTranscription1(ps interface{}) (interface{}, error) {
    var result strings.Builder
    result.WriteString(`<v-trx>`)
    for _, pitf := range ps.([]interface{}) {
        switch p := pitf.(type) {
        case string:
            result.WriteString(p)
        case []byte:
            r, _ := utf8.DecodeRune(p)
            if r == utf8.RuneError {
                return nil, fmt.Errorf("decode rune from %v", p)
            }
            if rm, ok := transcriptionMap[r]; ok {
                result.WriteString(rm)
            } else {
                result.WriteRune(r)
            }
        default:
            return nil, fmt.Errorf("unknown transcription type %T of %v", p, p)
        }
    }
    result.WriteString(`</v-trx>`)
    return result.String(), nil
}

func (p *parser) callonTranscription1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTranscription1(stack["ps"])
}

func (c *current) onOpenTag1(tnitf, targ interface{}) (interface{}, error) {
    tn := tnitf.(string)
    switch tn {
    case "c":
        if targ == nil {
            targ = "darkgreen"
        }
        return fmt.Sprintf(`<span style="color: %s">`, targ), nil
    case "m1":
        return `<p class="ml-0">`, nil
    case "m2":
        return `<p class="ml-2">`, nil
    case "m3":
        return `<p class="ml-4">`, nil
    case "b":
        return `<strong>`, nil
    case "sup":
        return `<sup>`, nil
    case "p":
        return ``, nil
    case "ex":
        return `<v-ex>`, nil
    case "lang":
        return ``, nil
    case "i":
        return `<i>`, nil
    case "com":
        return `<span style="color: darkgreen">`, nil
    case "trn":
        return ``, nil
    case "*":
        return ``, nil
    case "ref":
        return ``, nil
    default:
        return "unknowntag:"+tn, nil
    }
}

func (p *parser) callonOpenTag1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOpenTag1(stack["tnitf"], stack["targ"])
}

func (c *current) onCloseTag1(tnitf interface{}) (interface{}, error) {
    tn := tnitf.(string)
    switch tn {
    case "c":
        return `</span>`, nil
    case "b":
        return `</strong>`, nil
    case "sup":
        return `</sup>`, nil
    case "p":
        return ``, nil
    case "ex":
        return `</v-ex>`, nil
    case "lang":
        return ``, nil
    case "i":
        return `</i>`, nil
    case "com":
        return `</span>`, nil
    case "trn":
        return ``, nil
    case "m":
        return ``, nil
    case "*":
        return ``, nil
    case "ref":
        return ``, nil
    default:
        return "unknowntag:"+tn, nil
    }
}

func (p *parser) callonCloseTag1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCloseTag1(stack["tnitf"])
}

func (c *current) onTagname1(n interface{}) (interface{}, error) {
    return string(n.([]byte)), nil
}

func (p *parser) callonTagname1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagname1(stack["n"])
}

func (c *current) onTagarg1(v interface{}) (interface{}, error) {
    var result strings.Builder
    for _, vitf := range v.([]interface{}) {
        result.Write(vitf.([]byte))
    }
    return result.String(), nil
}

func (p *parser) callonTagarg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagarg1(stack["v"])
}

func (c *current) onNewline1() (interface{}, error) {
    return `</p>`, nil
}

func (p *parser) callonNewline1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNewline1()
}

func (c *current) onAnything1(ch interface{}) (interface{}, error) {
    return string(ch.([]byte)), nil
}

func (p *parser) callonAnything1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAnything1(stack["ch"])
}


var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule          = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch         = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos    position
	expr   interface{}
	run    func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs: new(errList),
		data: b,
		pt: savepoint{position: position{line: 1}},
		recover: true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v interface{}
	b bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug bool
	depth  int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules  map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth) + ">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth) + "<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth) + "MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth) + "MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}

