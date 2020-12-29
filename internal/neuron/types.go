package neuron

import log "github.com/anoop2811/cortex/logger"

type Neuron struct {
	logger               *log.StandardLogger
	Name                 string         `yaml:"name"`
	Type                 string         `yaml:"type"`
	Description          string         `yaml:"description"`
	ExecFile             string         `yaml:"exec_file"`
	PreExecDebug         string         `yaml:"pre_exec_debug"`
	AssertExitStatus     []string       `yaml:"assert_exit_status"`
	PostExecSuccessDebug string         `yaml:"post_exec_success_debug"`
	PostExecFailDebug    map[int]string `yaml:"post_exec_fail_debug"`
}

type Definition struct {
	Name   string `yaml:"neuron"`
	Config Config `yaml:"config"`
}

type Config struct {
	Path string         `yaml:"path"`
	Fix  map[int]string `yaml:"fix"`
}
