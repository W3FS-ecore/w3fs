package mock

import (
	"context"
	"fmt"
	"go.uber.org/multierr"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// A Hook is a pair of start and stop callbacks, either of which can be nil,
// plus a string identifying the supplier of the hook.
type Hook struct {
	OnStart func(context.Context) error
	OnStop  func(context.Context) error
	caller  string
}

// Lifecycle coordinates application lifecycle hooks.
type Lifecycle struct {
	logger     *Logger
	hooks      []Hook
	numStarted int
}

// New constructs a new Lifecycle.
func New(logger *Logger) *Lifecycle {
	if logger == nil {
		logger = NewLogger()
	}
	return &Lifecycle{logger: logger}
}

// Append adds a Hook to the lifecycle.
func (l *Lifecycle) Append(hook Hook) {
	hook.caller = Caller()
	l.hooks = append(l.hooks, hook)
}

// Start runs all OnStart hooks, returning immediately if it encounters an
// error.
func (l *Lifecycle) Start(ctx context.Context) error {
	for _, hook := range l.hooks {
		if hook.OnStart != nil {
			l.logger.Printf("START\t\t%s()", hook.caller)
			if err := hook.OnStart(ctx); err != nil {
				return err
			}
		}
		l.numStarted++
	}
	return nil
}

// Stop runs any OnStop hooks whose OnStart counterpart succeeded. OnStop
// hooks run in reverse order.
func (l *Lifecycle) Stop(ctx context.Context) error {
	var errs []error
	// Run backward from last successful OnStart.
	for ; l.numStarted > 0; l.numStarted-- {
		hook := l.hooks[l.numStarted-1]
		if hook.OnStop == nil {
			continue
		}
		l.logger.Printf("STOP\t\t%s()", hook.caller)
		if err := hook.OnStop(ctx); err != nil {
			// For best-effort cleanup, keep going after errors.
			errs = append(errs, err)
		}
	}
	return multierr.Combine(errs...)
}

// Match from beginning of the line until the first `vendor/` (non-greedy)
var vendorRe = regexp.MustCompile("^.*?/vendor/")

// sanitize makes the function name suitable for logging display. It removes
// url-encoded elements from the `dot.git` package names and shortens the
// vendored paths.
func sanitize(function string) string {
	// Use the stdlib to un-escape any package import paths which can happen
	// in the case of the "dot-git" postfix. Seems like a bug in stdlib =/
	if unescaped, err := url.QueryUnescape(function); err == nil {
		function = unescaped
	}

	// strip everything prior to the vendor
	return vendorRe.ReplaceAllString(function, "vendor/")
}

// Caller returns the formatted calling func name
func Caller() string {
	// Ascend at most 8 frames looking for a caller outside fx.
	pcs := make([]uintptr, 8)

	// Don't include this frame.
	n := runtime.Callers(2, pcs)
	if n == 0 {
		return "n/a"
	}

	frames := runtime.CallersFrames(pcs)
	for f, more := frames.Next(); more; f, more = frames.Next() {
		if shouldIgnoreFrame(f) {
			continue
		}
		return sanitize(f.Function)
	}
	return "n/a"
}

// FuncName returns a funcs formatted name
func FuncName(fn interface{}) string {
	fnV := reflect.ValueOf(fn)
	if fnV.Kind() != reflect.Func {
		return "n/a"
	}

	function := runtime.FuncForPC(fnV.Pointer()).Name()
	return fmt.Sprintf("%s()", sanitize(function))
}

func isErr(t reflect.Type) bool {
	errInterface := reflect.TypeOf((*error)(nil)).Elem()
	return t.Implements(errInterface)
}

// Ascend the call stack until we leave the Fx production code. This allows us
// to avoid hard-coding a frame skip, which makes this code work well even
// when it's wrapped.
func shouldIgnoreFrame(f runtime.Frame) bool {
	if strings.Contains(f.File, "_test.go") {
		return false
	}
	if strings.Contains(f.File, "go.uber.org/fx") {
		return true
	}
	return false
}
