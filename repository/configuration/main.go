package configuration

import (
	"flag"
	"log"
	"os"
	"os/user"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Configuration interface {
	Get(string) (string, error)
	Init(set *flag.FlagSet)
}

type ViperConfiguration struct {
}

func (vc *ViperConfiguration) setDefaults() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("setDefaults: %+v\n", err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("setDefaults: %+v\n", err)
	}
	viper.SetDefault("name", "Mr. Keeper")
	viper.SetDefault("email", usr.Name+"@"+hostname)
}

func (vc *ViperConfiguration) Init(cmd *flag.FlagSet) {
	vc.setDefaults()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("an error occured while running viper.BindPFlags(): %+v\n", err)
	}
}

func (vc *ViperConfiguration) Get(param string) (string, error) {
	return viper.GetString(param), nil
}

func NewConfiguration() (cfg *ViperConfiguration) {
	cfg = &ViperConfiguration{}
	return

}
