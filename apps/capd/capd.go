package main

import (
	"fmt"
	"os"
	"syscall"
	"github.com/spf13/cobra"
	"github.com/prometheus/common/log"
	"github.com/judwhite/go-svc/svc"
	"path/filepath"
	"go-cap/cap"
	"github.com/olebedev/emitter"
	"go-cap/config"
)

// service example: https://github.com/nsqio/nsq/blob/master/apps/nsqd/nsqd.go

func setupFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringP("config", "f", "", "config file")
	flags.BoolP("verbose", "v", false,"Print create_ap output")
	for _, o := range cap.GetOptions() {
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

func initConfig(cmd *cobra.Command) *config.Config {
	cfg := config.New()
	flags := cmd.Flags()

	cfg.AddFile("/etc/capd/capd.yaml", true)
	cfg.AddFile("./capd.yaml", true)

	cfgFile, _ := flags.GetString("config")
	if cfgFile != "" {
		cfg.AddFile(cfgFile, false)
	}

	//if cfgFile != "" {
	//	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
	//		log.Fatalf("config file (%s) not exist", cfgFile)
	//	}
	//} else {
	//	cfgFile = "./capd.yaml"
	//	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
	//		cfgFile = "/etc/capd/capd.yaml"
	//	}
	//}
	//
	//// config file exist
	//if _, err := os.Stat(cfgFile); err == nil {
	//	viper.SetConfigFile(cfgFile)
	//
	//}
	//
	//if err := viper.ReadInConfig(); err != nil {
	//	log.Warn(err)
	//	//os.Exit(1)
	//}


	var val interface{}
	for _, o := range cap.GetOptions() {
		if o.Type == "bool" {
			val, _ = flags.GetBool(o.Name)
		} else if o.Type == "string" {
			val, _ = flags.GetString(o.Name)
		}
		if val == nil {
			val = o.Default
		}
		if val != nil {
			cfg.Set(o.Name, val)
		}
	}

	return cfg
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

	// create command
	cmd := &cobra.Command{
		Use:   "capd",
		Short: "capd is the create_ap daemon",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := initConfig(cmd)
			ap, err := cap.CreateAP(cfg.AllSettings())
			if err != nil {
				log.Fatal(err)
			}
			p.ap = ap

			ap.On("stderr", func (event *emitter.Event) {
				lines := event.Args[0].([]string)
				for _, line := range lines {
					fmt.Println(line)
				}
			})

			flags := cmd.Flags()
			if ok, _ := flags.GetBool("verbose"); ok {
				ap.On("stdout", func (event *emitter.Event) {
					lines := event.Args[0].([]string)
					for _, line := range lines {
						fmt.Println(line)
					}
				})
			}

			ap.Start()

			err = ap.Wait()
			if err != nil {
				fmt.Errorf("%s", err)
				os.Exit(1)
				return
			}
			os.Exit(0)
		},
	}

	// init flags
	setupFlags(cmd)

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
