package policy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ory/hydra/pkg"
	"github.com/ory/ladon"
	"github.com/pkg/errors"
)

type jsonPolicy struct {
	ID          string          `json:"id" gorethink:"id"`
	Description string          `json:"description" gorethink:"description"`
	Subjects    []string        `json:"subjects" gorethink:"subjects"`
	Effect      string          `json:"effect" gorethink:"effect"`
	Resources   []string        `json:"resources" gorethink:"resources"`
	Actions     []string        `json:"actions" gorethink:"actions"`
	Conditions  json.RawMessage `json:"conditions" gorethink:"conditions"`
}

type HTTPManager struct {
	Endpoint *url.URL
	Dry      bool
	Client   *http.Client
}

// Create persists the policy.
func (m *HTTPManager) Create(policy ladon.Policy) error {
	var r = pkg.NewSuperAgent(m.Endpoint.String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Create(policy)
}

// Update the policy.
func (m *HTTPManager) Update(policy ladon.Policy) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, policy.GetID()).String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Update(policy)
}

// Get retrieves a policy.
func (m *HTTPManager) List(limit, offset int64) (ladon.Policies, error) {
	var policies []*ladon.DefaultPolicy
	var r = pkg.NewSuperAgent(m.Endpoint.String() + fmt.Sprintf("?limit=%d&offset=%d", limit, offset))
	r.Client = m.Client
	r.Dry = m.Dry
	if err := r.Get(&policies); err != nil {
		return nil, errors.WithStack(err)
	}

	ret := make(ladon.Policies, len(policies))
	for k, p := range policies {
		ret[k] = ladon.Policy(p)
	}
	return ret, nil
}

// Get retrieves a policy.
func (m *HTTPManager) Get(id string) (ladon.Policy, error) {
	var policy = ladon.DefaultPolicy{
		Conditions: ladon.Conditions{},
	}
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	r.Dry = m.Dry
	if err := r.Get(&policy); err != nil {
		return nil, errors.WithStack(err)
	}

	return &policy, nil
}

// Delete removes a policy.
func (m *HTTPManager) Delete(id string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Delete()
}

// Finds all policies associated with the subject.
func (m *HTTPManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	var policies []*ladon.DefaultPolicy
	var r = pkg.NewSuperAgent(m.Endpoint.String() + "?subject=" + subject)
	r.Client = m.Client
	r.Dry = m.Dry
	if err := r.Get(&policies); err != nil {
		return nil, errors.WithStack(err)
	}

	ret := make(ladon.Policies, len(policies))
	for k, p := range policies {
		ret[k] = ladon.Policy(p)
	}
	return ret, nil
}
