package main

import (
	"fmt"
	"os"
	"syscall"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/mitchellh/go-homedir"
	"github.com/prometheus/common/log"
	"github.com/judwhite/go-svc/svc"
	"go-cap/cap"
	"path/filepath"
)

// service example: https://github.com/nsqio/nsq/blob/master/apps/nsqd/nsqd.go

func setupFlags(cmd *cobra.Command, cfgFile *string) {
	flags := cmd.Flags()
	flags.StringVarP(cfgFile, "config", "f", "", "config file (default is $HOME/.capd.yaml)")
	for _, o := range cap.Options {
		val := o.Default
		if o.Type == "bool" {
			if val == nil {
				val = false
			}
			flags.BoolP(o.Name, o.Shorthand, val.(bool), o.Usage)
		} else if o.Type == "string" {
			if val == nil {
				val = ""
			}
			flags.StringP(o.Name, o.Shorthand, val.(string), o.Usage)
		}
	}
}

func initConfig(cmd *cobra.Command, cfgFile string) {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("capd")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Warn(err)
		//os.Exit(1)
	}

	flags := cmd.Flags()

	var val interface{}
	for _, o := range cap.Options {
		if o.Type == "bool" {
			val, _ = flags.GetBool(o.Name)
		} else if o.Type == "string" {
			val, _ = flags.GetString(o.Name)
		}
		if val == nil {
			val = o.Default
		}
		if val != nil {
			viper.Set(o.Name, val)
		}
	}
}

type program struct {
	ap *cap.AP
}


func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Start() error {
	var cfgFile string

	// create command
	cmd := &cobra.Command{
		Use:   "capd",
		Short: "capd is the create_ap daemon",
		Run: func(cmd *cobra.Command, args []string) {
			initConfig(cmd, cfgFile)
			ap, err := cap.CreateAP(viper.AllSettings())
			if err != nil {
				log.Fatal(err)
			}
			p.ap = ap
			ap.Start()
		},
	}

	// init flags
	setupFlags(cmd, &cfgFile)

	// execute
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

	if p.ap == nil {
		os.Exit(0)
	}

	return nil
}

func (p *program) Stop() error {
	if p.ap != nil {
		return p.ap.Stop()
	}
	return nil
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}
