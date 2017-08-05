package sequence

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrNotSeqfile    = errors.New("not a sequence file")
	ErrFrameExists   = errors.New("frame exists")
	ErrNegativeFrame = errors.New("nagative frame")
)

// Splitter is a file name splitter.
type Splitter struct {
	re *regexp.Regexp
}

// reSplit is default regular expression for Splitter.
var reSplit = regexp.MustCompile(`(.*\D)*(\d+)(.*?)$`)

// NewSplitter creates a new splitter.
func NewSplitter() *Splitter {
	return &Splitter{
		re: reSplit,
	}
}

// SetRegexp let users of this package to make their own splitter.
// Splitter always assume that the regular expression is right.
// So who makes their own splitter should ensure that it is right.
func (s *Splitter) SetRegexp(re *regexp.Regexp) {
	s.re = re
}

// Split takes file name and splits it into 3 parts,
// which is pre, digits, and post.
func (s *Splitter) Split(fname string) (pre, digits, post string, err error) {
	m := s.re.FindStringSubmatch(fname)
	if m == nil {
		return "", "", "", ErrNotSeqfile
	}
	return m[1], m[2], m[3], nil
}

// Fmt{Sharp, DollarF, PrecentD} are pre-defined formatter, that covers most user's need.
var (
	FmtSharp = func(pre, digits, post string) string {
		return pre + strings.Repeat("#", len(digits)) + post
	}
	FmtDollarF = func(pre, digits, post string) string {
		return pre + "$F" + strconv.Itoa(len(digits)) + post
	}
	FmtPercentD = func(pre, digits, post string) string {
		return pre + "%0" + strconv.Itoa(len(digits)) + "d" + post
	}
)

type Manager struct {
	Seqs map[string]*Seq

	splitter   *Splitter
	formatting func(pre, digits, post string) string
}

func NewManager(splitter *Splitter, formatting func(pre, digits, post string) string) *Manager {
	return &Manager{
		Seqs:       make(map[string]*Seq),
		splitter:   splitter,
		formatting: formatting,
	}
}

func (m *Manager) Add(fname string) error {
	pre, digits, post, err := m.splitter.Split(fname)
	if err != nil {
		return err
	}

	name := m.formatting(pre, digits, post)
	frame, _ := strconv.Atoi(digits)

	s, ok := m.Seqs[name]
	if !ok {
		s = NewSeq()
		m.Seqs[name] = s
	}
	return s.AddFrame(frame)
}

func (m *Manager) SeqNames() []string {
	names := []string{}
	for n := range m.Seqs {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func (m *Manager) String() string {
	str := ""
	for _, n := range m.SeqNames() {
		if str != "" {
			str += "\n"
		}
		str += fmt.Sprintf("%s %s", n, m.Seqs[n])
	}
	return str
}

type Seq struct {
	frames map[int]struct{}
}

func NewSeq() *Seq {
	return &Seq{
		frames: make(map[int]struct{}),
	}
}

func (s *Seq) AddFrame(f int) error {
	if f < 0 {
		return ErrNegativeFrame
	}
	if _, ok := s.frames[f]; ok {
		return ErrFrameExists
	}
	s.frames[f] = struct{}{}
	return nil
}

func (s *Seq) Ranges() []*Range {
	if len(s.frames) == 0 {
		return []*Range{}
	}

	frames := []int{}
	for f := range s.frames {
		frames = append(frames, f)
	}
	sort.Ints(frames)

	rngs := []*Range{}
	r := NewRange(frames[0])
	rngs = append(rngs, r)
	for _, f := range frames[1:] {
		ok := r.Extend(f)
		if !ok {
			r = NewRange(f)
			rngs = append(rngs, r)
		}
	}
	return rngs
}

func (s *Seq) String() string {
	str := ""
	rngs := s.Ranges()
	for _, r := range rngs {
		if str != "" {
			str += " "
		}
		if r.minf == r.maxf {
			str += fmt.Sprintf("%d", r.minf)
		} else {
			str += fmt.Sprintf("%d-%d", r.minf, r.maxf)
		}
	}
	return str
}

type Range struct {
	minf int
	maxf int
}

func NewRange(f int) *Range {
	return &Range{
		minf: f,
		maxf: f,
	}
}

func (r *Range) Extend(f int) bool {
	if f != r.maxf+1 {
		return false
	}
	r.maxf = f
	return true
}
