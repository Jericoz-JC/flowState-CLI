# flowstate-cli

A unified terminal productivity system for notes, todos, and focus sessions with semantic search.

## Installation

```bash
npm install -g flowstate-cli
```

## Usage

```bash
flowstate
```

If `flowstate` command is not found, try:
```bash
npx flowstate
```

## Supported Platforms

| Platform | Architecture | Status |
|----------|--------------|--------|
| Windows | x64 | Supported |
| Windows | ARM64 | Supported |
| macOS | Intel (x64) | Supported |
| macOS | Apple Silicon (ARM64) | Supported |
| Linux | x64 | Supported |
| Linux | ARM64 | Supported |

> **Note**: 32-bit systems are not supported. Requires 64-bit Node.js.

## Features

- **Notes**: Quick capture with markdown preview, wikilinks `[[Note Title]]`, and `#hashtag`/`@mention` tagging
- **Todos**: Task management with priorities, due dates, status badges, and multiple sort/filter modes
- **Focus Sessions**: Pomodoro-style timer with configurable durations and session history
- **Mind Map**: Visual graph of note connections
- **Semantic Search**: Local ONNX-powered semantic search
- **Linking System**: Connect notes and todos through bidirectional relationships

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Ctrl+X` | Quick capture note |
| `Ctrl+N` | Notes screen |
| `Ctrl+T` | Todos screen |
| `Ctrl+F` | Focus session screen |
| `Ctrl+/` | Semantic search |
| `Ctrl+G` | Mind map |
| `?` | Help modal |
| `q` | Quit |

## Troubleshooting

### "ia32" or "x86" architecture error
You have 32-bit Node.js installed. Please install [64-bit Node.js](https://nodejs.org/).

### "flowstate: command not found"
npm global bin is not in your PATH. Options:
1. Use `npx flowstate` (always works)
2. Add to PATH:
   ```bash
   # Add to ~/.bashrc or ~/.zshrc:
   export PATH="$(npm config get prefix)/bin:$PATH"
   ```

### Old version running
Reinstall to fix PATH conflicts:
```bash
npm uninstall -g flowstate-cli && npm install -g flowstate-cli
```

## More Information

For full documentation and source code, visit [GitHub](https://github.com/Jericoz-JC/flowState-CLI).

## License

MIT
