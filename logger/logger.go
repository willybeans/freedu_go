package logger

import (
	"context"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once

var log zerolog.Logger

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Get() zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = int(zerolog.InfoLevel) // default to INFO
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FieldsExclude: []string{
				"user_agent",
				"git_revision",
				"go_version",
			},
		}

		if os.Getenv("APP_ENV") != "development" {
			fileLogger := &lumberjack.Logger{
				Filename:   "server.log",
				MaxSize:    5, //
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.GoVersion).
			Logger()
	})

	zerolog.DefaultContextLogger = &log

	return log
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := Get()

		correlationID := xid.New().String()

		ctx := context.WithValue(r.Context(), "correlation_id", correlationID)

		r = r.WithContext(ctx)

		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("correlation_id", correlationID)
		})

		w.Header().Add("X-Correlation-ID", correlationID)

		lrw := newLoggingResponseWriter(w)

		defer func() {
			panicVal := recover()
			if panicVal != nil {
				lrw.statusCode = http.StatusInternalServerError // ensure that the status code is updated
				panic(panicVal)                                 // continue panicking
			}
			l.
				Info().
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Str("user_agent", r.UserAgent()).
				Dur("elapsed_ms", time.Since(start)).
				Int("status_code", lrw.statusCode).
				Msg("incoming request")
		}()
		next.ServeHTTP(lrw, r)
	})
}

// uses their hlog helper for logging instead of writing it yourself
// func RequestLogger(next http.Handler) http.Handler {
// 	l := Get()

// 	h := hlog.NewHandler(l)

// 	accessHandler := hlog.AccessHandler(
// 		func(r *http.Request, status, size int, duration time.Duration) {
// 			hlog.FromRequest(r).Info().
// 				Str("method", r.Method).
// 				Stringer("url", r.URL).
// 				Int("status_code", status).
// 				Int("response_size_bytes", size).
// 				Dur("elapsed_ms", duration).
// 				Msg("incoming request")
// 		},
// 	)

// 	userAgentHandler := hlog.UserAgentHandler("http_user_agent")

// 	return h(accessHandler(userAgentHandler(next)))
// }
