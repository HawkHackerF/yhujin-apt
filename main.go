package main

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// ── Config ───────────────────────────────────────────────────────────────────

var (
	BASE_DIR    string
	STORAGE_DIR string
	AUTH_TOKEN  string
	DEV_MODE    bool
	BUILD_HASH  string
)

const (
	MAX_UPLOAD_SIZE = 1024 * 1024 * 1024 // 1 GB
	IDLE_TIMEOUT    = 30 * time.Minute
)

func init() {
	exe, err := os.Executable()
	if err != nil {
		exe = "."
	}
	BASE_DIR = filepath.Dir(exe)
	STORAGE_DIR = filepath.Join(BASE_DIR, "storage")
	os.MkdirAll(STORAGE_DIR, 0755)

	DEV_MODE = os.Getenv("YHUJIN_DEV") == "true" || os.Getenv("YHUJIN_DEV") == "1"

	if DEV_MODE {
		b := make([]byte, 6)
		rand.Read(b)
		BUILD_HASH = fmt.Sprintf("%x", b)
	} else {
		BUILD_HASH = "v2"
	}

	AUTH_TOKEN = os.Getenv("YHUJIN_TOKEN")
	if AUTH_TOKEN == "" {
		b := make([]byte, 24)
		rand.Read(b)
		AUTH_TOKEN = base64.URLEncoding.EncodeToString(b)
	}
}

// ── Session ───────────────────────────────────────────────────────────────────

var (
	sessionsMu sync.Mutex
	sessions   = map[string]time.Time{}
)

func newSession() string {
	b := make([]byte, 24)
	rand.Read(b)
	sid := base64.URLEncoding.EncodeToString(b)
	sessionsMu.Lock()
	sessions[sid] = time.Now().Add(24 * time.Hour)
	sessionsMu.Unlock()
	return sid
}

func checkSession(r *http.Request) bool {
	c, err := r.Cookie("yhujin_sid")
	if err != nil {
		return false
	}
	sessionsMu.Lock()
	exp, ok := sessions[c.Value]
	if ok {
		sessions[c.Value] = time.Now().Add(24 * time.Hour) // refresh
	}
	sessionsMu.Unlock()
	return ok && time.Now().Before(exp)
}

func cleanSessions() {
	for {
		time.Sleep(10 * time.Minute)
		sessionsMu.Lock()
		now := time.Now()
		for sid, exp := range sessions {
			if now.After(exp) {
				delete(sessions, sid)
			}
		}
		sessionsMu.Unlock()
	}
}

// ── Middleware ────────────────────────────────────────────────────────────────

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dev mode: no cache
		if DEV_MODE {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			w.Header().Set("Pragma", "no-cache")
		}

		// Login/logout selalu boleh
		if r.URL.Path == "/login" || r.URL.Path == "/logout" {
			next.ServeHTTP(w, r)
			return
		}

		// Cek session
		if !checkSession(r) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(401)
				w.Write([]byte(`{"error":"Unauthorized"}`))
			} else {
				http.Redirect(w, r, "/login", http.StatusFound)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ── Terminal Idle Cleanup ─────────────────────────────────────────────────────

func cleanIdleTerminals() {
	for {
		time.Sleep(5 * time.Minute)
		termSessionsMu.Lock()
		now := time.Now()
		for sid, ts := range termSessions {
			ts.mu.Lock()
			idle := now.Sub(ts.lastActivity)
			ts.mu.Unlock()
			if idle > IDLE_TIMEOUT {
				ts.Close()
				delete(termSessions, sid)
				log.Printf("[TERMINAL] Session %s closed (idle %.0fm)", sid[:8], idle.Minutes())
			}
		}
		termSessionsMu.Unlock()
	}
}

// ── Login / Logout ────────────────────────────────────────────────────────────

const loginHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Yhujin — Login</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{background:#18181d;color:#dde1e8;font-family:'JetBrains Mono',monospace;height:100vh;display:flex;align-items:center;justify-content:center}
.box{background:#202027;border:1px solid #38383f;border-radius:12px;padding:40px;width:340px}
.brand{color:#4a9eff;font-size:1.1rem;font-weight:700;letter-spacing:2px;margin-bottom:28px;text-align:center}
label{font-size:.75rem;color:#8e95a3;display:block;margin-bottom:6px}
input{width:100%;background:#28282f;border:1px solid #38383f;color:#dde1e8;font-family:inherit;font-size:.85rem;padding:9px 13px;border-radius:6px;outline:none;margin-bottom:18px}
input:focus{border-color:#4a9eff}
button{width:100%;background:#4a9eff;border:none;color:#000;font-family:inherit;font-size:.85rem;font-weight:600;padding:10px;border-radius:6px;cursor:pointer;transition:background .15s}
button:hover{background:#5aa8ff}
.err{color:#f85149;font-size:.78rem;margin-bottom:14px;display:none}
</style>
</head>
<body>
<div class="box">
  <div class="brand">YHUJIN</div>
  <div class="err" id="err">Token salah, coba lagi.</div>
  <label>Access Token</label>
  <input type="password" id="tok" placeholder="••••••••" autofocus onkeydown="if(event.key==='Enter')doLogin()">
  <button onclick="doLogin()">Login</button>
</div>
<script>
function doLogin(){
  const tok=document.getElementById('tok').value;
  fetch('/login',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({token:tok})})
    .then(r=>r.json()).then(d=>{
      if(d.ok){ window.location.href='/'; }
      else{ document.getElementById('err').style.display='block'; }
    });
}
</script>
</body>
</html>`

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(loginHTML))
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}
	var body struct {
		Token string `json:"token"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if subtle.ConstantTimeCompare([]byte(body.Token), []byte(AUTH_TOKEN)) != 1 {
		time.Sleep(500 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":false}`))
		return
	}
	sid := newSession()
	http.SetCookie(w, &http.Cookie{
		Name:     "yhujin_sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("yhujin_sid")
	if err == nil {
		sessionsMu.Lock()
		delete(sessions, c.Value)
		sessionsMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "yhujin_sid",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}



// ── Terminal Session ─────────────────────────────────────────────────────────

type TerminalSession struct {
	SessionID      string
	ptmx           *os.File
	running        bool
	outputBuffer   []byte
	scrollback     []byte // persistent scrollback for reconnect
	lastActivity   time.Time
	mu             sync.Mutex
}

var (
	termSessions   = map[string]*TerminalSession{}
	termSessionsMu sync.Mutex
)

func newTerminalSession(id string) *TerminalSession {
	return &TerminalSession{
		SessionID:    id,
		running:      false,
		lastActivity: time.Now(),
	}
}

func (ts *TerminalSession) Start(cols, rows uint16) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	cmd := exec.Command(shell, "--login")
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
	)
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: rows, Cols: cols})
	if err != nil {
		return err
	}
	ts.ptmx = ptmx
	ts.running = true
	go ts.readOutput()
	return nil
}

func (ts *TerminalSession) readOutput() {
	buf := make([]byte, 4096)
	for {
		n, err := ts.ptmx.Read(buf)
		if n > 0 {
			ts.mu.Lock()
			ts.outputBuffer = append(ts.outputBuffer, buf[:n]...)
			// Keep last 512KB in scrollback for reconnect
			ts.scrollback = append(ts.scrollback, buf[:n]...)
			const maxScrollback = 512 * 1024
			if len(ts.scrollback) > maxScrollback {
				ts.scrollback = ts.scrollback[len(ts.scrollback)-maxScrollback:]
			}
			ts.lastActivity = time.Now()
			ts.mu.Unlock()
		}
		if err != nil {
			break
		}
	}
	ts.mu.Lock()
	ts.running = false
	ts.mu.Unlock()
}

func (ts *TerminalSession) ReadOutput() []byte {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	data := ts.outputBuffer
	ts.outputBuffer = nil
	return data
}

func (ts *TerminalSession) Write(data []byte) bool {
	ts.mu.Lock()
	running := ts.running
	ts.mu.Unlock()
	if !running || ts.ptmx == nil {
		return false
	}
	_, err := ts.ptmx.Write(data)
	if err != nil {
		ts.mu.Lock()
		ts.running = false
		ts.mu.Unlock()
		return false
	}
	ts.lastActivity = time.Now()
	return true
}

func (ts *TerminalSession) Resize(cols, rows uint16) {
	if ts.ptmx != nil {
		pty.Setsize(ts.ptmx, &pty.Winsize{Rows: rows, Cols: cols})
	}
}

func (ts *TerminalSession) Close() {
	ts.mu.Lock()
	ts.running = false
	ts.mu.Unlock()
	if ts.ptmx != nil {
		ts.ptmx.Close()
	}
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func naturalSize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), suffixes[exp])
}

func fmtUptime(secs float64) string {
	d := int(secs / 86400)
	h := int(secs/3600) % 24
	m := int(secs/60) % 60
	if d > 0 {
		return fmt.Sprintf("%dd %dh %dm", d, h, m)
	}
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func jsonOK(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func securePath(base, rel string) (string, bool) {
	// Kalau rel adalah absolute path (mulai dari /), gunakan langsung
	// tapi tetap clean untuk mencegah path traversal
	var full string
	if filepath.IsAbs(rel) {
		full = filepath.Clean(rel)
	} else {
		rel = strings.TrimPrefix(rel, "/")
		full = filepath.Join(base, rel)
	}
	// Pastikan tidak ada null bytes atau karakter berbahaya
	if strings.ContainsAny(full, "\x00") {
		return "", false
	}
	return full, true
}

// securePathUnder memastikan path berada di bawah base directory (strict)
func securePathUnder(base, rel string) (string, bool) {
	rel = strings.TrimPrefix(rel, "/")
	full := filepath.Join(base, rel)
	if !strings.HasPrefix(filepath.Clean(full)+"/", filepath.Clean(base)+"/") {
		return "", false
	}
	return full, true
}

// ── System Info ───────────────────────────────────────────────────────────────

func getSystemInfo() map[string]interface{} {
	result := map[string]interface{}{}

	// CPU
	cpuPct, _ := cpu.Percent(100*time.Millisecond, false)
	cpuPerCoreRaw, _ := cpu.Percent(100*time.Millisecond, true)
	cpuPerCore := make([]float64, len(cpuPerCoreRaw))
	for i, v := range cpuPerCoreRaw {
		cpuPerCore[i] = math.Round(v*10) / 10
	}
	cpuInfo, _ := cpu.Info()
	freqs, _ := cpu.Percent(0, false)
	_ = freqs
	logicalCores, _ := cpu.Counts(true)
	physCores, _ := cpu.Counts(false)
	processor := runtime.GOARCH
	if len(cpuInfo) > 0 {
		processor = cpuInfo[0].ModelName
	}
	var cpuFreq float64
	if len(cpuInfo) > 0 {
		cpuFreq = cpuInfo[0].Mhz
	}
	var cpuPercent float64
	if len(cpuPct) > 0 {
		cpuPercent = math.Round(cpuPct[0]*10) / 10
	}

	// Memory
	vmem, _ := mem.VirtualMemory()
	swap, _ := mem.SwapMemory()

	// Disk
	diskStat, _ := disk.Usage(BASE_DIR)

	// Network
	netIO, _ := net.IOCounters(false)
	var bytesSent, bytesRecv uint64
	if len(netIO) > 0 {
		bytesSent = netIO[0].BytesSent
		bytesRecv = netIO[0].BytesRecv
	}

	// Net IPs
	interfaces, _ := net.Interfaces()
	var ips []string
	for _, iface := range interfaces {
		for _, addr := range iface.Addrs {
			ip := strings.Split(addr.Addr, "/")[0]
			if !strings.HasPrefix(ip, "127.") && strings.Contains(ip, ".") {
				ips = append(ips, ip)
			}
		}
	}

	// Uptime
	bootTime, _ := host.BootTime()
	uptime := float64(time.Now().Unix()) - float64(bootTime)
	hostname, _ := os.Hostname()

	// Load avg
	loadAvg := []float64{0, 0, 0}
	if la, err := load.Avg(); err == nil {
		loadAvg = []float64{la.Load1, la.Load5, la.Load15}
	}

	// Users
	users, _ := host.Users()

	var diskPct float64
	if diskStat != nil {
		diskPct = math.Round(diskStat.UsedPercent*10) / 10
	}

	result["hostname"] = hostname
	result["platform"] = runtime.GOOS
	result["platform_release"] = ""
	result["architecture"] = runtime.GOARCH
	result["processor"] = processor
	result["cpu_cores"] = physCores
	result["cpu_threads"] = logicalCores
	result["cpu_percent"] = cpuPercent
	result["cpu_per_core"] = cpuPerCore
	result["cpu_freq_current"] = math.Round(cpuFreq)
	result["ip_addresses"] = ips
	result["net_bytes_sent"] = naturalSize(bytesSent)
	result["net_bytes_recv"] = naturalSize(bytesRecv)
	result["uptime"] = uptime
	result["boot_time"] = time.Unix(int64(bootTime), 0).Format(time.RFC3339)
	result["users"] = len(users)
	result["load_avg"] = loadAvg
	if vmem != nil {
		result["memory"] = map[string]interface{}{
			"total":     naturalSize(vmem.Total),
			"available": naturalSize(vmem.Available),
			"used":      naturalSize(vmem.Used),
			"percent":   math.Round(vmem.UsedPercent*10) / 10,
		}
	}
	if swap != nil {
		result["swap"] = map[string]interface{}{
			"total":   naturalSize(swap.Total),
			"used":    naturalSize(swap.Used),
			"percent": math.Round(swap.UsedPercent*10) / 10,
		}
	}
	if diskStat != nil {
		result["disk"] = map[string]interface{}{
			"total":       naturalSize(diskStat.Total),
			"used":        naturalSize(diskStat.Used),
			"free":        naturalSize(diskStat.Free),
			"percent":     diskPct,
			"mount_point": BASE_DIR,
		}
	}
	return result
}

// ── Processes ─────────────────────────────────────────────────────────────────

func getProcesses() []map[string]interface{} {
	procs, err := process.Processes()
	if err != nil {
		return nil
	}
	var result []map[string]interface{}
	for _, p := range procs {
		pid := p.Pid
		name, _ := p.Name()
		username, _ := p.Username()
		status, _ := p.Status()
		cpuPct, _ := p.CPUPercent()
		memPct, _ := p.MemoryPercent()
		memInfo, _ := p.MemoryInfo()
		createTime, _ := p.CreateTime()
		cmdline, _ := p.CmdlineSlice()

		rss := uint64(0)
		if memInfo != nil {
			rss = memInfo.RSS
		}
		cmdStr := name
		if len(cmdline) > 0 {
			end := len(cmdline)
			if end > 3 {
				end = 3
			}
			cmdStr = strings.Join(cmdline[:end], " ")
		}
		started := "-"
		if createTime > 0 {
			started = time.Unix(createTime/1000, 0).Format("15:04:05")
		}
		statusStr := "-"
		if len(status) > 0 {
			statusStr = status[0]
		}
		result = append(result, map[string]interface{}{
			"pid":            pid,
			"name":           name,
			"username":       username,
			"status":         statusStr,
			"cpu_percent":    math.Round(cpuPct*10) / 10,
			"memory_percent": math.Round(float64(memPct)*10) / 10,
			"memory_rss":     naturalSize(rss),
			"started":        started,
			"cmdline":        cmdStr,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		a := result[i]["cpu_percent"].(float64)
		b := result[j]["cpu_percent"].(float64)
		return a > b
	})
	if len(result) > 100 {
		result = result[:100]
	}
	return result
}

// ── File Helpers ──────────────────────────────────────────────────────────────

var archiveExts = map[string]bool{
	".zip": true, ".tar": true, ".gz": true, ".bz2": true,
	".xz": true, ".7z": true, ".rar": true, ".tgz": true,
	".tbz2": true, ".txz": true,
}

func isArchive(name string) bool {
	lower := strings.ToLower(name)
	for _, suffix := range []string{".tar.gz", ".tar.bz2", ".tar.xz"} {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	ext := strings.ToLower(filepath.Ext(name))
	return archiveExts[ext]
}

func getFilePermissions(path string) string {
	info, err := os.Lstat(path)
	if err != nil {
		return "----------"
	}
	mode := info.Mode()
	perms := ""
	if mode.IsDir() {
		perms = "d"
	} else if mode&os.ModeSymlink != 0 {
		perms = "l"
	} else {
		perms = "-"
	}
	bits := []struct {
		mask os.FileMode
		char string
	}{
		{0400, "r"}, {0200, "w"}, {0100, "x"},
		{0040, "r"}, {0020, "w"}, {0010, "x"},
		{0004, "r"}, {0002, "w"}, {0001, "x"},
	}
	for _, b := range bits {
		if mode&b.mask != 0 {
			perms += b.char
		} else {
			perms += "-"
		}
	}
	return perms
}

func getOwnerGroup(path string) (string, string) {
	info, err := os.Lstat(path)
	if err != nil {
		return "-", "-"
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "-", "-"
	}
	owner := strconv.Itoa(int(stat.Uid))
	group := strconv.Itoa(int(stat.Gid))
	if u, err := user.LookupId(owner); err == nil {
		owner = u.Username
	}
	if g, err := user.LookupGroupId(group); err == nil {
		group = g.Name
	}
	return owner, group
}

func getDirectoryListing(dirPath string) ([]map[string]interface{}, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var items []map[string]interface{}
	for _, entry := range entries {
		name := entry.Name()
		itemPath := filepath.Join(dirPath, name)
		info, err := os.Lstat(itemPath)
		if err != nil {
			continue
		}
		ext := strings.ToLower(filepath.Ext(name))
		isDir := info.IsDir()
		isLink := info.Mode()&os.ModeSymlink != 0
		mimeType := guessMime(ext)
		perms := getFilePermissions(itemPath)
		owner, group := getOwnerGroup(itemPath)
		items = append(items, map[string]interface{}{
			"name":           name,
			"is_dir":         isDir,
			"is_file":        !isDir,
			"is_link":        isLink,
			"is_hidden":      strings.HasPrefix(name, "."),
			"is_archive":     isArchive(name),
			"size":           info.Size(),
			"size_formatted": naturalSize(uint64(info.Size())),
			"modified":       info.ModTime().Format(time.RFC3339),
			"permissions":    perms,
			"owner":          owner,
			"group":          group,
			"extension":      ext,
			"mime_type":      mimeType,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		ai := items[i]["is_dir"].(bool)
		aj := items[j]["is_dir"].(bool)
		if ai != aj {
			return ai
		}
		return strings.ToLower(items[i]["name"].(string)) < strings.ToLower(items[j]["name"].(string))
	})
	return items, nil
}

func guessMime(ext string) string {
	mimes := map[string]string{
		".html": "text/html", ".css": "text/css", ".js": "application/javascript",
		".json": "application/json", ".png": "image/png", ".jpg": "image/jpeg",
		".jpeg": "image/jpeg", ".gif": "image/gif", ".webp": "image/webp",
		".svg": "image/svg+xml", ".pdf": "application/pdf",
		".mp4": "video/mp4", ".webm": "video/webm",
		".mp3": "audio/mpeg", ".wav": "audio/wav", ".ogg": "audio/ogg",
		".txt": "text/plain", ".md": "text/markdown",
		".zip": "application/zip", ".tar": "application/x-tar",
		".gz": "application/gzip",
	}
	if m, ok := mimes[ext]; ok {
		return m
	}
	return ""
}

func extractArchive(src, dest string) error {
	lower := strings.ToLower(src)
	os.MkdirAll(dest, 0755)

	isTar := strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz") ||
		strings.HasSuffix(lower, ".tar.bz2") || strings.HasSuffix(lower, ".tbz2") ||
		strings.HasSuffix(lower, ".tar.xz") || strings.HasSuffix(lower, ".txz") ||
		strings.HasSuffix(lower, ".tar")

	if strings.HasSuffix(lower, ".zip") {
		return extractZip(src, dest)
	}
	if isTar {
		return extractTar(src, dest)
	}
	if strings.HasSuffix(lower, ".gz") {
		return extractGzip(src, dest)
	}
	if strings.HasSuffix(lower, ".bz2") {
		return extractBzip2(src, dest)
	}
	// fallback to system commands
	if strings.HasSuffix(lower, ".7z") {
		return exec.Command("7z", "x", src, "-o"+dest, "-y").Run()
	}
	if strings.HasSuffix(lower, ".rar") {
		return exec.Command("unrar", "x", src, dest).Run()
	}
	return exec.Command("tar", "xf", src, "-C", dest).Run()
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		outPath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(outPath, filepath.Clean(dest)) {
			continue
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(outPath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(outPath), 0755)
		rc, err := f.Open()
		if err != nil {
			continue
		}
		out, err := os.Create(outPath)
		if err != nil {
			rc.Close()
			continue
		}
		io.Copy(out, rc)
		out.Close()
		rc.Close()
	}
	return nil
}

func extractTar(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	lower := strings.ToLower(src)
	var tr *tar.Reader
	if strings.HasSuffix(lower, ".gz") || strings.HasSuffix(lower, ".tgz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer gr.Close()
		tr = tar.NewReader(gr)
	} else if strings.HasSuffix(lower, ".bz2") || strings.HasSuffix(lower, ".tbz2") {
		tr = tar.NewReader(bzip2.NewReader(f))
	} else {
		tr = tar.NewReader(f)
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		outPath := filepath.Join(dest, hdr.Name)
		if !strings.HasPrefix(outPath, filepath.Clean(dest)) {
			continue
		}
		if hdr.Typeflag == tar.TypeDir {
			os.MkdirAll(outPath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(outPath), 0755)
		out, err := os.Create(outPath)
		if err != nil {
			continue
		}
		io.Copy(out, tr)
		out.Close()
	}
	return nil
}

func extractGzip(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()
	base := filepath.Base(src)
	if strings.HasSuffix(base, ".gz") {
		base = base[:len(base)-3]
	}
	out, err := os.Create(filepath.Join(dest, base))
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, gr)
	return nil
}

func extractBzip2(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	base := filepath.Base(src)
	if strings.HasSuffix(base, ".bz2") {
		base = base[:len(base)-4]
	}
	out, err := os.Create(filepath.Join(dest, base))
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, bzip2.NewReader(f))
	return nil
}

func toggleHidden(fullPath string) (string, error) {
	name := filepath.Base(fullPath)
	parent := filepath.Dir(fullPath)
	var newName string
	if strings.HasPrefix(name, ".") {
		newName = name[1:]
	} else {
		newName = "." + name
	}
	return newName, os.Rename(fullPath, filepath.Join(parent, newName))
}

// ── HTTP Handlers ─────────────────────────────────────────────────────────────

func handleIndex(w http.ResponseWriter, r *http.Request) {
	html := strings.ReplaceAll(htmlTemplate, "__BASEDIR__", BASE_DIR)
	html = strings.ReplaceAll(html, "__STORAGE__", STORAGE_DIR)
	html = strings.ReplaceAll(html, "__BUILD_HASH__", BUILD_HASH)

	devFlag := "false"
	if DEV_MODE {
		devFlag = "true"
	}
	html = strings.ReplaceAll(html, "__DEV_MODE__", devFlag)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, getSystemInfo())
}

func handleProcesses(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, getProcesses())
}

func handleKillProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	pidStr := strings.TrimPrefix(r.URL.Path, "/api/process/kill/")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		jsonErr(w, "Invalid PID", 400)
		return
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		jsonErr(w, fmt.Sprintf("Process %d not found", pid), 404)
		return
	}
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	jsonOK(w, map[string]string{"message": fmt.Sprintf("Process %d terminated", pid)})
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		jsonErr(w, "Bad request", 400)
		return
	}
	host := data["host"]
	if host == "" {
		jsonErr(w, "Host required", 400)
		return
	}
	start := time.Now()
	cmd := exec.Command("ping", "-c", "1", "-W", "2", host)
	err := cmd.Run()
	latency := int(time.Since(start).Milliseconds())
	reachable := err == nil
	jsonOK(w, map[string]interface{}{
		"reachable": reachable,
		"latency":   latency,
		"host":      host,
	})
}

func handleListFiles(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/files")
	if path == "" || path == "/" {
		path = "/"
	}
	fullPath, ok := securePath(BASE_DIR, path)
	if !ok {
		jsonErr(w, "Invalid path", 400)
		return
	}
	if !fileExists(fullPath) {
		jsonErr(w, "Path not found", 404)
		return
	}
	info, _ := os.Stat(fullPath)
	if info != nil && !info.IsDir() {
		http.ServeFile(w, r, fullPath)
		return
	}
	items, err := getDirectoryListing(fullPath)
	if err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	jsonOK(w, items)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(32 * 1024 * 1024); err != nil {
		jsonErr(w, "Upload too large or invalid", 413)
		return
	}
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		jsonErr(w, "No files provided", 400)
		return
	}
	uploadPath := r.FormValue("path")
	if uploadPath == "" {
		uploadPath = "/"
	}
	fullPath, ok := securePath(BASE_DIR, uploadPath)
	if !ok {
		jsonErr(w, "Invalid path", 400)
		return
	}
	var results []map[string]interface{}
	for _, fh := range files {
		filename := filepath.Base(fh.Filename)
		src, err := fh.Open()
		if err != nil {
			continue
		}
		dst, err := os.Create(filepath.Join(fullPath, filename))
		if err != nil {
			src.Close()
			continue
		}
		io.Copy(dst, src)
		dst.Close()
		src.Close()
		results = append(results, map[string]interface{}{"filename": filename, "success": true})
	}
	jsonOK(w, map[string]interface{}{
		"files":   results,
		"message": fmt.Sprintf("%d file(s) uploaded", len(results)),
	})
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/download")
	fullPath, ok := securePath(BASE_DIR, path)
	if !ok || !fileExists(fullPath) {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(fullPath))
	http.ServeFile(w, r, fullPath)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/delete")
	fullPath, ok := securePath(BASE_DIR, path)
	if !ok || !fileExists(fullPath) {
		jsonErr(w, "File not found", 404)
		return
	}
	info, _ := os.Lstat(fullPath)
	var err error
	if info != nil && info.IsDir() {
		err = os.RemoveAll(fullPath)
	} else {
		err = os.Remove(fullPath)
	}
	if err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	jsonOK(w, map[string]string{"message": "Deleted"})
}

func handleRename(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	oldName := data["old_name"]
	newName := data["new_name"]
	dirPath := data["path"]
	if oldName == "" || newName == "" {
		jsonErr(w, "Names required", 400)
		return
	}
	base, ok := securePath(BASE_DIR, dirPath)
	if !ok {
		jsonErr(w, "Invalid path", 400)
		return
	}
	err := os.Rename(filepath.Join(base, oldName), filepath.Join(base, newName))
	if err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	jsonOK(w, map[string]string{"message": "Renamed"})
}

// handleMove — pindahkan satu atau banyak file/folder ke dest
func handleMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	var data struct {
		Src  []string `json:"src"`  // absolute paths
		Dest string   `json:"dest"` // absolute dest dir
	}
	json.NewDecoder(r.Body).Decode(&data)
	if len(data.Src) == 0 || data.Dest == "" {
		jsonErr(w, "src dan dest diperlukan", 400)
		return
	}
	destDir, ok := securePath(BASE_DIR, data.Dest)
	if !ok {
		jsonErr(w, "Dest path invalid", 400)
		return
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	for _, s := range data.Src {
		src, ok := securePath(BASE_DIR, s)
		if !ok {
			jsonErr(w, "Src path invalid: "+s, 400)
			return
		}
		dst := filepath.Join(destDir, filepath.Base(src))
		if src == dst {
			continue
		}
		// Coba rename dulu (same filesystem, cepat)
		if err := os.Rename(src, dst); err != nil {
			// Cross-device: copy lalu delete
			if err2 := copyPath(src, dst); err2 != nil {
				jsonErr(w, err2.Error(), 500)
				return
			}
			os.RemoveAll(src)
		}
	}
	jsonOK(w, map[string]string{"message": "Dipindahkan"})
}

// handleCopy — salin satu atau banyak file/folder ke dest
func handleCopy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	var data struct {
		Src  []string `json:"src"`
		Dest string   `json:"dest"`
	}
	json.NewDecoder(r.Body).Decode(&data)
	if len(data.Src) == 0 || data.Dest == "" {
		jsonErr(w, "src dan dest diperlukan", 400)
		return
	}
	destDir, ok := securePath(BASE_DIR, data.Dest)
	if !ok {
		jsonErr(w, "Dest path invalid", 400)
		return
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	for _, s := range data.Src {
		src, ok := securePath(BASE_DIR, s)
		if !ok {
			jsonErr(w, "Src path invalid: "+s, 400)
			return
		}
		dst := filepath.Join(destDir, filepath.Base(src))
		// Kalau nama bentrok, tambah _copy
		if fileExists(dst) {
			ext := filepath.Ext(dst)
			base := strings.TrimSuffix(dst, ext)
			dst = base + "_copy" + ext
		}
		if err := copyPath(src, dst); err != nil {
			jsonErr(w, err.Error(), 500)
			return
		}
	}
	jsonOK(w, map[string]string{"message": "Disalin"})
}

// copyPath — rekursif copy file atau folder
func copyPath(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		if err := os.MkdirAll(dst, info.Mode()); err != nil {
			return err
		}
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, e := range entries {
			if err := copyPath(filepath.Join(src, e.Name()), filepath.Join(dst, e.Name())); err != nil {
				return err
			}
		}
		return nil
	}
	// File biasa
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func handleCreateFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	name := data["name"]
	if name == "" {
		jsonErr(w, "Name required", 400)
		return
	}
	dirPath, _ := securePath(BASE_DIR, data["path"])
	f, err := os.Create(filepath.Join(dirPath, name))
	if err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	f.Close()
	jsonOK(w, map[string]string{"message": "File created"})
}

func handleCreateFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	name := data["name"]
	if name == "" {
		jsonErr(w, "Name required", 400)
		return
	}
	dirPath, _ := securePath(BASE_DIR, data["path"])
	err := os.MkdirAll(filepath.Join(dirPath, name), 0755)
	if err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	jsonOK(w, map[string]string{"message": "Folder created"})
}

func handleContent(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/content")
	fullPath, ok := securePath(BASE_DIR, path)
	if !ok || !fileExists(fullPath) {
		jsonErr(w, "File not found", 404)
		return
	}
	info, _ := os.Stat(fullPath)
	if info != nil && info.IsDir() {
		jsonErr(w, "Not a file", 400)
		return
	}
	size := info.Size()
	if size > 5*1024*1024 {
		jsonErr(w, fmt.Sprintf("File too large to preview (%s). Download to view.", naturalSize(uint64(size))), 413)
		return
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	content := string(data)
	lines := strings.Count(content, "\n") + 1
	jsonOK(w, map[string]interface{}{
		"content":  content,
		"lines":    lines,
		"size":     naturalSize(uint64(size)),
		"encoding": "utf-8",
	})
}

func handleServe(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/serve")
	fullPath, ok := securePath(BASE_DIR, path)
	if !ok || !fileExists(fullPath) {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, fullPath)
}

func handleHide(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	base, ok := securePath(BASE_DIR, data["path"])
	if !ok {
		jsonErr(w, "Invalid path", 400)
		return
	}
	fullPath := filepath.Join(base, data["name"])
	if !fileExists(fullPath) {
		jsonErr(w, "Not found", 404)
		return
	}
	newName, err := toggleHidden(fullPath)
	if err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	action := "hidden"
	if !strings.HasPrefix(newName, ".") {
		action = "unhidden"
	}
	jsonOK(w, map[string]string{
		"message":  fmt.Sprintf("%q %s", newName, action),
		"new_name": newName,
	})
}

func handleExtract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	base, ok := securePath(BASE_DIR, data["path"])
	if !ok {
		jsonErr(w, "Invalid path", 400)
		return
	}
	src := filepath.Join(base, data["name"])
	if !fileExists(src) {
		jsonErr(w, "Archive not found", 404)
		return
	}
	rawDest := data["dest"]
	var dest string
	if rawDest != "" {
		dest = filepath.Join(base, rawDest)
	} else {
		stem := data["name"]
		for _, ext := range []string{".tar.gz", ".tar.bz2", ".tar.xz", ".tgz", ".tbz2", ".txz", ".zip", ".gz", ".bz2", ".xz", ".7z", ".rar", ".tar"} {
			if strings.HasSuffix(strings.ToLower(stem), ext) {
				stem = stem[:len(stem)-len(ext)]
				break
			}
		}
		dest = filepath.Join(base, stem)
	}
	os.MkdirAll(dest, 0755)
	if err := extractArchive(src, dest); err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	rel, _ := filepath.Rel(BASE_DIR, dest)
	jsonOK(w, map[string]string{"message": "Extracted", "dest": rel})
}

// ── Terminal Handlers ─────────────────────────────────────────────────────────

func handleTerminalStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	sid := generateUUID()
	ts := newTerminalSession(sid)
	if err := ts.Start(220, 50); err != nil {
		jsonErr(w, "Failed to start terminal: "+err.Error(), 500)
		return
	}
	termSessionsMu.Lock()
	termSessions[sid] = ts
	termSessionsMu.Unlock()

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	parts := strings.Split(shell, "/")
	jsonOK(w, map[string]string{"session_id": sid, "shell": parts[len(parts)-1]})
}

// List all active sessions (for reconnect on page reload)
func handleTerminalList(w http.ResponseWriter, r *http.Request) {
	termSessionsMu.Lock()
	defer termSessionsMu.Unlock()
	type sessionInfo struct {
		SessionID    string `json:"session_id"`
		Running      bool   `json:"running"`
		LastActivity string `json:"last_activity"`
	}
	var list []sessionInfo
	for sid, ts := range termSessions {
		ts.mu.Lock()
		running := ts.running
		la := ts.lastActivity
		ts.mu.Unlock()
		list = append(list, sessionInfo{
			SessionID:    sid,
			Running:      running,
			LastActivity: la.Format(time.RFC3339),
		})
	}
	jsonOK(w, list)
}

// Attach to existing session — returns scrollback so frontend can replay
func handleTerminalAttach(w http.ResponseWriter, r *http.Request) {
	sid := strings.TrimPrefix(r.URL.Path, "/api/terminal/attach/")
	termSessionsMu.Lock()
	ts := termSessions[sid]
	termSessionsMu.Unlock()
	if ts == nil {
		jsonErr(w, "Session not found", 404)
		return
	}
	ts.mu.Lock()
	scrollback := ts.scrollback
	running := ts.running
	ts.mu.Unlock()

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	parts := strings.Split(shell, "/")
	jsonOK(w, map[string]interface{}{
		"session_id": sid,
		"shell":      parts[len(parts)-1],
		"running":    running,
		"scrollback": base64.StdEncoding.EncodeToString(scrollback),
	})
}

func handleTerminalInput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	sid := strings.TrimPrefix(r.URL.Path, "/api/terminal/input/")
	termSessionsMu.Lock()
	ts := termSessions[sid]
	termSessionsMu.Unlock()
	if ts == nil {
		jsonErr(w, "Session not found", 404)
		return
	}
	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	raw := data["input"]

	// Try standard base64 first, then raw (no padding), then fallback to plain text
	var decoded []byte
	var err error
	decoded, err = base64.StdEncoding.DecodeString(raw)
	if err != nil {
		// Strip any existing padding and re-pad correctly
		stripped := strings.TrimRight(raw, "=")
		switch len(stripped) % 4 {
		case 2:
			stripped += "=="
		case 3:
			stripped += "="
		}
		decoded, err = base64.StdEncoding.DecodeString(stripped)
	}
	if err != nil {
		// Not base64, write as-is
		ts.Write([]byte(raw))
	} else {
		ts.Write(decoded)
	}
	jsonOK(w, map[string]bool{"success": true})
}

func handleTerminalOutput(w http.ResponseWriter, r *http.Request) {
	sid := strings.TrimPrefix(r.URL.Path, "/api/terminal/output/")
	termSessionsMu.Lock()
	ts := termSessions[sid]
	termSessionsMu.Unlock()
	if ts == nil {
		jsonErr(w, "Session not found", 404)
		return
	}
	raw := ts.ReadOutput()
	ts.mu.Lock()
	running := ts.running
	ts.mu.Unlock()
	encoded := base64.StdEncoding.EncodeToString(raw)
	jsonOK(w, map[string]interface{}{"output": encoded, "running": running})
}

func handleTerminalResize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	sid := strings.TrimPrefix(r.URL.Path, "/api/terminal/resize/")
	termSessionsMu.Lock()
	ts := termSessions[sid]
	termSessionsMu.Unlock()
	if ts == nil {
		jsonErr(w, "Session not found", 404)
		return
	}
	var data map[string]int
	json.NewDecoder(r.Body).Decode(&data)
	cols := uint16(data["cols"])
	rows := uint16(data["rows"])
	if cols == 0 {
		cols = 80
	}
	if rows == 0 {
		rows = 24
	}
	ts.Resize(cols, rows)
	jsonOK(w, map[string]bool{"success": true})
}

func handleTerminalClose(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	sid := strings.TrimPrefix(r.URL.Path, "/api/terminal/close/")
	termSessionsMu.Lock()
	ts := termSessions[sid]
	delete(termSessions, sid)
	termSessionsMu.Unlock()
	if ts != nil {
		ts.Close()
	}
	jsonOK(w, map[string]bool{"success": true})
}

func handleSaveFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonErr(w, "Method not allowed", 405)
		return
	}
	var data struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		jsonErr(w, "Bad request", 400)
		return
	}
	if data.Path == "" {
		jsonErr(w, "Path required", 400)
		return
	}
	fullPath, ok := securePath(BASE_DIR, data.Path)
	if !ok {
		jsonErr(w, "Invalid path", 400)
		return
	}
	if err := os.WriteFile(fullPath, []byte(data.Content), 0644); err != nil {
		jsonErr(w, err.Error(), 500)
		return
	}
	jsonOK(w, map[string]string{"message": "Saved"})
}



func router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("/logout", handleLogout)
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/system-info", handleSystemInfo)
	mux.HandleFunc("/api/processes", handleProcesses)
	mux.HandleFunc("/api/process/kill/", handleKillProcess)
	mux.HandleFunc("/api/ping", handlePing)
	mux.HandleFunc("/api/files", handleListFiles)
	mux.HandleFunc("/api/files/", handleListFiles)
	mux.HandleFunc("/api/upload", handleUpload)
	mux.HandleFunc("/api/download/", handleDownload)
	mux.HandleFunc("/api/delete/", handleDelete)
	mux.HandleFunc("/api/rename", handleRename)
	mux.HandleFunc("/api/move", handleMove)
	mux.HandleFunc("/api/copy", handleCopy)
	mux.HandleFunc("/api/create/file", handleCreateFile)
	mux.HandleFunc("/api/create/folder", handleCreateFolder)
	mux.HandleFunc("/api/content/", handleContent)
	mux.HandleFunc("/api/save", handleSaveFile)
	mux.HandleFunc("/api/serve/", handleServe)
	mux.HandleFunc("/api/hide", handleHide)
	mux.HandleFunc("/api/extract", handleExtract)
	mux.HandleFunc("/api/terminal/start", handleTerminalStart)
	mux.HandleFunc("/api/terminal/list", handleTerminalList)
	mux.HandleFunc("/api/terminal/attach/", handleTerminalAttach)
	mux.HandleFunc("/api/terminal/input/", handleTerminalInput)
	mux.HandleFunc("/api/terminal/output/", handleTerminalOutput)
	mux.HandleFunc("/api/terminal/resize/", handleTerminalResize)
	mux.HandleFunc("/api/terminal/close/", handleTerminalClose)
	return authMiddleware(mux)
}

// ── Utils ─────────────────────────────────────────────────────────────────────

func fileExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

var uuidCounter uint64
var uuidMu sync.Mutex

func generateUUID() string {
	uuidMu.Lock()
	uuidCounter++
	n := uuidCounter
	uuidMu.Unlock()
	return fmt.Sprintf("%016x-%d", time.Now().UnixNano(), n)
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	go cleanSessions()
	go cleanIdleTerminals()

	listenAddr := os.Getenv("YHUJIN_ADDR")
	if listenAddr == "" {
		listenAddr = "0.0.0.0:5000"
	}

	fmt.Println(strings.Repeat("=", 65))
	fmt.Println("YHUJIN CLOUDHOST v2.0 - Enterprise Server Manager (Go)")
	fmt.Println(strings.Repeat("=", 65))
	fmt.Printf("Server    : http://%s\n", listenAddr)
	fmt.Printf("Base Dir  : %s\n", BASE_DIR)
	fmt.Printf("Storage   : %s\n", STORAGE_DIR)
	fmt.Println(strings.Repeat("=", 65))

	info := getSystemInfo()
	if hostname, ok := info["hostname"].(string); ok {
		fmt.Printf("Hostname  : %s\n", hostname)
	}
	fmt.Printf("Platform  : %s\n", runtime.GOOS)
	if memMap, ok := info["memory"].(map[string]interface{}); ok {
		fmt.Printf("Memory    : %v / %v\n", memMap["used"], memMap["total"])
	}
	if diskMap, ok := info["disk"].(map[string]interface{}); ok {
		fmt.Printf("Disk      : %v / %v\n", diskMap["used"], diskMap["total"])
	}
	if ips, ok := info["ip_addresses"].([]string); ok && len(ips) > 0 {
		fmt.Printf("IP        : %s\n", strings.Join(ips, ", "))
	}
	fmt.Println(strings.Repeat("=", 65))

	if os.Getenv("YHUJIN_TOKEN") == "" {
		fmt.Printf("\n⚠  Token auto-generated (set YHUJIN_TOKEN env untuk permanen):\n")
		fmt.Printf("   Token : %s\n", AUTH_TOKEN)
	} else {
		fmt.Printf("\n✓  Token dari env YHUJIN_TOKEN\n")
	}
	fmt.Printf("   Login : http://%s/login\n\n", listenAddr)
	if DEV_MODE {
		log.Printf("⚠  DEV MODE — cache nonaktif, anti-debug nonaktif")
	}
	log.Fatal(http.ListenAndServe(listenAddr, router()))
}
