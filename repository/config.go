package repository

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/user"
)

type Configuration interface {
	Get(string) (string, error)
	Init(set *flag.FlagSet)
}


type viperConfiguration struct {

}

func (vc *viperConfiguration) setDefaults() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("setDefaults: %+v\n", err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("setDefaults: %+v\n", err)
	}
	viper.SetDefault("name", "Mr. Keeper")
	viper.SetDefault("email", usr.Name + "@" + hostname)
}

func (vc *viperConfiguration) Init(cmd *flag.FlagSet) {
	vc.setDefaults()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func (vc *viperConfiguration) Get(param string) (string, error) {
	return 	viper.GetString(param), nil
}

func NewConfiguration() (cfg *viperConfiguration){
	cfg = &viperConfiguration{}
	return

}
