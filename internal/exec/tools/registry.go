package tools

import (
	"strings"

	"github.com/pafthang/pocketagent/pkgs/ollama"
)

// Handler executes a tool from parsed arguments.
type Handler func(args map[string]interface{}) (string, error)

// Registry maps tool names to implementations.
type Registry map[string]Handler

// Set contains runtime tools and the LLM catalog.
type Set struct {
	Registry Registry
	Catalog  []ollama.Tool
	closers  []func()
}

// Build creates the tool registry and LLM catalog (builtin + MCP tools).
func Build(cfg Config) *Set {
	return BuildWithMCP(cfg, cfg.MCPServers)
}

// BuildMCPOnly connects MCP servers without builtin tools.
func BuildMCPOnly(servers []MCPServerConfig) *Set {
	return buildMCPServers(servers)
}

// BuildWithMCP creates builtins plus the provided MCP servers.
func BuildWithMCP(cfg Config, servers []MCPServerConfig) *Set {
	mcpBindings, mcpCatalog, clients := connectMCPServers(servers)

	reg := Registry{
		"search_web": func(args map[string]interface{}) (string, error) {
			return searchWeb(cfg, ArgString(args, "query", "q", "input"))
		},
		"scrape_page": func(args map[string]interface{}) (string, error) {
			return scrapePage(ArgString(args, "url", "input"))
		},
	}

	if cfg.CodeExecEnabled {
		reg["code_exec"] = func(args map[string]interface{}) (string, error) {
			return codeExec(cfg, args)
		}
	}

	for _, binding := range mcpBindings {
		b := binding
		reg[b.name] = func(args map[string]interface{}) (string, error) {
			return b.client.callTool(b.tool, args)
		}
	}

	closers := make([]func(), 0, len(clients))
	for _, client := range clients {
		c := client
		closers = append(closers, func() { _ = c.close() })
	}

	return &Set{
		Registry: reg,
		Catalog:  append(builtinCatalog(cfg), mcpCatalog...),
		closers:  closers,
	}
}

func buildMCPServers(servers []MCPServerConfig) *Set {
	mcpBindings, mcpCatalog, clients := connectMCPServers(servers)

	reg := Registry{}
	for _, binding := range mcpBindings {
		b := binding
		reg[b.name] = func(args map[string]interface{}) (string, error) {
			return b.client.callTool(b.tool, args)
		}
	}

	closers := make([]func(), 0, len(clients))
	for _, client := range clients {
		c := client
		closers = append(closers, func() { _ = c.close() })
	}

	return &Set{
		Registry: reg,
		Catalog:  mcpCatalog,
		closers:  closers,
	}
}

// Close releases MCP subprocess resources.
func (s *Set) Close() {
	if s == nil {
		return
	}
	for _, closeFn := range s.closers {
		if closeFn != nil {
			closeFn()
		}
	}
}

// DefaultRegistry returns builtin tools with environment-based configuration.
func DefaultRegistry() Registry {
	return Build(LoadFromEnv()).Registry
}

// Catalog returns builtin + MCP tool definitions for the LLM.
func Catalog(cfg Config) []ollama.Tool {
	return Build(cfg).Catalog
}

// Execute runs a tool by name.
func (r Registry) Execute(name, args string) string {
	name = strings.TrimSpace(name)
	fn, ok := r[name]
	if !ok {
		return "Unknown tool: " + name
	}

	result, err := fn(ParseArgs(args))
	if err != nil {
		return "Error: " + err.Error()
	}
	return result
}
