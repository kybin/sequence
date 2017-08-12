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
	// re catches sequence file name
	// and groups the file name into 3 parts (pre, digits, post).
	// When it does not match, the file will treated as non-sequece file.
	//
	// Note: If it does not have 3 sub groups, it will panic.
	re *regexp.Regexp
}

// reDefaultSplit is regular expression for DefaultSplitter.
// It finds right most digit strings and it's pre and post parts.
var reDefaultSplit = regexp.MustCompile(`(.*\D)*(\d+)(.*?)$`)

// DefaultSplitter is a default splitter for this package.
//
// User could create their own splitter. See NewSplitter.
var DefaultSplitter = NewSplitter(reDefaultSplit)

// NewSplitter creates a new splitter.
//
// Splitter assumes it's regular expression could catch sequence file name
// and groups the file name into 3 parts (pre, digits, post).
// When it does not match, the file will treated as non-sequece file.
//
// Note: If the regular expression does not have 3 sub groups, it will panic.
func NewSplitter(re *regexp.Regexp) *Splitter {
	return &Splitter{
		re: re,
	}
}

// Split takes a file name and splits it into 3 parts,
// which is pre, digits, and post.
// It returns error if the file name does not look like a sequence file.
func (s *Splitter) Split(fname string) (pre, digits, post string, err error) {
	m := s.re.FindStringSubmatch(fname)
	if m == nil {
		return "", "", "", ErrNotSeqfile
	}
	return m[1], m[2], m[3], nil
}

// Fmt{Sharp, DollarF, PrecentD} are pre-defined formatter,
// that covers most user's need.
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

// A Manager is a sequence manager.
type Manager struct {
	Seqs map[string]*Seq

	splitter   *Splitter
	formatting func(pre, digits, post string) string
}

// NewManager creates a new sequence manager.
func NewManager(splitter *Splitter, formatting func(pre, digits, post string) string) *Manager {
	return &Manager{
		Seqs:       make(map[string]*Seq),
		splitter:   splitter,
		formatting: formatting,
	}
}

// Add adds a file to the manager.
//
// If the file's sequence is not exist yet,
// it will create a new sequence automatically.
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

// SeqNames returns it's sequence names in ascending order.
func (m *Manager) SeqNames() []string {
	names := []string{}
	for n := range m.Seqs {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// String returns a string that shows it's sequences.
//
// It will be multiple lines if it has more than one sequence.
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

// A Seq is a frame sequence. It does not hold a sequence name.
type Seq struct {
	frames map[int]struct{}
}

// NewSeq creates a new sequence.
func NewSeq() *Seq {
	return &Seq{
		frames: make(map[int]struct{}),
	}
}

// AddFrame adds a frame into sequence.
//
// It treats negative frames are invalid.
// So returns ErrNegativeFrame when it takes a negative frame.
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

// Ranges converts a sequence to several contiguous ranges.
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

// String expresses a sequence using ranges.
func (s *Seq) String() string {
	str := ""
	rngs := s.Ranges()
	for _, r := range rngs {
		if str != "" {
			str += " "
		}
		str += r.String()
	}
	return str
}

// Range is a contiguous frame range,
// which includes Max frame.
type Range struct {
	Min int
	Max int
}

// NewRange creates a new range.
func NewRange(f int) *Range {
	return &Range{
		Min: f,
		Max: f,
	}
}

// Extend extends a range by one, only if,
// input frame is bigger than current max frame by 1.
// When it extends, it returns true, or it returns false.
func (r *Range) Extend(f int) bool {
	if f != r.Max+1 {
		return false
	}
	r.Max = f
	return true
}

// String expresses the range with dash. Like "1-10".
// But if the min and max are same, it will just show one. Like "5".
func (r *Range) String() string {
	if r.Min == r.Max {
		return fmt.Sprintf("%d", r.Min)
	}
	return fmt.Sprintf("%d-%d", r.Min, r.Max)
}
