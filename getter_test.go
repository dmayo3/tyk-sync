package tyk_git

import (
	"fmt"
	"github.com/TykTechnologies/tyk/apidef"
	"testing"
)

var REPO string = "https://github.com/lonelycode/integration-test.git"

type MockPublisher struct{}

var mockPublish MockPublisher = MockPublisher{}

func (mp MockPublisher) Create(apiDef *apidef.APIDefinition) (string, error) {
	newID := "654321"
	fmt.Printf("Creating API ID: %v (on: %vm to: %v\n)",
		newID,
		apiDef.Proxy.ListenPath,
		apiDef.Proxy.TargetURL)
	return newID, nil
}

func (mp MockPublisher) Update(id string, apiDef *apidef.APIDefinition) error {
	fmt.Printf("Updating API ID: %v (on: %vm to: %v\n)",
		apiDef.APIID,
		apiDef.Proxy.ListenPath,
		apiDef.Proxy.TargetURL)

	return nil
}

func TestNewGGetter(t *testing.T) {
	_, e := NewGGetter(REPO, "refs/heads/master", []byte{}, mockPublish)
	if e != nil {
		t.Fatal(e)
	}
}

func TestGitGetter_FetchRepo(t *testing.T) {
	g, e := NewGGetter(REPO, "refs/heads/master", []byte{}, mockPublish)
	if e != nil {
		t.Fatal(e)
	}

	e = g.FetchRepo()
	if e != nil {
		t.Fatal(e)
	}
}

func TestGitGetter_FetchTykSpec(t *testing.T) {
	g, e := NewGGetter(REPO, "refs/heads/master", []byte{}, mockPublish)
	if e != nil {
		t.Fatal(e)
	}

	e = g.FetchRepo()
	if e != nil {
		t.Fatal(e)
	}

	ts, err := g.FetchTykSpec()
	if err != nil {
		t.Fatal(err)
	}

	if ts.Type != TYPE_APIDEF {
		t.Fatalf("Spec Type is invalid: %v expected: '%v'", ts.Type, TYPE_APIDEF)
	}
}

func TestGitGetter_FetchAPIDef(t *testing.T) {
	g, e := NewGGetter(REPO, "refs/heads/master", []byte{}, mockPublish)
	if e != nil {
		t.Fatal(e)
	}

	e = g.FetchRepo()
	if e != nil {
		t.Fatal(e)
	}

	ts, err := g.FetchTykSpec()
	if err != nil {
		t.Fatal(err)
	}

	ad, err := g.FetchAPIDef(ts)
	if err != nil {
		t.Fatal(err)
	}

	if ad.APIID != ts.Meta.APIID {
		t.Fatalf("APIID Was not properly set, expected: %v, got %v", ts.Meta.APIID, ad.APIID)
	}
}

func TestGitGetter_FetchAPIDef_Swagger(t *testing.T) {
	g, e := NewGGetter(REPO, "refs/heads/swagger-test", []byte{}, mockPublish)
	if e != nil {
		t.Fatal(e)
	}

	e = g.FetchRepo()
	if e != nil {
		t.Fatal(e)
	}

	ts, err := g.FetchTykSpec()
	if err != nil {
		t.Fatal(err)
	}

	if ts.Type != TYPE_OAI {
		t.Fatalf("Spec type setting is unexpected expected: 'oas', got %v", ts.Type)
	}

	ad, err := g.FetchAPIDef(ts)
	if err != nil {
		t.Fatal(err)
	}

	if ad.Name != "Swagger Petstore" {
		t.Fatalf("Name Was not properly set, expected: 'Swagger Petstore', got %v", ad.Name)
	}

	if ad.APIID != ts.Meta.APIID {
		t.Fatalf("APIID Was not properly set, expected: %v, got %v", ts.Meta.APIID, ad.APIID)
	}

	if ad.Proxy.TargetURL != ts.Meta.OAS.OverrideTarget {
		t.Fatalf("Target Was not properly set, got: %v, expected %v", ad.Proxy.TargetURL, ts.Meta.OAS.OverrideTarget)
	}

	if ad.Proxy.ListenPath != ts.Meta.OAS.OverrideListenPath {
		t.Fatalf("Target Was not properly set, expected: %v, got %v", ad.Proxy.ListenPath, ts.Meta.OAS.OverrideListenPath)
	}
}
