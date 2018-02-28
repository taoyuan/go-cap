package cap

import (
	"github.com/go-cmd/cmd"
	"errors"
	"strconv"
	"go-cap/config"
	"github.com/olebedev/emitter"
	"time"
)

var CommandCreateAP = "create_ap"

type AP struct {
	emitter.Emitter
	cmd *cmd.Cmd
	cch <-chan cmd.Status
	wch chan int
}

func CreateAP(m map[string]interface{}) (*AP, error) {
	c := config.CreateMap(m)
	return CreateAPWithConfig(&c)
}

func CreateAPWithConfig(c config.Provider) (*AP, error) {
	var args []string
	options := GetOptions()

	for _, o := range options {
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

	ap := &AP{
		emitter.Emitter{},
		cmd.NewCmd(CommandCreateAP, args...),
		nil,
		nil,
	}
	ap.Use("*", emitter.Void)
	return ap, nil
}

func (ap *AP) Start() <-chan cmd.Status {
	if ap.IsRunning() {
		return ap.cch
	}
	ap.cch = ap.cmd.Start()
	ap.wch = make(chan int)

	go func() {
		outpos := 0
		errpos := 0
		for range time.NewTicker(time.Duration(200) * time.Millisecond).C {
			st := ap.Status()

			if outpos < len(st.Stdout) {
				output := st.Stdout[outpos:]
				outpos = len(st.Stdout)
				ap.Emit("stdout", output)
			}

			if errpos < len(st.Stderr) {
				output := st.Stderr[errpos:]
				errpos = len(st.Stderr)
				ap.Emit("stderr", output)
			}

			if !ap.IsRunning() && outpos >= len(st.Stdout) && errpos >= len(st.Stderr) {
				break
			}
		}
		ap.wch <- 0
	}()

	return ap.cch
}

func (ap *AP) Stop() error {
	return ap.cmd.Stop()
}

func (ap *AP) Wait() error {
	if ap.cch == nil {
		return nil
	}
	status := <-ap.cch
	<- ap.wch
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

func (ap *AP) IsRunning() bool {
	st := ap.cmd.Status()
	return st.StopTs <= 0 && ap.cch != nil
}
