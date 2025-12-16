# 🚀 txt2llm

*Your ultimate text file aggregator for AI conversations*

Stop copying and pasting files one by one! `txt2llm` takes your scattered text files and turns them into a clean, organised feast that AI models absolutely love to consume.

## ✨ What it does

Ever tried feeding multiple text files to ChatGPT? Copy-paste hell, right? Whether it's code, docs, configs, or any readable text files, `txt2llm` solves this by:

- 📁 **Concatenates multiple files** into one clean output
- 🏷️ **Adds clear file markers** so the AI knows what's what
- 🔍 **Supports glob patterns** because wildcards are life
- 📂 **Handles directories recursively** (when you want it to)
- 🎯 **Customizable markers** for different AI/file preferences

## 🏃‍♂️ Quick Start

```bash
# Grab all Go files in your project
txt2llm "**/*.go"

# Mix code, docs, and configs
txt2llm main.go pkg/**/*.go README.md config.yaml

# All text files in a directory
txt2llm --recursive --relative src/

# Documentation files only
txt2llm "**/*.{md,txt,rst}"
```

## 🛠️ Installation

### Pre-built binaries
Check the [releases page](https://github.com/matthewchivers/txt2llm/releases) for ready-to-go binaries.

### Build from source
```bash
git clone https://github.com/matthewchivers/txt2llm
cd txt2llm
make build
```

## 🎛️ Options

| Flag | Description | Default |
|------|-------------|---------|
| `--recursive` | Dive into subdirectories | `false` |
| `--relative` | Use relative paths in output | `false` |
| `--marker-prefix` | Start marker prefix | `<<<` |
| `--marker-suffix` | End marker suffix | `>>>` |

## 💡 Pro Tips

**Perfect for code reviews:**
```bash
txt2llm --relative "src/**/*.{js,ts,tsx}" > review.txt
```

**Gather all documentation directly to clipboard:**
```bash
txt2llm "*.md" "docs/**/*.md" | pbcopy  # macOS
txt2llm "*.md" "docs/**/*.md" | xclip   # Linux
```

**Configuration analysis:**
```bash
txt2llm "*.{json,yaml,toml,ini}" --relative
```

**Research papers and notes:**
```bash
txt2llm "research/**/*.{txt,md}" "notes/*.txt"
```

**Custom markers for specific output types:**
```bash
txt2llm --marker-prefix "```" --marker-suffix "```" *.py
```

## 🎯 Example Output

Command: `txt2llm ./main.go ./utils/helper.go --relative`


```
Each section below represents text output from one file.
Delimiters: <<<START:{filename}>>> ... <<<END:{filename}>>>

<<<START:main.go>>>
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
<<<END:main.go >>>

<<<START:utils/helper.go >>>
package utils

func Helper() string {
    return "I'm helping!"
}
<<<END:utils/helper.go >>>
```

## 🤝 Contributing

Found a bug? Have a feature idea? PRs welcome! This is a simple tool with a simple mission - make it easier to work with AI models.

## 📄 License

MIT - because sharing is caring.

---

*Built with ❤️ for anyone who needs to feed text files to AI models*
