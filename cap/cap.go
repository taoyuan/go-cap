package cap

import (
	"github.com/go-cmd/cmd"
	"errors"
	"strconv"
	"go-cap/config"
	"github.com/olebedev/emitter"
)

var CommandCreateAP = "create_ap"

type AP struct {
	cmd  *cmd.Cmd
	ch   <-chan cmd.Status
	curr int
	e    *emitter.Emitter
}

func CreateAP(m map[string]interface{}) (*AP, error) {
	c := config.CreateMap(m)
	return CreateAPWithConfig(&c)
}

func CreateAPWithConfig(c config.Provider) (*AP, error) {
	var args []string

	for _, o := range Options {
		if !c.IsSet(o.Name) {
			continue
		}

		var err error
		var v = c.Get(o.Name)
		if o.Resolve != nil {
			v, err = o.Resolve(v)
			if err != nil {
				return nil, err
			}
		}

		if o.CapOpt == "" {
			if v.(string) != "" {
				args = append(args, v.(string))
			}
			continue
		}

		if o.Type == "bool" && v.(bool) {
			args = append(args, o.CapOpt)
			continue
		}

		if o.Type == "string" && v.(string) != "" {
			args = append(args, o.CapOpt, v.(string))
			continue
		}
	}

	return &AP{
		cmd:  cmd.NewCmd(CommandCreateAP, args...),
		ch:   nil,
		curr: 0,
	}, nil

}

func (ap *AP) Start() <-chan cmd.Status {
	if ap.ch != nil {
		return ap.ch
	}
	ap.ch = ap.cmd.Start()
	ap.curr = 0
	return ap.ch
}

func (ap *AP) Stop() error {
	if ap.ch == nil {
		return errors.New("ap: not started")
	}
	return ap.cmd.Stop()
}

func (ap *AP) Wait() error {
	if ap.ch == nil {
		return errors.New("ap: not started")
	}
	status := <-ap.ch
	if status.Error != nil {
		return status.Error
	}
	if status.Exit != 0 {
		return errors.New("exit status " + strconv.Itoa(status.Exit))
	}
	return nil
}

func (ap *AP) Status() cmd.Status {
	return ap.cmd.Status()
}

func (ap *AP) Output() []string {
	status := ap.cmd.Status()
	if ap.curr < len(status.Stdout) {
		output := status.Stdout[ap.curr:]
		ap.curr = len(status.Stdout)
		return output
	}
	return nil
}
