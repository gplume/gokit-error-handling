## [errs] Yet Another Golang Error Package... (YAGERP?)

Errors are suposed to makes us better right? Well let's try to make make errors better then...

This package is primarily a go-kit errors package that allows you to either define error content and message, level at runtime or use your own predefined custom errors. The caller field will display accurately the real error caller inside your application. Finally no more needs for the `golang.org/pkg/errors`

Advantages:
- Always provide real error caller and allow to click this info in IDE
- Unclutter your Sumologic logs and allows better error tracking
- Define logging activation by error Level
- Has an `error` compliant type
- Structured `*Err` type that ease unit testing error checking even for customs one.

Let's asume an uppersace string service with this content check for example purpose:
```
func (stringService) Uppercase(s string) (string, error) {
	if _, err := strconv.Atoi(s); err == nil {
		return s, errs.New(http.StatusBadRequest, "uppercase some numbers dude, really??", errs.Info)
	}
    return strings.ToUpper(s), nil
}
```
The `errs.New()` function has a variadic input parameter where you can put an `error` type, `string` message, `int` for level (0-7) and an http.Status code to return.

You can define you error error like this and these are already ready for use inside `errs` package for your convenience (and please PR are accepted to include more of them) but they can be defined anywhere and customised as you wish:
```
var (
	// ErrInternalServer the basic 500 throw it all error...
	ErrInternalServer = &Err{
		Message: "an internal server error occurred please contact the server's administrator",
		Code:    http.StatusInternalServerError,
		Level:   Error,
	}

	ErrNotFound = &Err{
		Message: "not found",
		Code:    http.StatusNotFound,
		Level:   Info,
	}

	// ErrInvalidBody is used when the payload is wrong
	ErrInvalidBody = &Err{
		Message: "invalid body",
		Code:    http.StatusBadRequest,
		Level:   Info,
	}

	// ErrInvalidParameter throwed when a query required parameter is missing or wrong
	ErrInvalidParameter = &Err{
		Message: "invalid parameter",
		Code:    http.StatusBadRequest,
		Level:   Info,
	}

	// ErrEmpty is returned when an input string is empty.
	ErrEmptyParam = &Err{
		Message: "empty parameter",
		Code:    http.StatusBadRequest,
		Err:     errors.New("empty parameter"),
		Level:   Info,
	}
)
```

so a simple straightforward usage would be:
```
return errs.New(errs.ErrEmptyParam)
```
would log this:
```json
{
    "caller": "/home/gp/go/src/github.com/gplume/gokit-error-handling/api/service_validation.go:21",
    "code": 400,
    "error": "empty parameter",
    "level": "Info",
    "message": "empty parameter",
    "http.method": "POST",
    "http.proto": "HTTP/1.1",
    "http.url": "/uppercase",
    "http.user_agent": "PostmanRuntime/7.6.0"
}
```
`caller` is clean –because 99.9% of the time you don't need the full stack– has the right caller and is ctrl-click-able into your console, it will open the appropriate file in you prefered IDE.

`code` will be translated in the right `http` status code to the client.

`message` contain the message displayed to the client under the field `error`

`error` contains the std go `err.Error()` string.

`level` is parametrable at api startup and is defaulted to only log under the Info level.


## Standard Error level:
| Standard | intLevel |
|----------|----------|
| Undefined | 0 |
| Fatal 	| 1 |
| Error 	| 2 |
| Warn      | 3 |
| Info 	    | 4 |
| Debug 	| 5 |
| Trace 	| 6 |


## Errs type
```
type Err struct {
	Err     error
	Message string
	Code    int
	Level   level // from 1 to 7
	Caller  string
	Stack   *Stack
}
```

The `*errs.Err` type has the ability to pass the standard library `error` type by the trick of this line (that you don't need to declare, it's already there in the package): 
```
var _ error = &Err{}
```

so actually:

```
var (
    ErrCustomInvalidParam = &errs.Err{
        Message: "invalid parameter",
        Code:    http.StatusBadRequest,
        Level:   Info,
    }
)

func doSomeService(() error {
    return errs.New(ErrCustomInvalidParam)
}

```
...works !
That is very helpful to incorporate `errs` package in existing codebase and also `go-kit` that force you to use those `Endpoint` `error` type...


## Compositing

You can compose your own specific error at runtime like this:

```
reutrn errs.New(err, "string message displayed to client", http.StatusNotAcceptable, errs.Error)

```
Orders don't matters because of variadics. If you don't throw an http.Status code or message, these will be defaulted to the common error 500 code and message that you can finc in the `errors_def.go` file.




