package main

import (
	"context"
	"github.com/open-policy-agent/opa/rego"
)

const opa_module_prelude = `
package discourse

import future.keywords.if
import future.keywords.in

default trigger := false

`


type policyInput struct {
	TeamId int64
	TeamScore int64
	AskgodFlags map[string][]int64
	Post *post
}


type opa struct {
	PreparedQuery rego.PreparedEvalQuery
	ctx           context.Context
}

func PreparePolicy(policy string) (*opa, error) {
	ctx := context.TODO()
	query, err := rego.New(
		rego.Query("data.discourse.trigger"),
		rego.Module("discourse.rego", opa_module_prelude+policy),
	).PrepareForEval(ctx)

	if err != nil {
		return nil, err
	}

	return &opa{PreparedQuery: query, ctx: ctx}, nil
}

func (o *opa) Eval(in policyInput) (bool, error) {
	input := map[string]interface{}{
		"team":         map[string]interface{}{"id": in.TeamId, "score": in.TeamScore},
		"askgod_flags": in.AskgodFlags,
		"post":         in.Post,
	}

	result, err := o.PreparedQuery.Eval(o.ctx, rego.EvalInput(input))

	if err != nil {
		// log somethin
		return false, err
	}

	return result.Allowed(), nil
}
