package rule

import (
	v1 "go.anx.io/go-anxcloud/pkg/apis/lbaas/v1"
	"go.anx.io/go-anxcloud/pkg/lbaas/common"
)

type Rule = v1.Rule

// Definition describes the Rule object that should be created
type Definition struct {
	Name             string       `json:"name,omitempty"`
	State            common.State `json:"state,omitempty"`
	RuleType         string       `json:"rule_type,omitempty"`
	ParentType       string       `json:"parent_type,omitempty"`
	Frontend         *string      `json:"frontend,omitempty"`
	Backend          *string      `json:"backend,omitempty"`
	Index            int          `json:"index"`
	Condition        string       `json:"condition,omitempty"`
	ConditionTest    string       `json:"condition_test,omitempty"`
	Type             string       `json:"type,omitempty"`
	Action           string       `json:"action,omitempty"`
	RedirectionType  string       `json:"redirection_type,omitempty"`
	RedirectionValue string       `json:"redirection_value,omitempty"`
	RedirectionCode  string       `json:"redirection_code,omitempty"`
}
