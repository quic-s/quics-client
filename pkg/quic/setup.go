package quic

import "github.com/quic-s/quics-client/pkg/utils"

var host string
var port string

func init() {
	host = utils.GetViperEnvVariables("QUICS_SERVER_HOST")
	port = utils.GetViperEnvVariables("QUICS_SERVER_PORT")
}
