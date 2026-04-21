package rule

import (
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

type Rule = v1.Rule

// Definition describes the Rule object that should be created
type Definition struct {
	Name             string       `json:"name"`
	State            common.State `json:"state"`
	RuleType         string       `json:"rule_type"`
	ParentType       string       `json:"parent_type"`
	Frontend         *string      `json:"frontend,omitempty"`
	Backend          *string      `json:"backend,omitempty"`
	Index            int          `json:"index"`
	Condition        string       `json:"condition"`
	ConditionTest    string       `json:"condition_test"`
	Type             string       `json:"type"`
	Action           string       `json:"action"`
	RedirectionType  string       `json:"redirection_type"`
	RedirectionValue string       `json:"redirection_value"`
	RedirectionCode  string       `json:"redirection_code"`
}
