# PropX

PropX is a property-based testing library for Go that allows you to test properties of your code by generating random test cases and automatically shrinking counterexamples when failures are found.

## Features

- **Property-based testing** with automatic test case generation
- **Intelligent shrinking** to find minimal counterexamples
- **Rich generator library** for common data types
- **Custom generators** with user-defined shrinking logic
- **Parallel execution** for faster test runs
- **Command-line configuration** via flags
- **Domain-specific generators** (e.g., CPF validation)
- **State machine testing** for complex stateful systems

## Installation

```bash
go get arcsyn.io/propx
```

## Quick Start

Follow these simple steps to get started with PropX:

### Step 1: Create a test file
Create a new file called `example_test.go` in your project directory.

### Step 2: Add the test code
Copy and paste the following code into your test file:

```go
package main

import (
    "testing"
    "arcsyn.io/propx"
)

func TestAdditionCommutativity(t *testing.T) {
    propx.ForAll(t, propx.Default(),
        propx.Int(propx.Size{Max: 1000}),
        propx.Int(propx.Size{Max: 1000}),
    )(func(t *testing.T, a, b int) {
        if a+b != b+a {
            t.Errorf("addition is not commutative: %d + %d = %d, but %d + %d = %d",
                a, b, a+b, b, a, b+a)
        }
    })
}

```

### Step 3: Run the tests
Execute the tests using the standard Go test command:

```bash
# Run all tests
go test

# Run tests with verbose output to see individual test cases
go test -v

# Run a specific test
go test -run TestAdditionCommutativity

# Run tests with custom PropX configuration
go test -propx.examples=1000 -propx.maxshrink=50
```

### Step 4: Understanding the output
When you run the tests, PropX will:
- Generate random test cases automatically
- Run each test case to verify the property
- If a test fails, it will shrink the input to find a minimal counterexample
- Report the results with clear error messages

That's it! You're now ready to use PropX for property-based testing in your Go projects.

### Reproducing Failed Tests

When a property-based test fails, PropX provides a command to reproduce the exact failure:

```bash
# Example output from a failed test:
[propx] property failed; seed=12345; examples_run=42; shrunk_steps=15
counterexample (min): [1, 2, 3]
replay: go test -run '^TestMyProperty$/ex#l2(/|$)' -propx.seed=12345

# To reproduce the failure:
go test -run '^TestMyProperty$/ex#l2(/|$)' -propx.seed=12345
```

## Command Line Flags

PropX supports several command-line flags for configuring property-based tests:

| Flag | Description | Default |
|------|-------------|---------|
| `-propx.seed` | Random seed for test case generation (0 = random) | 0 |
| `-propx.examples` | Number of test cases to generate | 100 |
| `-propx.maxshrink` | Maximum number of shrinking steps | 400 |
| `-propx.shrink.strategy` | Shrinking strategy: "bfs" or "dfs" | "bfs" |
| `-propx.shrink.subtests` | Use Go's subtest functionality | true |
| `-propx.shrink.parallel` | Number of parallel workers | 1 |

### Usage Examples

```bash
# Run with more test cases
go test -propx.examples=1000

# Use depth-first shrinking strategy
go test -propx.shrink.strategy=dfs

# Run with specific seed for reproducible results
go test -propx.seed=12345

# Use parallel execution with 4 workers
go test -propx.shrink.parallel=4

# Combine multiple flags
go test -propx.examples=500 -propx.maxshrink=200 -propx.shrink.strategy=dfs -propx.shrink.parallel=2
```

## State Machine Testing
- See [State Machine Testing Doc](docs/state-machine.md) - Testing stateful systems

## Shrinking Strategies: BFS vs DFS

PropX supports two different shrinking strategies, each with distinct characteristics:

#### BFS (Breadth-First Search) - Default
- **Behavior**: Explores all candidates at the current "level" before moving deeper
- **Advantages**:
  - Finds counterexamples that are "closer" to the original failing input
  - More predictable shrinking path
  - Better for understanding why a property fails
- **Use when**: You want to understand the minimal change that breaks your property
- **Example**: If your original input was `[1, 2, 3, 4, 5]`, BFS might find `[1, 2, 3, 4]` before trying `[1, 2, 3]`

#### DFS (Depth-First Search)
- **Behavior**: Explores one shrinking path as deeply as possible before trying alternatives
- **Advantages**:
  - Often finds smaller counterexamples faster
  - More aggressive shrinking
  - Better for finding the absolute minimum failing case
- **Use when**: You want the smallest possible counterexample, regardless of how different it is from the original
- **Example**: If your original input was `[1, 2, 3, 4, 5]`, DFS might quickly find `[1]` or `[0]` as the minimal case

#### Choosing the Right Strategy

```bash
# Use BFS for debugging and understanding failures
go test -propx.shrink.strategy=bfs

# Use DFS for finding the smallest possible counterexample
go test -propx.shrink.strategy=dfs
```

> **Recommendation**: Start with BFS (default) for most use cases, then try DFS if you need more aggressive shrinking.

## Examples

See the `examples/` directory for comprehensive usage examples including:
- Basic property testing
- Custom generators
- CPF validation testing
- String and integer property tests
- State machine testing (BankAccount, Counter, Cache)

## Documentation

- [Architecture Decision Records](docs/adrs) - Design decisions and rationale

## License

This project is licensed under the MIT License.

