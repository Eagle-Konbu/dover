package octopus

import (
	"strings"
	"testing"

	graphql "github.com/hasura/go-graphql-client"
)

func TestObtainTokenMutationTypeName(t *testing.T) {
	var m struct {
		ObtainKrakenToken struct {
			Token string
		} `graphql:"obtainKrakenToken(input: $input)"`
	}
	vars := map[string]any{
		"input": obtainJSONWebTokenInput{},
	}
	query, err := graphql.ConstructMutation(&m, vars)
	if err != nil {
		t.Fatalf("ConstructMutation: %v", err)
	}
	if !strings.Contains(query, "$input:ObtainJSONWebTokenInput!") {
		t.Errorf("mutation = %q, want it to contain %q", query, "$input:ObtainJSONWebTokenInput!")
	}
}
