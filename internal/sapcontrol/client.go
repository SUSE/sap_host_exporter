package sapcontrol

import (
	"context"
	"net"
	"net/http"

	"github.com/hooklift/gowsdl/soap"
	"github.com/spf13/viper"
)

func NewSoapClient(config *viper.Viper) *soap.Client {
	socket := config.GetString("sap-control-uds")

	if socket != "" {
		udsClient := &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					d := net.Dialer{}
					return d.DialContext(ctx, "unix", socket)
				},
			},
		}
		return soap.NewClient("", soap.WithHTTPClient(udsClient))
	}

	client := soap.NewClient(
		config.GetString("sap-control-url"),
		soap.WithBasicAuth(
			config.GetString("sap-control-user"),
			config.GetString("sap-control-password"),
		),
	)
	return client
}
