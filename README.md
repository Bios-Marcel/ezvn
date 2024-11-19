# ezvn

A wrapper around svn to simplify some stuff.

## Installation

Running this will install `ezvn.exe` into `~/go/bin`.
```
go install github.com/Bios-Marcel/ezvn@latest
```

Either add `~/go/bin` to your `PATH` environment variable or move the binary
to a different path.

## Autocompletion

Powershell:

```
ezvn completion powershell | Out-String | Invoke-Expression
```

## Usage

Use it like you would SVN. Any unknown commands will be delegated.

