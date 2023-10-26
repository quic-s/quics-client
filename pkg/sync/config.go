package sync

import (
	"fmt"

	"github.com/quic-s/quics-client/pkg/viper"
)

func ConfigServer(host string, port string) string {
	viper.WriteViperEnvVariables("QUICS_SERVER_HOST", host)
	viper.WriteViperEnvVariables("QUICS_SERVER_PORT", port)
	return fmt.Sprintf("Changed Host : %s , Changed Port :%s", viper.GetViperEnvVariables("QUICS_SERVER_HOST"), viper.GetViperEnvVariables("QUICS_SERVER_HOST"))
}
