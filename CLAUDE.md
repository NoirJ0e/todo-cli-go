# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture

This is a simple todo CLI application written in Go that manages tasks through JSON file persistence. The application follows a single-file architecture with all core functionality in `main.go`.

### Core Components

- **Task Structure**: The `taskStruct` type defines the data model with ID (UUID), content, creation/completion dates, and completion status
- **File Operations**: Tasks are persisted to `tasks.json` using JSON serialization
- **Command Interface**: The application accepts commands through command-line arguments with a `run()` function that validates against allowed commands: `add`, `remove`, `complete`

### Key Functions

- `createTask()`: Creates new tasks with UUID and timestamps
- `addTask()`: Appends tasks to the task slice
- `saveTasksToFile()` / `loadTasksFromFile()`: Handle JSON persistence
- `run()`: Main command dispatcher (currently validates but doesn't implement command logic)

## Development Commands

### Build

```bash
go build -o todo
```

### Run Tests

```bash
go test -v
```

### Run Single Test

```bash
go test -v -run TestName
```

## Current State

The application is in development with basic data structures and file I/O complete. The command handling in `run()` function validates commands but doesn't implement the actual functionality yet. Test coverage exists for all implemented functions.

## Dependencies

- `github.com/google/uuid` - UUID generation for task IDs
- Standard Go library for JSON, file I/O, and time handling

## Project Purpose

- Taught me how to write go like a friendly mentor while maintain the TDD behavior: guide me how to create tests, and
  then how to write the code. NEVER show me the implementaion of code unless I ask you to do so
