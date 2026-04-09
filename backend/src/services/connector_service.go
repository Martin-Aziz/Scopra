package services

import (
	"context"
	"errors"
	"fmt"
)

type Connector interface {
	Execute(ctx context.Context, action string, payload map[string]any) (map[string]any, error)
}

type ConnectorRegistry struct {
	connectors map[string]Connector
}

func NewConnectorRegistry() *ConnectorRegistry {
	return &ConnectorRegistry{
		connectors: map[string]Connector{
			"github": newStaticConnector("github", []string{"repo.read", "issues.write", "pr.read"}),
			"slack":  newStaticConnector("slack", []string{"channels.read", "chat.write", "users.read"}),
			"jira":   newStaticConnector("jira", []string{"issue.read", "issue.write"}),
			"notion": newStaticConnector("notion", []string{"pages.read", "pages.write"}),
			"linear": newStaticConnector("linear", []string{"issue.read", "issue.write"}),
		},
	}
}

func (r *ConnectorRegistry) Execute(ctx context.Context, tool string, action string, payload map[string]any) (map[string]any, error) {
	connector, ok := r.connectors[tool]
	if !ok {
		return nil, ErrUnsupportedTool
	}
	return connector.Execute(ctx, action, payload)
}

type staticConnector struct {
	name           string
	supportedScope map[string]struct{}
}

func newStaticConnector(name string, allowed []string) *staticConnector {
	supported := make(map[string]struct{}, len(allowed))
	for _, action := range allowed {
		supported[action] = struct{}{}
	}
	return &staticConnector{name: name, supportedScope: supported}
}

func (c *staticConnector) Execute(_ context.Context, action string, payload map[string]any) (map[string]any, error) {
	if _, ok := c.supportedScope[action]; !ok {
		return nil, fmt.Errorf("%w: %s does not support action %s", ErrUnsupportedAction, c.name, action)
	}

	return map[string]any{
		"tool":    c.name,
		"action":  action,
		"payload": payload,
		"status":  "executed",
	}, nil
}

var (
	ErrUnsupportedTool   = errors.New("unsupported tool")
	ErrUnsupportedAction = errors.New("unsupported action")
)
