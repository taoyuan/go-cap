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
	cmd  *cmd.Cmd
	ch   <-chan cmd.Status
}

type OutListener func (lines []string)

func wrapMiddleware(listener OutListener) func (event *emitter.Event) {
	return func (event *emitter.Event) {
		if len(event.Args) > 0 {
			listener(event.Args[0].([]string))
		}
	}
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
	}
	ap.Use("*", emitter.Void)
	return ap, nil
}

func (ap *AP) Start() <-chan cmd.Status {
	if ap.IsRunning() {
		return ap.ch
	}
	ap.ch = ap.cmd.Start()

	go func() {
		curr := 0
		for range time.NewTicker(time.Duration(200) * time.Millisecond).C {
			if st := ap.cmd.Status(); ap.IsRunning() && curr < len(st.Stdout) {
				output := st.Stdout[curr:]
				curr = len(st.Stdout)
				ap.Emit("stdout", output)
			}
		}

	}()

	return ap.ch
}

func (ap *AP) Stop() error {
	if !ap.IsRunning() {
		return errors.New("ap: not started")
	}
	return ap.cmd.Stop()
}

func (ap *AP) Wait() error {
	if !ap.IsRunning() {
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

func (ap *AP) IsRunning() bool {
	return ap.ch != nil
}
