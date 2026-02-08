package uci

import (
	"context"
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/orca/chess"
	"github.com/leonhfr/orca/search"
)

// compile time check that commandUCI implements command.
var _ command = commandUCI{}

func TestCommandUCI(t *testing.T) {
	t.Parallel()
	name, author := "NAME", "AUTHOR"
	e := search.NewEngine()
	w := &strings.Builder{}
	c := NewController(name, author, w)

	expected := concatenateResponses(c, []response{
		responseID{name, author},
		availableSearchOptions[0].response(),
		availableSearchOptions[1].response(),
		availableUCIOptions[0].response(),
		responseUCIOK{},
	})

	commandUCI{}.run(context.Background(), e, c)

	assert.Equal(t, expected, w.String())
}

// compile time check that commandDebug implements command.
var _ command = commandDebug{}

func TestCommandDebug(t *testing.T) {
	t.Parallel()
	c := NewController("", "", io.Discard)

	for _, tt := range []bool{true, false} {
		t.Run(strconv.FormatBool(tt), func(t *testing.T) {
			t.Parallel()
			e := search.NewEngine()
			commandDebug{on: tt}.run(context.Background(), e, c)

			assert.Equal(t, tt, c.debug)
		})
	}
}

// compile time check that commandIsReady implements command.
var _ command = commandIsReady{}

func TestCommandIsReady(t *testing.T) {
	t.Parallel()
	tests := []struct {
		rr []response
	}{
		{[]response{responseReadyOK{}}},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			e := search.NewEngine()

			c := NewController("", "", io.Discard)
			want := concatenateResponses(c, tt.rr)
			w := newMockWaitWriter(len(tt.rr))
			c.writer = w

			commandIsReady{}.run(context.Background(), e, c)
			w.Wait()

			assert.Equal(t, want, w.String())
		})
	}
}

// compile time check that commandSetOption implements command.
var _ command = commandSetOption{}

func TestCommandSetOption(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args commandSetOption
		want []string
	}{
		{
			"valid option",
			commandSetOption{"Hash", "64"},
			[]string{},
		},
		{
			"invalid option",
			commandSetOption{"NAME", "VALUE"},
			[]string{"info string option name not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := search.NewEngine()
			w := &strings.Builder{}
			c := NewController("", "", w)

			commandSetOption{tt.args.name, tt.args.value}.run(context.Background(), e, c)

			assert.Equal(t, concatenateStrings(tt.want), w.String())
		})
	}
}

// compile time check that commandUCINewGame implements command.
var _ command = commandUCINewGame{}

// compile time check that commandPosition implements command.
var _ command = commandPosition{}

func TestCommandPosition(t *testing.T) {
	t.Parallel()
	tests := []struct {
		c    commandPosition
		r    []string
		want string
	}{
		{
			commandPosition{startPos: true},
			[]string{},
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			commandPosition{fen: "2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23"},
			[]string{},
			"2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23",
		},
		{
			commandPosition{fen: "bad fen"},
			[]string{"info string invalid fen (bad fen), must have 6 fields"},
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			commandPosition{
				fen:   "2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23",
				moves: []string{"a4b3"},
			},
			[]string{},
			"2r3k1/1q1nbppp/r3p3/3pP3/2pP4/PpQ2N2/2RN1PPP/2R4K w - - 0 24",
		},
		{
			commandPosition{
				fen:   "2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23",
				moves: []string{"bad move"},
			},
			[]string{"info string invalid move in UCI notation"},
			"2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			e := search.NewEngine()
			w := &strings.Builder{}
			c := NewController("", "", w)

			tt.c.run(context.Background(), e, c)

			assert.Equal(t, tt.want, c.position.String())
			assert.Equal(t, concatenateStrings(tt.r), w.String())
		})
	}
}

// compile time check that commandGo implements command.
var _ command = commandGo{}

func TestCommandGo(t *testing.T) {
	t.Parallel()
	m1 := chess.Move(chess.A2) ^ chess.Move(chess.A3)<<6 ^ chess.Move(chess.NoPiece)<<20
	m2 := chess.Move(chess.E2) ^ chess.Move(chess.E4)<<6 ^ chess.Move(chess.NoPiece)<<20
	m3 := chess.Move(chess.A7) ^ chess.Move(chess.A6)<<6 ^ chess.Move(chess.NoPiece)<<20

	output1 := search.Output{Depth: 1, Nodes: 71, Score: 51, PV: []chess.Move{m1}}
	output2 := search.Output{Depth: 2, Nodes: 391, Score: 6, PV: []chess.Move{m2, m3}}

	tests := []struct {
		c  commandGo
		oo []search.Output
		rr []response
	}{
		{
			commandGo{depth: 2, nodes: 2048},
			[]search.Output{output1, output2},
			[]response{
				responseOutput{Output: output1, time: 1 * time.Nanosecond},
				responseOutput{Output: output2, time: 1 * time.Nanosecond},
				responseBestMove{m2},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			e := search.NewEngine()

			c := NewController("", "", io.Discard)
			expected := concatenateResponses(c, tt.rr)
			w := newMockWaitWriter(len(tt.rr))
			c.writer = w

			tt.c.run(context.Background(), e, c)

			w.Wait()

			timeRegex := regexp.MustCompile(`time \d+`)
			got := timeRegex.ReplaceAllString(w.String(), "time 0")
			assert.Equal(t, expected, got)
		})
	}
}

// compile time check that commandStop implements command.
var _ command = commandStop{}

func TestCommandStop(t *testing.T) {
	t.Parallel()
	e := search.NewEngine()
	c := NewController("", "", io.Discard)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	var stopCalled bool

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		case <-c.stop:
			stopCalled = true
		}
	}()

	time.Sleep(10 * time.Millisecond)
	commandStop{}.run(context.Background(), e, c)
	wg.Wait()

	assert.True(t, stopCalled)
}

// compile time check that commandQuit implements command.
var _ command = commandQuit{}

func TestCommandQuit(t *testing.T) {
	t.Parallel()
	e := search.NewEngine()
	c := NewController("", "", io.Discard)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	var stopCalled bool

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		case <-c.stop:
			stopCalled = true
		}
	}()

	time.Sleep(10 * time.Millisecond)
	commandQuit{}.run(context.Background(), e, c)
	wg.Wait()

	assert.True(t, stopCalled)
}

// concatenateStrings concatenate strings and adds a newline.
func concatenateStrings(ss []string) string {
	res := make([]string, len(ss))
	for i, s := range ss {
		res[i] = s + "\n"
	}
	return strings.Join(res, "")
}

// concatenateResponses concatenate responses and adds a newline.
func concatenateResponses(c *Controller, responses []response) string {
	s := make([]string, len(responses))
	for i, r := range responses {
		s[i] = r.format(c) + "\n"
	}
	return strings.Join(s, "")
}

// mockWaitWriter associates a string builder and a wait group.
type mockWaitWriter struct {
	b  *strings.Builder
	wg *sync.WaitGroup
}

// newMockWaitWriter creates a new mockWaitWriter.
func newMockWaitWriter(lim int) *mockWaitWriter {
	ms := &mockWaitWriter{
		b:  &strings.Builder{},
		wg: &sync.WaitGroup{},
	}
	ms.wg.Add(lim)
	return ms
}

// Write implements the io.Writer interface.
func (ms *mockWaitWriter) Write(p []byte) (n int, err error) {
	n, err = ms.b.Write(p)
	if strings.Contains(string(p), "\n") {
		ms.wg.Done()
	}
	return n, err
}

// String implements the fmt.Stringer interface.
func (ms *mockWaitWriter) String() string {
	return ms.b.String()
}

// Wait waits until the expected string length has been written.
func (ms *mockWaitWriter) Wait() {
	ms.wg.Wait()
}
