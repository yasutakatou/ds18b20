package main

import (
    "context"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
	"github.com/yryz/ds18b20"
)

//npx @modelcontextprotocol/inspector --config ./config.json --server temperature
//現在温度を教えてください
//現在温度は暑いでしょうか、寒いでしょうか。何月位の温度ですか？

func main() {
	sensors, err := ds18b20.Sensors()
    if err != nil {
        panic(err)
    }

    s := server.NewMCPServer("MCP温度計", "0.0.1")

    currentTimeTool := mcp.NewTool("temperature",
        mcp.WithDescription("現在温度を返します"),
    )

    s.AddTool(currentTimeTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		for _, sensor := range sensors {
			t, err := ds18b20.Temperature(sensor)
			if err == nil {
				message := fmt.Sprintf("temperature: %.2f°C", t)
				return mcp.NewToolResultText(message), nil
			}
		}
		return mcp.NewToolResultText("no sensor found."), nil
    })

    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("サーバーエラー: %v\n", err)
    }
}