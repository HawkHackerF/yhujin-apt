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

Server berjalan di `http://localhost:5000`

## Perubahan dari versi Python

| Python (Flask)         | Go (net/http)                          |
|------------------------|----------------------------------------|
| `flask`                | `net/http` (stdlib)                   |
| `psutil`               | `github.com/shirou/gopsutil/v3`       |
| `pty` (Python)         | `github.com/creack/pty`               |
| `humanize`             | Fungsi `naturalSize()` custom         |
| `zipfile`, `tarfile`   | `archive/zip`, `archive/tar` (stdlib) |
| Session Flask          | Map `termSessions` dengan mutex       |
| `threading.Thread`     | Goroutine                             |
| `werkzeug.secure_filename` | `filepath.Base()`               |

## Fitur

- File Manager (list, upload, download, rename, delete, hide/unhide, extract)
- Terminal web (PTY multi-tab)
- Process Manager (kill process)
- System Monitor (CPU, RAM, disk, network)
- Server Manager (SSH shortcuts, ping)
- Settings panel
