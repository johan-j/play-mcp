#!/bin/bash

# VS Code MCP Configuration Helper

echo "🔧 Setting up VS Code for MCP Server connection..."

# Create .vscode directory if it doesn't exist
if [ ! -d ".vscode" ]; then
    echo "📁 Creating .vscode directory..."
    mkdir -p .vscode
fi

# Copy MCP configuration to VS Code settings
echo "⚙️  Configuring VS Code settings..."
cp vscode-mcp-config.json .vscode/settings.json

echo "✅ VS Code configuration complete!"
echo ""
echo "📋 Next steps:"
echo "1. Open this project in VS Code: code ."
echo "2. Install an MCP extension from the marketplace"
echo "3. Start the MCP server: ./bin/mcp-server"
echo "4. Look for MCP connection indicator in VS Code status bar"
echo ""
echo "🧪 Test connection:"
echo "- Start server: ./bin/mcp-server"
echo "- Open VS Code in this directory: code ."
echo "- Try asking Copilot: 'Get Apple stock data'"
echo ""
echo "🔍 Troubleshooting:"
echo "- Check server is running: curl http://localhost:8080/health"
echo "- Check VS Code Developer Console for MCP logs"
echo "- Run server in debug mode: ./bin/mcp-server -debug"
