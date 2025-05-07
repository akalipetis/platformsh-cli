package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/platformsh/cli/internal/config"
	"github.com/spf13/cobra"
)

type LegacyCLIMCPHandler struct {
	c *config.Config
}

func (h *LegacyCLIMCPHandler) Exec(ctx context.Context, args ...string) (*mcp.CallToolResult, error) {
	// Capture the output in a buffer
	var b bytes.Buffer
	c := makeLegacyCLIWrapper(h.c, &b, nil, nil)

	// Execute the command
	if err := c.Exec(ctx, args...); err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	// Return the output as text
	return mcp.NewToolResultText(b.String()), nil
}

// newMCPCommand creates a new MCP command that exposes MCP server functionality
func newMCPCommand(cnf *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: fmt.Sprintf("MCP server for the %s CLI", cnf.Service.Name),
		Long:  fmt.Sprintf("Run the MCP (Model Context Protocol) server for the %s CLI.", cnf.Service.Name),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create MCP server
			s := server.NewMCPServer(
				fmt.Sprintf("%s CLI", cnf.Service.Name),
				cnf.Metadata.Version,
			)

			h := &LegacyCLIMCPHandler{c: cnf}

			// Add list commands tool
			listTool := mcp.NewTool("list_commands",
				mcp.WithDescription("List available CLI commands"),
				mcp.WithBoolean("all",
					mcp.Description("Show all commands, including hidden ones"),
					mcp.DefaultBool(false),
				),
				mcp.WithString("namespace",
					mcp.Description("Filter commands by namespace"),
					mcp.DefaultString(""),
				),
			)
			s.AddTool(listTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				args := []string{"list", "--format=json"}
				if request.Params.Arguments["all"].(bool) {
					args = append(args, "--all")
				}
				if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
					args = append(args, ns)
				}
				return h.Exec(ctx, args...)
			})

			// Add projects tool
			projectsTool := mcp.NewTool("list_projects",
				mcp.WithDescription("List available Upsun projects. Running this command will return a list of all projects that are available to the user."),
			)
			s.AddTool(projectsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return h.Exec(ctx, "project:list")
			})

			// Add project info tool
			projectInfoTool := mcp.NewTool("get_project_info",
				mcp.WithDescription("Get information about a specific Upsun project. If the project ID is not provided, list all projects and then try to search for the correct project using the returned data. The returned information includes project data, billing data, subscription data, etc."),
				mcp.WithString("project_id",
					mcp.Required(),
					mcp.Description("The ID of the project to get information about"),
				),
			)
			s.AddTool(projectInfoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				projectID, ok := request.Params.Arguments["project_id"].(string)
				if !ok || projectID == "" {
					return nil, errors.New("project_id must be a non-empty string")
				}
				return h.Exec(ctx, "project:info", "--project", projectID)
			})

			// Start the stdio server
			if err := server.ServeStdio(s); err != nil {
				fmt.Printf("Server error: %v\n", err)
				return err
			}

			return nil
		},
	}

	return cmd
}
