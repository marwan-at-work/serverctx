// Package serverctx provides net/http Server
// utilities for handling context cancellation signals.
package serverctx

import (
	"context"
	"net/http"
	"time"
)

// Run calls ListenAndServeTLS on s
// and will gracefully shut it down
// when the given ctx is done. Any
// errors will be returned whether it's
// on startup or on shutdown.
func Run(ctx context.Context, s *http.Server, timeout time.Duration) error {
	return RunTLS(ctx, s, timeout, "", "")
}

// RunTLS is like Run but calls ListenAndServeTLS instead.
func RunTLS(ctx context.Context, s *http.Server, timeout time.Duration, certFile, keyFile string) error {
	serverErr := make(chan error)
	go func() {
		if certFile != "" && keyFile != "" {
			serverErr <- s.ListenAndServeTLS(certFile, keyFile)
		} else {
			serverErr <- s.ListenAndServe()
		}
	}()
	var err error
	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err = s.Shutdown(ctx)
	case err = <-serverErr:
	}
	return err
}
