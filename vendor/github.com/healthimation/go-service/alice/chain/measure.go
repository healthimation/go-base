package chain

import (
	"net/http"

	"github.com/healthimation/go-service/alice/middleware"
	"github.com/justinas/alice"
)

// Measurer can setup a Measured chain
type Measurer interface {
	Measure(name string, handler http.Handler) http.HandlerFunc
}

type base struct {
	baseChain alice.Chain
	timer     middleware.Timer
}

// NewBase gets a new measurer with the provided base chain
// Expected usage:
// t, err := middleware.NewNewrelicTimer(env, serviceName, nrKey)
// if err != nil {
// 	log.Fatalf("Could not instantiate newrelic timer: %v", err)
// }
// b := chain.NewBase(alice.New(), t, middleware.NewLogrusLogger(logrus.NewEntry(logrus.New())))
// router.Get("/user", b.Measure("get users", user.Get()))
//
func NewBase(b alice.Chain, timer middleware.Timer, logger middleware.Logger) Measurer {
	c := b.Append(middleware.Recovery, middleware.UserIDInjector, middleware.RequestID, logger.Log)
	return &base{baseChain: c, timer: timer}
}

// NewBaseWithExtras similar to NewBase but allows users to pass in a set of additional constructors to append the the base chain
func NewBaseWithExtras(b alice.Chain, timer middleware.Timer, logger middleware.Logger, constructors ...alice.Constructor) Measurer {
	c := b.Append(middleware.Recovery, middleware.UserIDInjector, middleware.RequestID, logger.Log)
	c = c.Append(constructors...)
	return &base{baseChain: c, timer: timer}
}

// Measure returns a chain that will have metrics measured
func (b *base) Measure(name string, handler http.Handler) http.HandlerFunc {
	if b.timer != nil {
		return b.baseChain.Append(b.timer.Time(name)).Then(handler).ServeHTTP
	}
	return b.baseChain.Then(handler).ServeHTTP
}
