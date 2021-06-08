package configs

import (
	log "github.com/corgi-kx/logcustom"
	"github.com/spf13/viper"
)

type ConfigInfo struct {
	Network
	Blockchain
}

type Network struct {
	ProtocolType string
	ListenHost   string
	ListenPort   string
	ServerCrt    string
	ServerKey    string
}

type Blockchain struct {
	Keys   string
	Blocks string
}

var Configinfos ConfigInfo

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Fatal error config file: %s \n", err)
	}
	Configinfos.Network.ProtocolType = viper.GetString("network.protocol_type")
	Configinfos.Network.ListenHost = viper.GetString("network.listen_host")
	Configinfos.Network.ListenPort = viper.GetString("network.listen_port")
	Configinfos.Network.ServerCrt = viper.GetString("network.server_crt")
	Configinfos.Network.ServerKey = viper.GetString("network.server_key")

	Configinfos.Blockchain.Keys = viper.GetString("blockchain.keys")
	Configinfos.Blockchain.Blocks = viper.GetString("blockchain.blocks")
}
func GetConfig() ConfigInfo {
	return Configinfos
}
