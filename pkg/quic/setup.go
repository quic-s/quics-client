package quic

import "github.com/quic-s/quics-client/pkg/viper"

var host string
var port string

func init() {
	host = viper.GetViperEnvVariables("QUICS_SERVER_HOST")
	port = viper.GetViperEnvVariables("QUICS_SERVER_PORT")
}
