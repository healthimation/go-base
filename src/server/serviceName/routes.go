package <serviceName>

import (
	"github.com/divideandconquer/go-consul-client/src/balancer"
	"github.com/divideandconquer/go-consul-client/src/config"
	"github.com/healthimation/go-service/alice/chain"
	"github.com/healthimation/go-service/alice/middleware"
	"github.com/healthimation/go-service/service"
	"github.com/healthimation/go-client/client"
	"github.com/husobee/vestigo"
	"github.com/justinas/alice"
)

// config keys
const (
	configKeyDBUser     = "HMD_DB_USER"
	configKeyDBPassword = "HMD_DB_PASSWORD"
	configKeyUseCORS    = "HMD_USE_CORS"
	configKeySigningKey = "HMD_SIGNING_KEY"
)

// DefaultServiceName is used in 99% of cases
const DefaultServiceName = "<serviceName>"

type server struct {
	environment string
	serviceName string
	pathPrefix               string
	appendServiceNameToRoute bool
	router      *vestigo.Router
	conf        config.Loader
	balancer    balancer.DNS
}

// NewServer returns a Server
func NewServer(env, serviceName string, pathPrefix string, appendServiceNameToRoute bool, conf config.Loader, lb balancer.DNS) service.Server {
	ret := &server{environment: env, serviceName: serviceName, conf: conf, balancer: lb, pathPrefix: pathPrefix, appendServiceNameToRoute: appendServiceNameToRoute}
	ret.init()
	return ret
}

func (s *server) init() {
	dbUser := s.conf.MustGetString(configKeyDBUser)
	dbPass := s.conf.MustGetString(configKeyDBPassword)
	useCORS := s.conf.MustGetBool(configKeyUseCORS)
	signingKey := []byte(s.conf.MustGetString(configKeySigningKey))

	log := middleware.GetDefaultLogger(s.serviceName, s.environment)

	//initialize the db
	dbFactory := data.GetDBFactory(s.balancer, dbUser, dbPass, s.serviceName, log)

	// To track timer metrics setup and pass in a timer instead of nil
	b := chain.NewBase(alice.New(), nil, middleware.NewLogrusLogger(log, true))

	// error handlers
	vestigo.CustomNotFoundHandlerFunc(chain.NotFoundHandler)
	vestigo.CustomMethodNotAllowedHandlerFunc(chain.MethodNotAllowedHandler)
	router := vestigo.NewRouter()
	if useCORS {
		router.SetGlobalCors(&vestigo.CorsAccessControl{
			AllowOrigin:      []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
			AllowHeaders:     []string{"DNT", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false,
		})
	}

	// Setup routes
	router.Get(s.prefixRoute("/v1/test/ping"), b.Measure("ping", test.Ping()))

	s.router = router
}

func (s *server) GetRouter() *vestigo.Router {
	return s.router
}
func (s *server) GetEnvironment() string {
	return s.environment
}
func (s *server) GetServiceName() string {
	return s.serviceName
}

func (s *server) prefixRoute(route string) string {
	r := client.PrefixRoute(s.serviceName, s.pathPrefix, s.appendServiceNameToRoute, route)
	return r
}
