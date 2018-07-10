package log15

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"fmt"
)

// Predefined handlers
var (
	root          *logger
	StdoutHandler = StreamHandler(os.Stdout, LogfmtFormat())
	StderrHandler = StreamHandler(os.Stderr, LogfmtFormat())
)

func init() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		StdoutHandler = StreamHandler(colorable.NewColorableStdout(), TerminalFormat())
	}

	if isatty.IsTerminal(os.Stderr.Fd()) {
		StderrHandler = StreamHandler(colorable.NewColorableStderr(), TerminalFormat())
	}

	root = &logger{[]interface{}{}, new(swapHandler)}
	root.SetHandler(StdoutHandler)
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func New(ctx ...interface{}) Logger {
	return root.New(ctx...)
}

// Root returns the root logger
func Root() Logger {
	return root
}

// The following functions bypass the exported logger methods (logger.Debug,
// etc.) to keep the call depth the same for all paths to logger.write so
// runtime.Caller(2) always refers to the call site in client code.

// Debug is a convenient alias for Root().Debug
func Debug(msg string, ctx ...interface{}) {
	root.write(msg, LvlDebug, ctx)
}

// mimics logrus.Debugf() behaivour
func Debugf(format string, args ...interface{}){
	var emptyCtx []interface{}
	root.write(fmt.Sprintf(format, args...), LvlDebug, emptyCtx)
}

// Info is a convenient alias for Root().Info
func Info(msg string, ctx ...interface{}) {
	root.write(msg, LvlInfo, ctx)
}

// mimics logrus.Infof() behaivour
func Infof(format string, args ...interface{}){
	var emptyCtx []interface{}
	root.write(fmt.Sprintf(format, args...), LvlInfo, emptyCtx)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, ctx ...interface{}) {
	root.write(msg, LvlWarn, ctx)
}

// mimics logrus.Warnf() behaivour
func Warnf(format string, args ...interface{}){
	var emptyCtx []interface{}
	root.write(fmt.Sprintf(format, args...), LvlWarn, emptyCtx)
}

// Error is a convenient alias for Root().Error
func Error(msg string, ctx ...interface{}) {
	root.write(msg, LvlError, ctx)
}

// mimics logrus.Errorf() behaivour
func Errorf(format string, args ...interface{}){
	var emptyCtx []interface{}
	root.write(fmt.Sprintf(format, args...), LvlError, emptyCtx)
}

// Crit is a convenient alias for Root().Crit
func Crit(msg string, ctx ...interface{}) {
	root.write(msg, LvlCrit, ctx)
}

// mimics logrus.Critf() behaivour
func Critf(format string, args ...interface{}){
	var emptyCtx []interface{}
	root.write(fmt.Sprintf(format, args...), LvlCrit, emptyCtx)
}