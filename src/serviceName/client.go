package <serviceName>

import (
	"time"

	"github.com/divideandconquer/go-consul-client/src/balancer"
	"github.com/healthimation/go-client/client"
)

const serviceName = "<serviceName>"

// Client is a client that can interact with the profile service
type Client interface {
}

type serviceClient struct {
	c client.BaseClient
	prependServiceNameToRoute bool
}

// NewClient will create a new Client
func NewClient(lb balancer.DNS, useTLS bool) Client {
	// return &serviceClient{c: client.NewBaseClient(lb.GetHttpUrl, serviceName, useTLS, 10*time.Second)}
	return &serviceClient{
		c:                         client.NewBaseClient(lb.GetHttpUrl, serviceName, useTLS, 10*time.Second, nil),
		prependServiceNameToRoute: prependServiceNameToRoute,
	}
}
