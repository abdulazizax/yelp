package casbin

import "github.com/casbin/casbin"

// NewEnforcer initializes a Casbin enforcer using the provided model and policy files.
func NewEnforcer(modelPath, policyPath string) (*casbin.Enforcer, error) {
	// Create a new Casbin enforcer with the specified model and policy files.
	enforcer := casbin.NewEnforcer(modelPath, policyPath)

	// Load policies from the policy file.
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return enforcer, nil
}
