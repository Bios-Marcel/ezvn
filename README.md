# ezvn

A wrapper around svn to simplify some stuff.

## Installation

### Manually

Running this will install `ezvn.exe` into `~/go/bin`.
```
go install github.com/Bios-Marcel/ezvn@latest
```

Either add `~/go/bin` to your `PATH` environment variable or move the binary
to a different path.

### Scoop

Add the bucket (one time action):

```
scoop bucket add extras
scoop bucket add biosmarcel "https://github.com/Bios-Marcel/scoopbucket.git"
```

Install ezvn:

```
scoop install ezvn
```

## Autocompletion

Powershell:

```
ezvn completion powershell | Out-String | Invoke-Expression
```

## Usage

Use it like you would SVN. Any unknown commands will be delegated.

