package common

type ClientOpts struct {
	Environment Environment
}

type Environment string

const EnvironmentDev = Environment("dev")
const EnvironmentStage = Environment("stage")
const EnvironmentProd = Environment("prod")
