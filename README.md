## Groblin
> A smart, scalable Go crawler for e-comm

## Features
- Multi file suppport
- Concurrent Crawling
- Smart URL Discovery
- Color logs
- Fast
- Caching
- Scalable
- Performance Optimized

## Workflow
```mermaid
flowchart TD
    A[Start] --> B[Parse Input Arguments]
    B --> C[Load Domain List]
    C --> D[Load User-Agent Config]
    D --> E[Load Category Config]
    E --> F[Initialize Colly Crawler]
    F --> G[For Each Domain in Domain List]
    G --> H[Crawl Domain Using Goroutines]
    H --> I[Regex Match Product Links]
    I --> J[Store Results in Shared Map]
    J --> K[Wait for All Goroutines to Complete]
    K --> L[Save Results to JSON File]
    L --> M[End]

    B --> C
    B --> F
    G --> H
    G --> I
    I --> J
    H --> J
    K --> L
    B -.->|File Input| C
    B -.->|Concurrency Limit| F
```

<!--- ![Untitled diagram-2025-01-22-114427](https://github.com/user-attachments/assets/e7cfd74c-eb94-4291-82f1-63ab437d538c) --->

---

## Requirements

- Go (version 1.22 or higher)
- [Colly](https://github.com/gocolly/colly)
- [Charmbracelet/log](https://github.com/charmbracelet/log) 
- Configuration files:
  - `config/user_agent.json` - Contains a list of User-Agents.
  - `config/category.json` - Contains the categories for filtering product links.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/r3yc0n1c/groblin.git
   cd groblin
   ```
2. Install dependencies:
  ```bash
  go mod tidy
  ```
3. Build
  ```bash
  make build
  ```

## Usage

Command-Line Arguments:
- `--file`: Path to the input file containing the list of domains (CSV/JSON).
- `--n`: Number of domains to explore concurrently (default: 1).

### Start Crawling

```bash
bin/groblin --file ./dom.json --n 4
```

### Output
```json
{
  "example.com": [
    "https://example.com/items/1",
    "https://example.com/items/2"
  ],
  "anotherexample.com": [
    "https://anotherexample.com/products/1"
  ]
}
```
