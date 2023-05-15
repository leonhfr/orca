package uci

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/orca/chess"
)

func TestCommandUCI(t *testing.T) {
	name, author := "NAME", "AUTHOR"
	e := new(mockEngine)
	e.On("Options").Return([]Option{
		testOptions[OptionInteger],
	})
	w := &strings.Builder{}
	s := NewState(name, author, w)

	expected := concatenateResponses([]response{
		responseID{name, author},
		testOptions[OptionInteger],
		responseUCIOK{},
	})

	commandUCI{}.run(context.Background(), e, s)

	e.AssertExpectations(t)
	assert.Equal(t, expected, w.String())
}

func TestCommandDebug(t *testing.T) {
	s := NewState("", "", io.Discard)

	for _, tt := range []bool{true, false} {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			e := new(mockEngine)
			commandDebug{on: tt}.run(context.Background(), e, s)

			e.AssertExpectations(t)
			assert.Equal(t, tt, s.debug)
		})
	}
}

func TestCommandIsReady(t *testing.T) {
	tests := []struct {
		err  error
		logs []string
		rr   []response
	}{
		{nil, nil, []response{responseReadyOK{}}},
		{fmt.Errorf("ERROR"), []string{"info string ERROR"}, []response{responseReadyOK{}}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			e := new(mockEngine)
			e.On("Init").Return(tt.err)
			w := &strings.Builder{}
			s := NewState("", "", w)
			want := concatenateResponses(tt.rr)
			if len(tt.logs) > 0 {
				want = strings.Join(tt.logs, "\n") + "\n" + want
			}

			commandIsReady{}.run(context.Background(), e, s)
			time.Sleep(10 * time.Millisecond)

			e.AssertExpectations(t)
			assert.Equal(t, want, w.String())
		})
	}
}

func TestCommandSetOption(t *testing.T) {
	type args struct {
		cmd commandSetOption
		err error
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"valid option",
			args{commandSetOption{"NAME", "VALUE"}, nil},
			[]string{},
		},
		{
			"invalid option",
			args{commandSetOption{"NAME", "VALUE"}, errors.New("ERROR")},
			[]string{"info string ERROR"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(mockEngine)
			e.On("SetOption", tt.args.cmd.name, tt.args.cmd.value).Return(tt.args.err)
			w := &strings.Builder{}
			s := NewState("", "", w)

			commandSetOption{tt.args.cmd.name, tt.args.cmd.value}.run(context.Background(), e, s)

			e.AssertExpectations(t)
			assert.Equal(t, concatenateStrings(tt.want), w.String())
		})
	}
}

func TestCommandPosition(t *testing.T) {
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
			e := new(mockEngine)
			w := &strings.Builder{}
			s := NewState("", "", w)

			tt.c.run(context.Background(), e, s)

			e.AssertExpectations(t)
			assert.Equal(t, tt.want, s.position.String())
			assert.Equal(t, concatenateStrings(tt.r), w.String())
		})
	}
}

func TestCommandGo(t *testing.T) {
	m1 := chess.Move(chess.B1) ^ chess.Move(chess.A3)<<6 ^ chess.Move(chess.NoPiece)<<20
	m2 := chess.Move(chess.E6) ^ chess.Move(chess.E7)<<6 ^ chess.Move(chess.NoPiece)<<20

	output1 := &Output{Score: 1000, PV: []chess.Move{m1}}
	output2 := &Output{Score: 2000, PV: []chess.Move{m1, m2}}

	tests := []struct {
		c  commandGo
		oo []*Output
		rr []response
	}{
		{
			commandGo{},
			[]*Output{output1, output2},
			[]response{
				responseOutput{Output: output1, time: 1 * time.Nanosecond},
				responseOutput{Output: output2, time: 1 * time.Nanosecond},
				responseBestMove{m1},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			e := new(mockEngine)
			oc := make(chan *Output, len(tt.oo))
			for _, o := range tt.oo {
				oc <- o
			}
			close(oc)
			e.On("Search", mock.Anything, mock.Anything, mock.Anything).Return(oc)

			expected := concatenateResponses(tt.rr)
			w := newMockWaitWriter(len(expected))
			s := NewState("", "", w)

			tt.c.run(context.Background(), e, s)

			w.Wait()
			e.AssertExpectations(t)
			assert.Equal(t, expected, w.String())
		})
	}
}

func TestCommandStop(t *testing.T) {
	e := new(mockEngine)
	s := NewState("", "", io.Discard)

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
		case <-s.stop:
			stopCalled = true
		}
	}()

	time.Sleep(10 * time.Millisecond)
	commandStop{}.run(context.Background(), e, s)
	wg.Wait()

	assert.True(t, stopCalled)
}

func TestCommandQuit(t *testing.T) {
	e := new(mockEngine)
	e.On("Close")
	s := NewState("", "", io.Discard)

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
		case <-s.stop:
			stopCalled = true
		}
	}()

	time.Sleep(10 * time.Millisecond)
	commandQuit{}.run(context.Background(), e, s)
	wg.Wait()

	e.AssertExpectations(t)
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
func concatenateResponses(responses []response) string {
	s := make([]string, len(responses))
	for i, r := range responses {
		s[i] = r.String() + "\n"
	}
	return strings.Join(s, "")
}

// mockWaitWriter associates a string builder and a wait group.
type mockWaitWriter struct {
	b   *strings.Builder
	wg  *sync.WaitGroup
	lim int
}

// newMockWaitWriter creates a new mockWaitWriter.
func newMockWaitWriter(lim int) *mockWaitWriter {
	ms := &mockWaitWriter{
		b:   &strings.Builder{},
		wg:  &sync.WaitGroup{},
		lim: lim,
	}
	ms.wg.Add(1)
	return ms
}

// Write implements the io.Writer interface.
func (ms *mockWaitWriter) Write(p []byte) (n int, err error) {
	n, err = ms.b.Write(p)
	if ms.b.Len() >= ms.lim {
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
