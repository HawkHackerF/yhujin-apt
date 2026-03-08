# Yhujin Cloudhost v2.0 — Go Port

Konversi dari `yhujin.py` (Flask/Python) ke Go murni menggunakan `net/http`.

## Dependensi

```bash
go get github.com/creack/pty
go get github.com/shirou/gopsutil/v3
```

## Build & Run

```bash
go mod tidy
go build -o yhujin .
./yhujin
```

Atau langsung:

```bash
go run .
```

Server berjalan di `http://0.0.0.0:5000`

## Perubahan dari versi Python

| Go (net/http)                          |
|----------------------------------------|
| `net/http` (stdlib)                   |
| `github.com/shirou/gopsutil/v3`       |
| `github.com/creack/pty`               |
| Fungsi `naturalSize()` custom         |
| `archive/zip`, `archive/tar` (stdlib) |
| Map `termSessions` dengan mutex       |
| Goroutine                             |
| `filepath.Base()`               |

## Fitur

- File Manager (list, upload, download, rename, delete, hide/unhide, extract, copy, cut)
- Terminal web (PTY multi-tab)
- Process Manager (kill process)
- System Monitor (CPU, RAM, disk, network)
- Server Manager (SSH shortcuts, ping)
- Settings panel
