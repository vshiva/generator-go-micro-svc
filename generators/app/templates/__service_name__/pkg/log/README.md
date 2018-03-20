# wercker/pkg/log

The package contains several

# Usage

...

# Guidelines

The following are the guidelines that Wercker uses when it comes to logging.
This should be read by Wercker employees, and can be used during code reviews.

### Only use `github.com/wercker/pkg/log`

Always use `github.com/wercker/pkg/log` when writing logs. Do not use
`github.com/sirupsen/logrus`, `log` or `fmt.Print*` for logs.

### Use `log.FromContext`

When a `context.Context` is available always use `log.FromContext` to get a log
instance. The logger retrieved from this func will include all the fields
defined on the global logger, and any fields defined in the context.

Best practice: always instantiate a `logger` as the first line of a context
aware function, then use this `logger` throughout this function. Example:

```go
func (s *SomeService) GetByID(ctx context.Context, id string) (*SomeStruct, error) {
  logger := log.FromContext(ctx)
  
  logger.Info("SomeService.GetByID was called")

  // ...
}
```

### Log messages should be a human-readable sentence

The message that is logged as part of the log event (ie.
`log.Info("<message>")`) should be a normal sentence meant for humans. It
should start with a capitalized letter, but should not end with a `.`. Any
dynamic properties should be favoured to use fields, as these allow for easier
filtering.

Avoid:

```go
logger := log.FromContext(ctx)
logger.Infof("SomeService.GetByID was called; runID: %s, applicationID: %s, userID: %s.", rundID, applicationID, userID)
```

Instead:

```go
logger := log.FromContext(ctx)
logger.WithFields(log.Fields{"runID": rundID, "applicationID": applicationID, "userID": userID}).Infof("SomeService.GetByID  was called")
```

### When returning errors, do not capitalize the first letter

The error message should not capitalize its first letter as it might be wrapped
in another error. In the situtation the error message will be delimited by a
colon, so it will look weird when it contains a capitalized first letter.

Never:

```go
return errors.New("Something failed")
```

Instead:

```go
return errors.New("something failed")
```

### Use errors.Wrap when returning errors

When checking for a error, try to wrap the error with `errors.Wrap` to add more
context to the error before returning it. See
https://godoc.org/github.com/pkg/errors

```go
import "github.com/pkg/errors"

func (s *SomeService) SomeFunc(ctx context.Context, args string) error {
  logger := log.FromContext(ctx)

  result1, err := s.service1.get(id)
  if err != nil {
    return errors.Wrap(err, "something bad happened in service1")
  }

  result2, err := s.service2.Something(args)
  if err != nil {
    return errors.Wrap(err, "something bad happened in service2")
  }

  // ...
}
```

### Always use log.WithError

When logging an error, always add the error object as a field using
log.WithError.

Never:

```go
func main() {
  err := someFunc()
  if err != nil {
    log.Errorf("Unable to execute someFunc: %s", err)
    os.Exit(1)
  }

  // ...
}
```

Instead:

```go
func main() {
  err := someFunc()
  if err != nil {
    log.WithError(err).Error("Unable to execute someFunc")
    os.Exit(1)
  }

  // ...
}
```

## Fields

### Field keys should follow camel casing rules

The keys of the fields should always follow the lower camel casing rules, and
they should only contain letters and numbers. Words in names that are
initialisms or acronyms (e.g. "URL" or "NATO") should have a consistent case
though.

Never:

```go
log.WithField("some_field", "")

log.WithField("UserId", "")
```

Instead:

```go
log.WithField("someField", "")

log.WithField("userID", "")
```

### Field keys can be namespaced

Field keys can optionally be namespaced by dots. The segements should still
follow the previous rule. Try to apply namespaced fields when there are
multiple fields that have the same prefix, unless it diverges from consistent
fields names in other services.

```go
log.WithField("keen.projectID")
```

### Add request specific data to the context

When a context.Context is available and there are some request specific fields,
add them using log.AddFieldToCtx. This allows fields to be logged in unrelated
services. Use this for fields that are relevant for a request.

Example (where `TraceMiddleware` will be called before `MongoStore`)

```go
func (s *TraceMiddleware) Middleware(ctx context.Context) {
  ctx, logger := AddFieldToCtx("traceID", generateTraceID())
  // ...
}

func (s *MongoStore) Func(ctx context.Context) {
  logger := log.FromContext(ctx)
  // logger now contains the traceID field
}
```

### Use consistent names for field keys

Use consistent field keys for the same data that other services are using.

NOTE(bvdberg): how do we keep track and decide which field keys to use?

Avoid:

```go
log.WithField("runIdentifier", "")

log.WithField("run.id", "")
```

Better:

```go
log.WithField("runID", "")
```

## Levels

### Error

Use Error when something failed and the operation was not able to continue.

### Warn

Use Warn when something happened that was not expected, but the operation was able to finish.

### Info

Use Info for normal messages.

### Debug

Use Debug for verbose message, it is possible that these will not be displayed.

## Dates

Note, logging a `time.Time` struct directly as a field will not comply with the
Wercker guidelines.

### Always use the UTC timezone

Convert a time.Time to the UTC timezone before adding it to a field.

Never:

```go
log.WithField("createdAt", time.Now().Format(time.RFC3339))
```

Instead:

```go
log.WithField("createdAt", time.Now().UTC().Format(time.RFC3339))
```

### Always use the time.RFC3339 format

When logging a `time.Time` as a field, always format it to the RFC3339 format.

Never:

```go
log.WithField("createdAt", time.Now().UTC())
```

Instead:

```go
log.WithField("createdAt", time.Now().UTC().Format(time.RFC3339))

```
