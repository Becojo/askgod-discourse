package main

import (
	"testing"
)

func makePost(policy string) post {
	return post{
		Type:  "post",
		Topic: "topic",
		Trigger: &postTrigger{
			Type:   "opa",
			Policy: policy,
		},
	}
}

func TestEval(t *testing.T) {
	input := policyInput{
		TeamId: 105,
		TeamScore: 12,
		Post: nil,
		AskgodFlags: map[string][]int64{"flag1": []int64{105}},
	}
	trigger := `
      trigger if {
         input.team.id == 105
         input.team.score == 12
	     input.team.id in input.askgod_flags["flag1"]
      }
    `

	policy, err := PreparePolicy(trigger)

	if err != nil {
		t.Errorf("PreparePolicy err: %v", err)
	}

	allowed, err := policy.Eval(input)

	if err != nil {
		t.Errorf("policy.Eval err: %v", err)
	}

	if allowed == false {
		t.Errorf("policy should return true")
	}
}
