package skasaha

var (
	quietLoggerSingleton Logger = &quietLogger{}
)

// Logger is a basic logging interface.
type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// quietLogger is a no-op logger that implements the Logger interface.
type quietLogger struct{}

func (log *quietLogger) Fatal(v ...interface{}) {
	// No-op.
}

func (log *quietLogger) Fatalf(format string, v ...interface{}) {
	// No-op.
}

func (log *quietLogger) Fatalln(v ...interface{}) {
	// No-op.
}

func (log *quietLogger) Panic(v ...interface{}) {
	// No-op.
}

func (log *quietLogger) Panicf(format string, v ...interface{}) {
	// No-op.
}

func (log *quietLogger) Panicln(v ...interface{}) {
	// No-op.
}

func (log *quietLogger) Print(v ...interface{}) {
	// No-op.
}

func (log *quietLogger) Printf(format string, v ...interface{}) {
	// No-op.
}

func (log *quietLogger) Println(v ...interface{}) {
	// No-op.
}
