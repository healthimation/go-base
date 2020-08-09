package main

import (
	"log"
	"net/http"
	"os"

	"github.com/healthimation/go-aws-config/src/awsconfig"
	"github.com/healthimation/go-env-config/src/balancer"
	"github.com/healthimation/go-service/balancer"
)

// config keys
const (
	configKeyEnvironment = "HMD_ENVIRONMENT"
	configKeyPathPrefix        = "PATH_PREFIX"
	configKeyAppendServiceName = "APPEND_SERVICE_NAME"
)

func main() {
	// pull environment from env vars
	env := os.Getenv(configKeyEnvironment)
	if len(env) == 0 {
		log.Fatal("environment not set")
	}
	// use the default service name to load config
	conf := awsconfig.NewAWSLoader(env, <serviceName>.DefaultServiceName)
	err := conf.Initialize()
	if err != nil {
		log.Fatalf("Couldnt initialize config: %v", err)
	}

	appendServiceName := conf.MustGetBool(configKeyAppendServiceName)
	pathPrefix := conf.MustGetString(configKeyPathPrefix)

	b := balancer.NewSRVBalancer()
	svr := <serviceName>.NewServer(env, <serviceName>.DefaultServiceName, pathPrefix, appendServiceName, conf, b)

	// Start up the server
	log.Printf("Starting %s %s", env, <serviceName>.DefaultServiceName)
	log.Fatal(http.ListenAndServe(":8080", svr.GetRouter()))
}
