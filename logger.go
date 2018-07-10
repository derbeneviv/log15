package log15

import (
	"fmt"
	"time"

	"github.com/go-stack/stack"
)

const timeKey = "t"
const lvlKey = "lvl"
const msgKey = "msg"
const errorKey = "LOG15_ERROR"

// Lvl is a type for predefined log levels.
type Lvl int

// List of predefined log Levels
const (
	LvlCrit Lvl = iota
	LvlError
	LvlWarn
	LvlInfo
	LvlDebug
)

// Returns the name of a Lvl
func (l Lvl) String() string {
	switch l {
	case LvlDebug:
		return "dbug"
	case LvlInfo:
		return "info"
	case LvlWarn:
		return "warn"
	case LvlError:
		return "eror"
	case LvlCrit:
		return "crit"
	default:
		panic("bad level")
	}
}

// LvlFromString returns the appropriate Lvl from a string name.
// Useful for parsing command line args and configuration files.
func LvlFromString(lvlString string) (Lvl, error) {
	switch lvlString {
	case "debug", "dbug":
		return LvlDebug, nil
	case "info":
		return LvlInfo, nil
	case "warn":
		return LvlWarn, nil
	case "error", "eror":
		return LvlError, nil
	case "crit":
		return LvlCrit, nil
	default:
		return LvlDebug, fmt.Errorf("Unknown level: %v", lvlString)
	}
}

// A Record is what a Logger asks its handler to write
type Record struct {
	Time     time.Time
	Lvl      Lvl
	Msg      string
	Ctx      []interface{}
	Call     stack.Call
	KeyNames RecordKeyNames
}

// RecordKeyNames are the predefined names of the log props used by the Logger interface.
type RecordKeyNames struct {
	Time string
	Msg  string
	Lvl  string
}

// A Logger writes key/value pairs to a Handler
type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(lvl Lvl, ctx ...interface{}) Logger

	// GetHandler gets the handler associated with the logger.
	GetHandler() Handler

	// SetHandler updates the logger to write records to the specified handler.
	SetHandler(h Handler)

	// SetLevel updates the logger to set specific max level to write for
	SetLevel(maxLvl Lvl)

	// Log a message at the given level with context key/value pairs
	Debug(msg interface{}, ctx ...interface{})
	Debugf(format string, args ...interface{})
	Info(msg interface{}, ctx ...interface{})
	Infof(format string, args ...interface{})
	Warn(msg interface{}, ctx ...interface{})
	Warnf(format string, args ...interface{})
	Error(msg interface{}, ctx ...interface{})
	Errorf(format string, args ...interface{})
	Crit(msg interface{}, ctx ...interface{})
	Critf(format string, args ...interface{})
	Panic(msg interface{}, ctx ...interface{})
	Panicf(format string, args ...interface{})
}

type logger struct {
	maxLvl Lvl
	ctx []interface{}
	h   *swapHandler
}

func (l *logger) write(msg string, lvl Lvl, ctx []interface{}) {
	if lvl <= l.maxLvl {
		l.h.Log(&Record{
			Time: time.Now(),
			Lvl:  lvl,
			Msg:  msg,
			Ctx:  newContext(l.ctx, ctx),
			Call: stack.Caller(2),
			KeyNames: RecordKeyNames{
				Time: timeKey,
				Msg:  msgKey,
				Lvl:  lvlKey,
			},
		})
	}
}

func (l *logger) New(lvl Lvl, ctx ...interface{}) Logger {
	if lvl == 0 {
		lvl = l.maxLvl
	}
	child := &logger{lvl,newContext(l.ctx, ctx), new(swapHandler)}
	child.SetHandler(l.h)
	return child
}

func newContext(prefix []interface{}, suffix []interface{}) []interface{} {
	normalizedSuffix := normalize(suffix)
	newCtx := make([]interface{}, len(prefix)+len(normalizedSuffix))
	n := copy(newCtx, prefix)
	copy(newCtx[n:], normalizedSuffix)
	return newCtx
}

func (l *logger) Debug(msg interface{}, ctx ...interface{}) {
	l.write(fmt.Sprint(msg), LvlDebug, ctx)
}

func (l *logger) Debugf(format string, args ...interface{}){
	var emptyCtx []interface{}
	l.write(fmt.Sprintf(format, args...), LvlDebug, emptyCtx)
}


func (l *logger) Info(msg interface{}, ctx ...interface{}) {
	l.write(fmt.Sprint(msg), LvlInfo, ctx)
}

func (l *logger) Infof(format string, args ...interface{}){
	var emptyCtx []interface{}
	l.write(fmt.Sprintf(format, args...), LvlInfo, emptyCtx)
}

func (l *logger) Warn(msg interface{}, ctx ...interface{}) {
	l.write(fmt.Sprint(msg), LvlWarn, ctx)
}

func (l *logger) Warnf(format string, args ...interface{}){
	var emptyCtx []interface{}
	l.write(fmt.Sprintf(format, args...), LvlWarn, emptyCtx)
}

func (l *logger) Error(msg interface{}, ctx ...interface{}) {
	l.write(fmt.Sprint(msg), LvlError, ctx)
}

func (l *logger) Errorf(format string, args ...interface{}){
	var emptyCtx []interface{}
	l.write(fmt.Sprintf(format, args...), LvlError, emptyCtx)
}

func (l *logger) Crit(msg interface{}, ctx ...interface{}) {
	l.write(fmt.Sprint(msg), LvlCrit, ctx)
}

func (l *logger) Critf(format string, args ...interface{}){
	var emptyCtx []interface{}
	l.write(fmt.Sprintf(format, args...), LvlCrit, emptyCtx)
}

// these two use Crit level, but also panics and exits the program
func (l *logger) Panic(msg interface{}, ctx ...interface{}) {
	l.write(fmt.Sprint(msg), LvlCrit, ctx)
	// wait for 10 ms to smoothen the output
	time.Sleep(10 )
	panic(fmt.Sprint(msg))
}

func (l *logger) Panicf(format string, args ...interface{}){
	var emptyCtx []interface{}
	l.write(fmt.Sprintf(format, args...), LvlCrit, emptyCtx)
	// wait for 10 ms to smoothen the output
	time.Sleep(10)
	panic(fmt.Sprintf(format, args...))
}

func (l *logger) GetHandler() Handler {
	return l.h.Get()
}

func (l *logger) SetHandler(h Handler) {
	l.h.Swap(h)
}

func (l *logger) SetLevel(maxLvl Lvl) {
	l.maxLvl = maxLvl
}

func normalize(ctx []interface{}) []interface{} {
	// if the caller passed a Ctx object, then expand it
	if len(ctx) == 1 {
		if ctxMap, ok := ctx[0].(Ctx); ok {
			ctx = ctxMap.toArray()
		}
	}

	// ctx needs to be even because it's a series of key/value pairs
	// no one wants to check for errors on logging functions,
	// so instead of erroring on bad input, we'll just make sure
	// that things are the right length and users can fix bugs
	// when they see the output looks wrong
	if len(ctx)%2 != 0 {
		ctx = append(ctx, nil, errorKey, "Normalized odd number of arguments by adding nil")
	}

	return ctx
}

// Lazy allows you to defer calculation of a logged value that is expensive
// to compute until it is certain that it must be evaluated with the given filters.
//
// Lazy may also be used in conjunction with a Logger's New() function
// to generate a child logger which always reports the current value of changing
// state.
//
// You may wrap any function which takes no arguments to Lazy. It may return any
// number of values of any type.
type Lazy struct {
	Fn interface{}
}

// Ctx is a map of key/value pairs to pass as context to a log function
// Use this only if you really need greater safety around the arguments you pass
// to the logging functions.
type Ctx map[string]interface{}

func (c Ctx) toArray() []interface{} {
	arr := make([]interface{}, len(c)*2)

	i := 0
	for k, v := range c {
		arr[i] = k
		arr[i+1] = v
		i += 2
	}

	return arr
}
