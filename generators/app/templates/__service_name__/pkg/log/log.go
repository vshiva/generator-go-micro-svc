<%=licenseText%>
package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Basically doing what this suggests: http://stackoverflow.com/questions/30257622/golang-logrus-how-to-do-a-centralized-configuration

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

var (
	// std is the name of the standard logger in stdlib `log`
	entry = &Entry{logrus.NewEntry(logrus.StandardLogger())}
)

// Fields type, used to pass to `WithFields`
type Fields logrus.Fields

// Formatter interface is used to implement a custom Formatter
type Formatter logrus.Formatter

// Level type
type Level logrus.Level

// Hook to be fired when logging on the logging levels returned from
// `Levels()` on your implementation of the interface.
type Hook logrus.Hook

// Entry is the final or intermediate Logrus logging entry.
type Entry struct {
	*logrus.Entry
}

// WithField wraps logrus
func (e *Entry) WithField(key string, value interface{}) *Entry {
	return &Entry{e.Entry.WithField(key, value)}
}

// WithFields wraps logrus
func (e *Entry) WithFields(fields Fields) *Entry {
	return &Entry{e.Entry.WithFields(logrus.Fields(fields))}
}

// AddFields adds some fields globally to the standard entry
func AddFields(fields Fields) {
	entry = entry.WithFields(fields)
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) {
	logrus.SetFormatter(logrus.Formatter(formatter))
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) {
	logrus.SetLevel(logrus.Level(level))
}

// GetLevel returns the standard logger level.
func GetLevel() Level {
	return Level(logrus.GetLevel())
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook Hook) {
	logrus.AddHook(logrus.Hook(hook))
}

// WithError creates an entry from the standard entry and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	return entry.WithField(logrus.ErrorKey, err)
}

// WithField creates an entry from the standard entry and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	return entry.WithField(key, value)
}

// WithFields creates an entry from the standard entry and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *Entry {
	return entry.WithFields(fields)
}

// Debug logs a message at level Debug on the standard entry.
func Debug(args ...interface{}) {
	entry.Debug(args...)
}

// Print logs a message at level Info on the standard entry.
func Print(args ...interface{}) {
	entry.Print(args...)
}

// Info logs a message at level Info on the standard entry.
func Info(args ...interface{}) {
	entry.Info(args...)
}

// Warn logs a message at level Warn on the standard entry.
func Warn(args ...interface{}) {
	entry.Warn(args...)
}

// Warning logs a message at level Warn on the standard entry.
func Warning(args ...interface{}) {
	entry.Warning(args...)
}

// Error logs a message at level Error on the standard entry.
func Error(args ...interface{}) {
	entry.Error(args...)
}

// Panic logs a message at level Panic on the standard entry.
func Panic(args ...interface{}) {
	entry.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard entry.
func Fatal(args ...interface{}) {
	entry.Fatal(args...)
}

// Debugf logs a message at level Debug on the standard entry.
func Debugf(format string, args ...interface{}) {
	entry.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard entry.
func Printf(format string, args ...interface{}) {
	entry.Printf(format, args...)
}

// Infof logs a message at level Info on the standard entry.
func Infof(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard entry.
func Warnf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard entry.
func Warningf(format string, args ...interface{}) {
	entry.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard entry.
func Errorf(format string, args ...interface{}) {
	entry.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard entry.
func Panicf(format string, args ...interface{}) {
	entry.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard entry.
func Fatalf(format string, args ...interface{}) {
	entry.Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard entry.
func Debugln(args ...interface{}) {
	entry.Debugln(args...)
}

// Println logs a message at level Info on the standard entry.
func Println(args ...interface{}) {
	entry.Println(args...)
}

// Infoln logs a message at level Info on the standard entry.
func Infoln(args ...interface{}) {
	entry.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard entry.
func Warnln(args ...interface{}) {
	entry.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard entry.
func Warningln(args ...interface{}) {
	entry.Warningln(args...)
}

// Errorln logs a message at level Error on the standard entry.
func Errorln(args ...interface{}) {
	entry.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard entry.
func Panicln(args ...interface{}) {
	entry.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard entry.
func Fatalln(args ...interface{}) {
	entry.Fatalln(args...)
}
