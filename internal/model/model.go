package model

type Neuron struct {
	Name                 string         `yaml:"name"`
	Type                 string         `yaml:"type"`
	Description          string         `yaml:"description"`
	ExecFile             string         `yaml:"exec_file"`
	PreExecDebug         string         `yaml:"pre_exec_debug"`
	AssertExitStatus     []int          `yaml:"assertExitStatus"`
	PostExecSuccessDebug string         `yaml:"post_exec_success_debug"`
	PostExecFailDebug    map[int]string `yaml:"post_exec_fail_debug"`
}
