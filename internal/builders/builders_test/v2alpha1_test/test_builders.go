package v2alpha1_test

import (
	"github.com/kyma-project/api-gateway/apis/gateway/shared"
	gatewayv2alpha1 "github.com/kyma-project/api-gateway/apis/gateway/v2alpha1"
	"k8s.io/utils/ptr"
	"net/http"
)

type RuleBuilder struct {
	rule *gatewayv2alpha1.Rule
}

func (r *RuleBuilder) WithPath(path string) *RuleBuilder {
	r.rule.Path = path
	return r
}

func (r *RuleBuilder) WithTimeout(timeout uint32) *RuleBuilder {
	r.rule.Timeout = ptr.To(gatewayv2alpha1.Timeout(timeout))
	return r
}

func (r *RuleBuilder) WithService(name, namespace string, port uint32) *RuleBuilder {
	r.rule.Service = &gatewayv2alpha1.Service{
		Name:      &name,
		Namespace: &namespace,
		Port:      &port,
	}
	return r
}

func (r *RuleBuilder) WithMethods(methods ...gatewayv2alpha1.HttpMethod) *RuleBuilder {
	r.rule.Methods = methods
	return r
}

func (r *RuleBuilder) NoAuth() *RuleBuilder {
	r.rule.NoAuth = ptr.To(true)
	return r
}

func (r *RuleBuilder) WithJWTAuthn(issuer, jwksUri string, fromHeaders []*shared.JwtHeader, fromParams []string) *RuleBuilder {
	if r.rule.Jwt == nil {
		r.rule.Jwt = &shared.JwtConfig{}
	}
	r.rule.Jwt.Authentications = append(r.rule.Jwt.Authentications, &shared.JwtAuthentication{
		Issuer:      issuer,
		JwksUri:     jwksUri,
		FromHeaders: fromHeaders,
		FromParams:  fromParams,
	})

	return r
}

func (r *RuleBuilder) WithJWTAuthz(requiredScopes []string, audiences []string) *RuleBuilder {
	if r.rule.Jwt == nil {
		r.rule.Jwt = &shared.JwtConfig{}
	}

	r.rule.Jwt.Authorizations = append(r.rule.Jwt.Authorizations, &shared.JwtAuthorization{
		RequiredScopes: requiredScopes,
		Audiences:      audiences,
	})

	return r
}

func (r *RuleBuilder) WithRequest(rm *gatewayv2alpha1.Request) *RuleBuilder {
	r.rule.Request = rm
	return r
}

func (r *RuleBuilder) WithExtAuth(auth *gatewayv2alpha1.ExtAuth) *RuleBuilder {
	r.rule.ExtAuth = auth
	return r
}

func NewRuleBuilder() *RuleBuilder {
	return &RuleBuilder{
		rule: &gatewayv2alpha1.Rule{},
	}
}

func (r *RuleBuilder) Build() *gatewayv2alpha1.Rule {
	return r.rule
}

type RequestBuilder struct {
	request *gatewayv2alpha1.Request
}

func (m *RequestBuilder) WithHeaders(headers map[string]string) *RequestBuilder {
	m.request.Headers = headers

	return m
}

func (m *RequestBuilder) WithCookies(cookies map[string]string) *RequestBuilder {
	m.request.Cookies = cookies

	return m
}

func NewRequestModifier() *RequestBuilder {
	return &RequestBuilder{
		request: &gatewayv2alpha1.Request{},
	}
}

func (m *RequestBuilder) Build() *gatewayv2alpha1.Request {
	return m.request
}

type ApiRuleBuilder struct {
	apiRule *gatewayv2alpha1.APIRule
}

func (a *ApiRuleBuilder) WithHost(host string) *ApiRuleBuilder {
	a.apiRule.Spec.Hosts = append(a.apiRule.Spec.Hosts, ptr.To(gatewayv2alpha1.Host(host)))
	return a
}

func (a *ApiRuleBuilder) WithHosts(hosts ...string) *ApiRuleBuilder {
	for _, host := range hosts {
		a.WithHost(host)
	}
	return a
}

func (a *ApiRuleBuilder) WithService(name, namespace string, port uint32) *ApiRuleBuilder {
	a.apiRule.Spec.Service = &gatewayv2alpha1.Service{
		Name:      &name,
		Namespace: &namespace,
		Port:      &port,
	}
	return a
}

func (a *ApiRuleBuilder) WithGateway(gateway string) *ApiRuleBuilder {
	a.apiRule.Spec.Gateway = ptr.To(gateway)
	return a
}

func (a *ApiRuleBuilder) WithCORSPolicy(policy gatewayv2alpha1.CorsPolicy) *ApiRuleBuilder {
	a.apiRule.Spec.CorsPolicy = &policy
	return a
}

func (a *ApiRuleBuilder) WithTimeout(timeout uint32) *ApiRuleBuilder {
	a.apiRule.Spec.Timeout = ptr.To(gatewayv2alpha1.Timeout(timeout))
	return a
}

func (a *ApiRuleBuilder) WithRule(rule gatewayv2alpha1.Rule) *ApiRuleBuilder {
	a.apiRule.Spec.Rules = append(a.apiRule.Spec.Rules, rule)
	return a
}

func (a *ApiRuleBuilder) WithRules(rules ...*gatewayv2alpha1.Rule) *ApiRuleBuilder {
	for _, rule := range rules {
		a.WithRule(*rule)
	}
	return a
}

func (a *ApiRuleBuilder) Build() *gatewayv2alpha1.APIRule {
	return a.apiRule
}

func NewAPIRuleBuilder() *ApiRuleBuilder {
	return &ApiRuleBuilder{
		apiRule: &gatewayv2alpha1.APIRule{},
	}
}

// NewAPIRuleBuilderWithDummyDataWithNoAuthRule returns an APIRuleBuilder pre-filled with placeholder data:
//
// Host: example-host.example.com
//
// Gateway: example-namespace/example-gateway
//
// Service: example-namespace/example-service:8080
//
// Rule: GET /
//
// Strategy: NoAuth
func NewAPIRuleBuilderWithDummyDataWithNoAuthRule() *ApiRuleBuilder {
	return NewAPIRuleBuilder().
		WithHost("example-host.example.com").
		WithGateway("example-namespace/example-gateway").
		WithService("example-service", "example-namespace", 8080).
		WithRule(*NewRuleBuilder().WithMethods(http.MethodGet).WithPath("/").NoAuth().Build())
}

func NewAPIRuleBuilderWithDummyData() *ApiRuleBuilder {
	return NewAPIRuleBuilder().
		WithHost("example-host.example.com").
		WithGateway("example-namespace/example-gateway").
		WithService("example-service", "example-namespace", 8080)
}

type corsPolicyBuilder struct {
	policy gatewayv2alpha1.CorsPolicy
}

func (c *corsPolicyBuilder) WithAllowOrigins(origins []map[string]string) *corsPolicyBuilder {
	c.policy.AllowOrigins = origins
	return c
}

func (c *corsPolicyBuilder) WithAllowMethods(methods []string) *corsPolicyBuilder {
	c.policy.AllowMethods = methods
	return c
}

func (c *corsPolicyBuilder) WithAllowHeaders(headers []string) *corsPolicyBuilder {
	c.policy.AllowHeaders = headers
	return c
}

func (c *corsPolicyBuilder) WithExposeHeaders(headers []string) *corsPolicyBuilder {
	c.policy.ExposeHeaders = headers
	return c
}

func (c *corsPolicyBuilder) WithMaxAge(maxAge uint64) *corsPolicyBuilder {
	c.policy.MaxAge = &maxAge
	return c
}

func (c *corsPolicyBuilder) WithAllowCredentials(allow bool) *corsPolicyBuilder {
	c.policy.AllowCredentials = &allow
	return c
}

func NewCorsPolicyBuilder() *corsPolicyBuilder {
	return &corsPolicyBuilder{
		policy: gatewayv2alpha1.CorsPolicy{},
	}
}

func (c *corsPolicyBuilder) Build() gatewayv2alpha1.CorsPolicy {
	return c.policy
}

type ExtAuthBuilder struct {
	extAuth *gatewayv2alpha1.ExtAuth
}

func NewExtAuthBuilder() *ExtAuthBuilder {
	return &ExtAuthBuilder{
		extAuth: &gatewayv2alpha1.ExtAuth{},
	}
}

func (e *ExtAuthBuilder) Build() *gatewayv2alpha1.ExtAuth {
	return e.extAuth
}

func (e *ExtAuthBuilder) WithAuthorizers(auths ...string) *ExtAuthBuilder {
	e.extAuth.ExternalAuthorizers = append(e.extAuth.ExternalAuthorizers, auths...)
	return e
}

func (e *ExtAuthBuilder) WithRestriction(config *shared.JwtConfig) *ExtAuthBuilder {
	e.extAuth.Restrictions = config
	return e
}
