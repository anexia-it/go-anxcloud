package common

type ClientOpts struct {
	Environment Environment
}

func (o *ClientOpts) ApplyToPathFormat(format string) {

}

type Environment string

const EnvironmentDev = Environment("dev")
const EnvironmentStage = Environment("stage")
const EnvironmentProd = Environment("prod")
