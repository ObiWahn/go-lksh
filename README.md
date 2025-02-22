# LKSH - Command Execution Utility

LKSH is a Go package that provides a convenient way to execute commands with
enhanced features like output capturing, timeout control, and decorated output.
But otherwise just **like** in a **shell**.

## Features

- Execute commands with context support
- Capture command output
- Set execution timeouts
- Path lookup for executables
- Decorated output mode with colorized formatting
- Exit code handling

## Configuration

The `Config` struct allows you to customize command execution:

- `Decorate`: Enable decorated output with colorized formatting
- `Pipe`: Enable piping command output to the terminal
- `MaxBufferSize`: Set the maximum buffer size for command output. If set to a
   positive value, the command output will be captured into a buffer and then
   printed to the terminal. NOTE: The buffer ist printed to stdout only.
- `KeepEnvVar`: Set the environment variables to keep from the current process.
- `AddEnvVar`: Set the environment variables to add to the command.
- `LookUpPath`: Enable path lookup for the command.
- `DefaultExit`: Set the default exit code.
