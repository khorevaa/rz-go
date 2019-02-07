// +build !binary_log

package pkgerrors

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/bloom42/rz-go"
	"github.com/pkg/errors"
)

func TestLogStack(t *testing.T) {
	rz.ErrorStackMarshaler = MarshalStack

	out := &bytes.Buffer{}
	log := rz.New(rz.Writer(out), rz.Timestamp(false))

	err := errors.Wrap(errors.New("error message"), "from error")
	log.Log("", func(e *rz.Event) {
		e.Stack().Err(err)
	})

	got := out.String()
	want := `\{"stack":\[\{"func":"TestLogStack","line":"20","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestContextStack(t *testing.T) {
	rz.ErrorStackMarshaler = MarshalStack

	out := &bytes.Buffer{}
	log := rz.New(
		rz.Writer(out),
		rz.Stack(true),
		rz.Timestamp(false),
	)

	err := errors.Wrap(errors.New("error message"), "from error")
	log.Log("", func(e *rz.Event) {
		e.Err(err)
	})

	got := out.String()
	want := `\{"stack":\[\{"func":"TestContextStack","line":"42","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func BenchmarkLogStack(b *testing.B) {
	rz.ErrorStackMarshaler = MarshalStack
	out := &bytes.Buffer{}
	log := rz.New(rz.Writer(out))
	err := errors.Wrap(errors.New("error message"), "from error")
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Log("", func(e *rz.Event) {
			e.Stack().Err(err)
		})
		out.Reset()
	}
}