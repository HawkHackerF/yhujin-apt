package main

var htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Yhujin Cloudhost</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@300;400;500;700&display=swap" rel="stylesheet">
<script src="https://cdn.jsdelivr.net/npm/xterm@5.3.0/lib/xterm.min.js"></script>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/xterm@5.3.0/css/xterm.css">
<script src="https://cdn.jsdelivr.net/npm/xterm-addon-fit@0.8.0/lib/xterm-addon-fit.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/xterm-addon-web-links@0.9.0/lib/xterm-addon-web-links.min.js"></script>
<style>
:root{
  --bg0:#18181d;--bg1:#202027;--bg2:#28282f;--bg3:#32323c;
  --tx0:#dde1e8;--tx1:#8e95a3;--tx2:#5a6070;
  --blue:#4a9eff;--green:#3fb950;--yellow:#d29922;--red:#f85149;--purple:#bc8cff;--cyan:#39d0d8;
  --border:#38383f;
}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:'JetBrains Mono',monospace;background:var(--bg0);color:var(--tx0);height:100vh;overflow:hidden;display:flex;flex-direction:column}
.tab-bar{display:flex;align-items:center;background:var(--bg1);border-bottom:1px solid var(--border);height:40px;flex-shrink:0;overflow-x:auto;overflow-y:hidden;scrollbar-width:none}
.tab-bar::-webkit-scrollbar{display:none}
.tab-brand{padding:0 16px;color:var(--blue);font-size:.9rem;font-weight:700;letter-spacing:1.5px;white-space:nowrap;border-right:1px solid var(--border);height:100%;display:flex;align-items:center;flex-shrink:0}
.tab-item{display:flex;align-items:center;gap:7px;padding:0 15px;height:100%;cursor:pointer;color:var(--tx1);font-size:.8rem;white-space:nowrap;border-right:1px solid var(--border);transition:background .15s,color .15s;position:relative;flex-shrink:0;user-select:none}
.tab-item:hover{background:var(--bg3);color:var(--tx0)}
.tab-item.active{background:var(--bg0);color:var(--blue)}
.tab-item.active::after{content:'';position:absolute;bottom:0;left:0;right:0;height:2px;background:var(--blue)}
.tab-add{padding:0 14px;height:100%;display:flex;align-items:center;cursor:pointer;color:var(--tx2);font-size:1.2rem;transition:color .15s}
.tab-add:hover{color:var(--tx0)}
.app-body{display:flex;flex:1;overflow:hidden}
.sidebar{width:48px;background:var(--bg1);border-right:1px solid var(--border);display:flex;flex-direction:column;align-items:center;padding:6px 0;gap:2px;flex-shrink:0}
.sb-icon{width:34px;height:34px;display:flex;align-items:center;justify-content:center;border-radius:7px;cursor:pointer;color:var(--tx2);transition:all .15s}
.sb-icon:hover{background:var(--bg3);color:var(--tx0)}
.sb-icon.active{background:rgba(74,158,255,.12);color:var(--blue)}
.sb-sep{width:26px;height:1px;background:var(--border);margin:4px 0}
.sb-bottom{margin-top:auto}
.panel-container{flex:1;overflow:hidden;position:relative}
.panel{position:absolute;inset:0;display:none;flex-direction:column;overflow:hidden}
.panel.active{display:flex}
.topbar{height:42px;background:var(--bg1);border-bottom:1px solid var(--border);display:flex;align-items:center;padding:0 10px;gap:7px;flex-shrink:0}
.nav-btn{width:27px;height:27px;display:flex;align-items:center;justify-content:center;border:none;background:var(--bg2);color:var(--tx1);border-radius:5px;cursor:pointer;transition:all .15s;flex-shrink:0}
.nav-btn:hover{background:var(--bg3);color:var(--tx0)}
.nav-btn:disabled{opacity:.3;cursor:default}
.breadcrumb{flex:1;display:flex;align-items:center;gap:3px;font-size:.78rem;overflow:hidden;min-width:0}
.bc-item{color:var(--tx1);cursor:pointer;white-space:nowrap;padding:2px 5px;border-radius:4px;transition:all .15s}
.bc-item:hover{background:var(--bg3);color:var(--tx0)}
.bc-item.last{color:var(--tx0);cursor:default}
.bc-item.last:hover{background:none}
.bc-sep{color:var(--tx2);font-size:.7rem}
.topbar-right{display:flex;align-items:center;gap:7px;margin-left:auto;flex-shrink:0}
.search-box{background:var(--bg2);border:1px solid var(--border);border-radius:5px;padding:3px 9px;display:flex;align-items:center;gap:5px}
.search-box input{background:none;border:none;color:var(--tx0);font-family:inherit;font-size:.78rem;width:160px;outline:none}
.search-box input::placeholder{color:var(--tx2)}
.vbtn{background:var(--bg2);border:1px solid var(--border);color:var(--tx1);padding:3px 9px;border-radius:4px;cursor:pointer;font-size:.75rem;font-family:inherit;transition:all .15s}
.vbtn:hover{background:var(--bg3);color:var(--tx0)}
.vbtn.active{background:var(--blue);border-color:var(--blue);color:#000}
.toolbar{padding:7px 12px;background:var(--bg1);border-bottom:1px solid var(--border);display:flex;gap:7px;flex-wrap:wrap;flex-shrink:0;align-items:center}
.btn{background:var(--bg2);border:1px solid var(--border);color:var(--tx1);padding:5px 11px;border-radius:5px;cursor:pointer;font-size:.76rem;font-family:inherit;display:flex;align-items:center;gap:5px;transition:all .15s;white-space:nowrap}
.btn:hover{background:var(--bg3);color:var(--tx0)}
.btn.primary{background:var(--blue);border-color:var(--blue);color:#000;font-weight:500}
.btn.primary:hover{background:#5aa8ff}
.btn.danger{background:rgba(248,81,73,.1);border-color:var(--red);color:var(--red)}
.btn.danger:hover{background:rgba(248,81,73,.2)}
.btn.ok{background:rgba(63,185,80,.1);border-color:var(--green);color:var(--green)}
.btn.ok:hover{background:rgba(63,185,80,.22)}
.btn.warn{background:rgba(210,153,34,.1);border-color:var(--yellow);color:var(--yellow)}
.btn.warn:hover{background:rgba(210,153,34,.22)}
.fl-header{display:grid;grid-template-columns:28px 20px 1fr 80px 100px 120px 155px;padding:6px 12px;background:var(--bg1);color:var(--tx2);font-size:.72rem;text-transform:uppercase;letter-spacing:.5px;border-bottom:1px solid var(--border);flex-shrink:0}
.fl-body{flex:1;overflow-y:auto}
.fi{display:grid;grid-template-columns:28px 20px 1fr 80px 100px 120px 155px;padding:5px 12px;border-bottom:1px solid rgba(56,56,63,.4);cursor:pointer;transition:background .1s;font-size:.83rem;align-items:center}
.fi:hover{background:var(--bg3)}
.fi.sel{background:rgba(74,158,255,.1)}
.fi.sel:hover{background:rgba(74,158,255,.15)}
.fi-icon{display:flex;align-items:center;justify-content:center;color:var(--tx1)}
.fi.hid{opacity:.43}
.fi-cb{display:flex;align-items:center;justify-content:center}
.fi-cb input,.fgi-cb input{width:13px;height:13px;accent-color:var(--blue);cursor:pointer}
.fi-name.hid{color:var(--tx2);font-style:italic}
.fgi.hid{opacity:.43}
.fgi-cb{position:absolute;top:5px;left:5px}
.ib{background:none;border:1px solid var(--border);color:var(--tx1);padding:1px 5px;border-radius:3px;cursor:pointer;font-size:.67rem;font-family:inherit;transition:all .1s;white-space:nowrap}
.ib:hover{background:var(--bg3);color:var(--tx0)}
.ib.del:hover{border-color:var(--red);color:var(--red)}
.ib.ex:hover{border-color:var(--green);color:var(--green)}
.ib.hd:hover{border-color:var(--yellow);color:var(--yellow)}
.sel-info{font-size:.72rem;color:var(--tx2);margin-left:3px}
.extract-log{background:#0c0c10;border:1px solid var(--border);border-radius:6px;padding:10px;font-size:.74rem;color:var(--green);min-height:55px;max-height:155px;overflow-y:auto;margin-top:10px;font-family:monospace;line-height:1.6}
.fi-icon.d{color:var(--blue)}
.fi-name{color:var(--tx0);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;padding-right:6px}
.fi-name.d{color:var(--blue)}
.fi-size{color:var(--tx1);font-size:.75rem}
.fi-perms{color:var(--tx2);font-size:.72rem;font-family:monospace}
.fi-own{color:var(--tx2);font-size:.75rem}
.fi-mod{color:var(--tx1);font-size:.72rem;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.fi-act{display:flex;gap:3px;flex-wrap:nowrap;overflow:hidden}
.fi-act button{background:none;border:1px solid var(--border);color:var(--tx1);padding:1px 5px;border-radius:3px;cursor:pointer;font-size:.67rem;font-family:inherit;transition:all .1s;white-space:nowrap;flex-shrink:0}
.fi-act button:hover{background:var(--bg3);color:var(--tx0)}
.fi-act button.del:hover{border-color:var(--red);color:var(--red)}
.file-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(105px,1fr));gap:9px;padding:12px}
.fgi{background:var(--bg1);border:1px solid var(--border);border-radius:8px;padding:12px 8px 8px;text-align:center;cursor:pointer;transition:all .15s;position:relative}
.fgi:hover{background:var(--bg2);border-color:var(--blue)}
.fgi.sel{border-color:var(--blue);background:rgba(74,158,255,.07)}
.fgi-icon{display:flex;align-items:center;justify-content:center;color:var(--blue);margin-bottom:7px}
.fgi-icon.f{color:var(--tx1)}
.fgi-name{font-size:.72rem;color:var(--tx0);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.fgi-size{font-size:.68rem;color:var(--tx2);margin-top:2px}
.empty-state{text-align:center;padding:60px 20px;color:var(--tx2);font-size:.88rem}
/* TERMINAL */
.term-tabbar{background:var(--bg1);border-bottom:1px solid var(--border);display:flex;align-items:center;height:35px;flex-shrink:0;overflow:hidden}
.term-tabbar-tabs{display:flex;align-items:center;flex:1;overflow-x:auto;scrollbar-width:none;height:100%}
.term-tabbar-tabs::-webkit-scrollbar{display:none}
.term-tab{display:flex;align-items:center;gap:6px;padding:0 11px;height:100%;cursor:pointer;color:var(--tx1);font-size:.75rem;white-space:nowrap;border-right:1px solid var(--border);transition:all .15s;position:relative;user-select:none;flex-shrink:0}
.term-tab:hover{background:var(--bg3);color:var(--tx0)}
.term-tab.active{background:var(--bg0);color:var(--green)}
.term-tab.active::after{content:'';position:absolute;bottom:0;left:0;right:0;height:2px;background:var(--green)}
.tdot{width:6px;height:6px;border-radius:50%;background:var(--red);flex-shrink:0;transition:background .3s}
.term-tab.active .tdot{background:var(--green);animation:blink 2s infinite}
.tcl{width:14px;height:14px;border-radius:3px;display:flex;align-items:center;justify-content:center;font-size:10px;color:var(--tx2);margin-left:2px;line-height:1}
.tcl:hover{background:rgba(248,81,73,.18);color:var(--red)}
.tnew{padding:0 12px;height:100%;display:flex;align-items:center;cursor:pointer;color:var(--tx2);font-size:1.1rem;transition:color .15s;flex-shrink:0}
.tnew:hover{color:var(--green)}
.term-killall{display:flex;align-items:center;gap:5px;margin:0 8px;padding:4px 10px;border-radius:5px;border:1px solid var(--border);background:var(--bg2);color:var(--tx2);font-size:.72rem;cursor:pointer;transition:all .15s;flex-shrink:0;font-family:inherit}
.term-killall:hover{border-color:var(--red);color:var(--red);background:rgba(248,81,73,.08)}
.term-instances{flex:1;overflow:hidden;position:relative;background:#0c0c10}
.term-inst{position:absolute;inset:0;display:none;padding:2px}
.term-inst.active{display:block}
/* Right-click context menu for terminal */
.ctx-menu{position:fixed;background:var(--bg2);border:1px solid var(--border);border-radius:7px;padding:4px 0;z-index:9000;min-width:150px;box-shadow:0 8px 24px rgba(0,0,0,.5);display:none}
.ctx-menu.active{display:block}
.ctx-item{padding:7px 16px;font-size:.78rem;color:var(--tx0);cursor:pointer;display:flex;align-items:center;gap:8px;transition:background .1s}
.ctx-item:hover{background:var(--bg3)}
.ctx-sep{height:1px;background:var(--border);margin:4px 0}
@keyframes blink{0%,100%{opacity:1}50%{opacity:.35}}
/* PROCESSES */
.proc-toolbar{padding:7px 12px;background:var(--bg1);border-bottom:1px solid var(--border);display:flex;align-items:center;gap:9px;flex-shrink:0}
.proc-search{background:var(--bg2);border:1px solid var(--border);color:var(--tx0);font-family:inherit;font-size:.8rem;padding:5px 10px;border-radius:5px;outline:none;width:210px}
.proc-count{color:var(--tx2);font-size:.75rem;margin-left:auto}
.proc-wrap{flex:1;overflow-y:auto}
.proc-table{width:100%;border-collapse:collapse;font-size:.8rem}
.proc-table th{position:sticky;top:0;background:var(--bg1);color:var(--tx2);text-transform:uppercase;font-size:.7rem;letter-spacing:.5px;padding:7px 10px;border-bottom:1px solid var(--border);text-align:left;cursor:pointer;user-select:none}
.proc-table th:hover{color:var(--tx0)}
.proc-table td{padding:5px 10px;border-bottom:1px solid rgba(56,56,63,.35);color:var(--tx1);white-space:nowrap;overflow:hidden;max-width:180px;text-overflow:ellipsis}
.proc-table tr:hover td{background:var(--bg3)}
.pid{color:var(--tx2);font-family:monospace}.pname{color:var(--tx0);font-weight:500}
.pcpu{color:var(--blue)}.pmem{color:var(--purple)}
.prun{color:var(--green)}.pslp{color:var(--tx2)}
.kill-btn{background:none;border:1px solid var(--border);color:var(--tx2);border-radius:3px;padding:1px 6px;cursor:pointer;font-size:.7rem;font-family:inherit;transition:all .1s}
.kill-btn:hover{border-color:var(--red);color:var(--red)}
/* MONITOR */
.mon-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(270px,1fr));gap:13px;padding:14px;overflow-y:auto;flex:1}
.mon-card{background:var(--bg1);border:1px solid var(--border);border-radius:10px;padding:15px}
.mon-title{font-size:.72rem;text-transform:uppercase;letter-spacing:.8px;color:var(--tx2);margin-bottom:11px;display:flex;align-items:center;gap:6px}
.mon-big{font-size:1.9rem;font-weight:700;color:var(--tx0);line-height:1;margin-bottom:5px}
.mon-sub{font-size:.75rem;color:var(--tx2);margin-bottom:9px}
.gauge{width:100%;height:7px;background:var(--bg2);border-radius:4px;overflow:hidden;margin:5px 0}
.gfill{height:100%;border-radius:4px;transition:width .6s ease}
.gfill.cpu{background:linear-gradient(90deg,var(--blue),#7ab8ff)}
.gfill.mem{background:linear-gradient(90deg,var(--purple),#d4b0ff)}
.gfill.dsk{background:linear-gradient(90deg,var(--yellow),#e8c86e)}
.gfill.swp{background:linear-gradient(90deg,var(--cyan),#7feaf0)}
.mon-det{display:grid;grid-template-columns:1fr 1fr;gap:5px;margin-top:9px}
.mon-det-i{font-size:.75rem}
.mon-lbl{color:var(--tx2)}.mon-val{color:var(--tx0);font-weight:500}
.core-grid{display:flex;flex-direction:column;gap:4px;margin-top:7px}
.core-i{display:grid;grid-template-columns:26px 1fr 34px;align-items:center;gap:6px;font-size:.68rem}
.core-lbl{color:var(--tx2)}
.core-bar{width:100%;height:5px;background:var(--bg2);border-radius:2px;overflow:hidden}
.core-fill{height:100%;background:var(--blue);border-radius:2px;transition:width .4s}
/* SETTINGS */
.sett-wrap{flex:1;overflow-y:auto;padding:18px}
.sett-info{background:rgba(74,158,255,.05);border:1px solid rgba(74,158,255,.18);border-radius:8px;padding:11px 15px;margin-bottom:14px;font-size:.8rem;color:var(--tx1);line-height:1.7}
.sett-info strong{color:var(--blue)}
.sett-sec{background:var(--bg1);border:1px solid var(--border);border-radius:9px;margin-bottom:14px;overflow:hidden}
.sett-sec-title{padding:10px 16px;border-bottom:1px solid var(--border);font-size:.73rem;text-transform:uppercase;letter-spacing:.8px;color:var(--tx2);display:flex;align-items:center;gap:7px}
.sett-row{padding:12px 16px;display:flex;align-items:center;border-bottom:1px solid rgba(56,56,63,.35);gap:12px}
.sett-row:last-child{border-bottom:none}
.sett-lbl{flex:1}.sett-lbl strong{display:block;font-size:.84rem;color:var(--tx0)}.sett-lbl span{font-size:.74rem;color:var(--tx2)}
.sett-inp{background:var(--bg2);border:1px solid var(--border);color:var(--tx0);font-family:inherit;font-size:.82rem;padding:5px 11px;border-radius:5px;outline:none;width:190px}
.sett-inp:focus{border-color:var(--blue)}
.toggle{width:42px;height:22px;background:var(--bg2);border:1px solid var(--border);border-radius:11px;cursor:pointer;position:relative;transition:background .2s;flex-shrink:0}
.toggle.on{background:var(--blue);border-color:var(--blue)}
.toggle::after{content:'';position:absolute;top:2px;left:2px;width:16px;height:16px;border-radius:50%;background:#fff;transition:left .2s}
.toggle.on::after{left:22px}
/* SERVERS PANEL */
.srv-wrap{flex:1;overflow-y:auto;padding:16px}
.srv-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(290px,1fr));gap:12px;margin-top:12px;min-width:0}
.srv-card{background:var(--bg1);border:1px solid var(--border);border-radius:10px;padding:14px;position:relative;transition:border-color .15s;overflow:hidden;min-width:0;word-break:break-word}
.srv-card:hover{border-color:var(--blue)}
.srv-card.connected{border-left:3px solid var(--green)}
.srv-card.disconnected{border-left:3px solid var(--tx2)}
.srv-card.error{border-left:3px solid var(--red)}
.srv-name{font-size:.88rem;font-weight:700;color:var(--tx0);margin-bottom:4px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;padding-right:24px}
.srv-host{font-size:.75rem;color:var(--tx2);margin-bottom:8px;font-family:monospace;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.srv-tags{display:flex;flex-wrap:wrap;gap:4px;margin-bottom:8px;width:100%;overflow:hidden}
.srv-status{display:inline-flex;align-items:center;gap:5px;font-size:.72rem;padding:2px 8px;border-radius:10px;background:var(--bg2)}
.srv-dot{width:6px;height:6px;border-radius:50%}
.srv-dot.green{background:var(--green);animation:blink 2s infinite}
.srv-dot.red{background:var(--red)}
.srv-dot.gray{background:var(--tx2)}
.srv-actions{display:flex;gap:6px;margin-top:10px}
.srv-del-btn{position:absolute;top:10px;right:10px;background:none;border:none;color:var(--tx2);cursor:pointer;font-size:14px;padding:3px 6px;border-radius:4px}
.srv-del-btn:hover{background:rgba(248,81,73,.15);color:var(--red)}
.srv-empty{text-align:center;padding:60px 20px;color:var(--tx2);font-size:.88rem}
/* MODAL */
.modal{display:none;position:fixed;inset:0;background:rgba(0,0,0,.75);z-index:1000;align-items:center;justify-content:center}
.modal.active{display:flex}
.modal-box{background:var(--bg1);border:1px solid var(--border);border-radius:10px;width:540px;max-width:92%;max-height:80vh;overflow:auto}
.modal-hd{padding:13px 17px;border-bottom:1px solid var(--border);display:flex;justify-content:space-between;align-items:center}
.modal-hd h3{font-size:.9rem;color:var(--tx0)}
.modal-cl{background:none;border:none;color:var(--tx2);font-size:1.2rem;cursor:pointer}
.modal-cl:hover{color:var(--red)}
.modal-bd{padding:16px}
.modal-ft{padding:11px 17px;border-top:1px solid var(--border);display:flex;justify-content:flex-end;gap:7px}
.upload-area{border:2px dashed var(--border);border-radius:8px;padding:32px;text-align:center;cursor:pointer;transition:all .2s;color:var(--tx1);font-size:.85rem}
.upload-area:hover{border-color:var(--blue);color:var(--tx0)}
.upload-area.dragover{border-color:var(--green);background:rgba(63,185,80,.04)}
.upload-list{max-height:170px;overflow-y:auto;margin-top:11px}
.upload-item{display:flex;justify-content:space-between;padding:4px 7px;border-bottom:1px solid var(--border);font-size:.79rem;color:var(--tx1)}
.form-row{margin-bottom:12px}
.form-row label{display:block;font-size:.75rem;color:var(--tx2);margin-bottom:5px}
.form-row input{width:100%;background:var(--bg2);border:1px solid var(--border);color:var(--tx0);font-family:inherit;font-size:.82rem;padding:7px 11px;border-radius:5px;outline:none}
.form-row input:focus{border-color:var(--blue)}
/* TOAST */
.toast-box{position:fixed;bottom:18px;right:18px;display:flex;flex-direction:column;gap:7px;z-index:9999}
.toast{background:var(--bg2);border:1px solid var(--border);border-radius:7px;padding:9px 14px;font-size:.79rem;color:var(--tx0);animation:toastin .2s ease;min-width:200px}
.toast.success{border-left:3px solid var(--green)}
.toast.error{border-left:3px solid var(--red)}
.toast.info{border-left:3px solid var(--blue)}
@keyframes toastin{from{transform:translateX(16px);opacity:0}to{transform:none;opacity:1}}
/* FILE VIEWER */
.viewer-modal{display:none;position:fixed;inset:0;background:rgba(0,0,0,.88);z-index:2000;flex-direction:column}
.viewer-modal.active{display:flex}
.viewer-topbar{height:46px;background:var(--bg1);border-bottom:1px solid var(--border);display:flex;align-items:center;padding:0 12px;gap:8px;flex-shrink:0}
.viewer-title{flex:1;font-size:.84rem;color:var(--tx0);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;min-width:0}
.viewer-meta{font-size:.72rem;color:var(--tx2);white-space:nowrap}
.viewer-body{flex:1;overflow:hidden;display:flex;flex-direction:column;min-height:0}
.viewer-code{flex:1;overflow:auto;background:#0c0c10;margin:0;padding:16px;font-family:'JetBrains Mono',monospace;font-size:.82rem;line-height:1.6;color:#dde1e8;white-space:pre;tab-size:4}
.viewer-img{flex:1;display:flex;align-items:center;justify-content:center;background:#0c0c10;overflow:auto;padding:16px}
.viewer-img img{max-width:100%;max-height:100%;object-fit:contain;border-radius:4px}
.viewer-media{flex:1;display:flex;align-items:center;justify-content:center;background:#0c0c10;padding:20px}
.viewer-media video,.viewer-media audio{max-width:100%;max-height:100%;outline:none}
.viewer-pdf{flex:1;display:flex;flex-direction:column;background:#0c0c10}
.viewer-pdf iframe{flex:1;border:none;width:100%;height:100%}
.viewer-unsupported{flex:1;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:14px;color:var(--tx2);font-size:.88rem}
.viewer-nav-btn{width:27px;height:27px;display:flex;align-items:center;justify-content:center;border:none;background:var(--bg2);color:var(--tx1);border-radius:5px;cursor:pointer;transition:all .15s;flex-shrink:0}
.viewer-nav-btn:hover{background:var(--bg3);color:var(--tx0)}
.viewer-nav-btn:disabled{opacity:.3;cursor:default}
.line-nums{counter-reset:line;display:block}
.line-nums .ln{counter-increment:line;display:block}
.line-nums .ln::before{content:counter(line);display:inline-block;width:3em;color:var(--tx2);margin-right:1em;text-align:right;user-select:none;border-right:1px solid var(--border);padding-right:.5em;font-size:.78rem}
.ib.ed:hover{border-color:var(--green);color:var(--green)}
::-webkit-scrollbar{width:6px;height:6px}
::-webkit-scrollbar-track{background:transparent}
::-webkit-scrollbar-thumb{background:var(--bg2);border-radius:3px}
::-webkit-scrollbar-thumb:hover{background:var(--bg3)}
</style>
<script>
// ── Mode & Cache ──────────────────────────────────────────────────────────────
const DEV_MODE = __DEV_MODE__;
const BUILD_HASH = '__BUILD_HASH__';

if(DEV_MODE){
  // Dev: log build hash setiap load untuk konfirmasi cache tidak dipakai
  console.info('%c[YHUJIN DEV] Build: ' + BUILD_HASH, 'color:#4a9eff;font-weight:bold');
  console.info('%c[YHUJIN DEV] Cache disabled — semua request fresh', 'color:#f0a500');
} else {
  // ── Anti-Debug (Production Only) ──────────────────────────────────────────
  // 1. Deteksi DevTools lewat debugger trap
  (function antiDebug(){
    function trap(){
      const start = performance.now();
      debugger; // jeda di sini kalau devtools terbuka
      if(performance.now() - start > 100){
        document.body.innerHTML = '';
        window.location.href = '/login';
      }
    }
    // Jalankan tiap 3 detik
    setInterval(trap, 3000);
    trap();
  })();

  // 2. Deteksi resize mendadak (devtools dock biasanya mengubah ukuran window)
  let _w = window.outerWidth, _h = window.outerHeight;
  setInterval(()=>{
    const threshold = 160;
    if(window.outerWidth - window.innerWidth > threshold ||
       window.outerHeight - window.innerHeight > threshold){
      document.body.innerHTML = '';
      window.location.href = '/login';
    }
  }, 1000);

  // 3. Disable klik kanan
  document.addEventListener('contextmenu', e => e.preventDefault());

  // 4. Disable F12, Ctrl+Shift+I, Ctrl+Shift+J, Ctrl+U
  document.addEventListener('keydown', e => {
    if(e.key === 'F12' ||
      (e.ctrlKey && e.shiftKey && ['I','J','C'].includes(e.key.toUpperCase())) ||
      (e.ctrlKey && e.key.toUpperCase() === 'U')){
      e.preventDefault();
      e.stopPropagation();
      return false;
    }
  });
}
</script>
</head>
<body>

<div class="tab-bar" id="tabBar">
  <div class="tab-brand">YHUJIN<span id="devBadge" style="display:none;margin-left:7px;font-size:.6rem;background:#f0a500;color:#000;padding:1px 6px;border-radius:3px;font-weight:700;letter-spacing:.5px">DEV</span></div>
  <div class="tab-item active" id="tab-files"     data-panel="files"     onclick="switchPanel('files')">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>File Manager
  </div>
  <div class="tab-item" id="tab-terminal"  data-panel="terminal"  onclick="switchPanel('terminal')">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>Terminal
  </div>
  <div class="tab-item" id="tab-processes" data-panel="processes" onclick="switchPanel('processes')">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="9" height="9" rx="1"/><rect x="13" y="2" width="9" height="9" rx="1"/><rect x="2" y="13" width="9" height="9" rx="1"/><rect x="13" y="13" width="9" height="9" rx="1"/></svg>Processes
  </div>
  <div class="tab-item" id="tab-monitor"   data-panel="monitor"   onclick="switchPanel('monitor')">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>Monitor
  </div>
  <div class="tab-item" id="tab-servers"   data-panel="servers"   onclick="switchPanel('servers')">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>Servers
  </div>
  <div class="tab-item" id="tab-settings"  data-panel="settings"  onclick="switchPanel('settings')">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/></svg>Settings
  </div>
  <div style="flex:1"></div>
  <div class="tab-item" onclick="doLogout()" title="Logout" style="color:var(--tx2);border-left:1px solid var(--border);flex-shrink:0">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>Logout
  </div>
</div>

<div class="app-body">
  <div class="sidebar">
    <div class="sb-icon active" id="sb-files"     onclick="switchPanel('files')"     title="File Manager">
      <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
    </div>
    <div class="sb-icon" id="sb-terminal"  onclick="switchPanel('terminal')"  title="Terminal">
      <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>
    </div>
    <div class="sb-icon" id="sb-processes" onclick="switchPanel('processes')" title="Processes">
      <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="9" height="9" rx="1"/><rect x="13" y="2" width="9" height="9" rx="1"/><rect x="2" y="13" width="9" height="9" rx="1"/><rect x="13" y="13" width="9" height="9" rx="1"/></svg>
    </div>
    <div class="sb-icon" id="sb-monitor"   onclick="switchPanel('monitor')"   title="System Monitor">
      <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
    </div>
    <div class="sb-icon" id="sb-servers"   onclick="switchPanel('servers')"   title="Servers">
      <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>
    </div>
    <div class="sb-sep"></div>
    <div class="sb-bottom">
      <div class="sb-icon" id="sb-settings" onclick="switchPanel('settings')" title="Settings">
        <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
      </div>
    </div>
  </div>

  <div class="panel-container">

    <!-- FILE MANAGER -->
    <div class="panel active" id="panel-files">
      <div class="topbar">
        <button class="nav-btn" id="btnBack" onclick="histBack()" disabled title="Back">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="15 18 9 12 15 6"/></svg>
        </button>
        <button class="nav-btn" id="btnFwd" onclick="histFwd()" disabled title="Forward">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
        </button>
        <button class="nav-btn" onclick="refreshFiles()" title="Refresh">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
        </button>
        <div class="breadcrumb" id="breadcrumb"></div>
        <div class="topbar-right">
          <div class="search-box">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
            <input id="searchInput" type="text" placeholder="Search..." oninput="filterFiles()">
          </div>
          <button class="vbtn active" id="vbl" onclick="setView('list',this)">List</button>
          <button class="vbtn" id="vbg" onclick="setView('grid',this)">Grid</button>
        </div>
      </div>
      <div class="toolbar">
        <button class="btn primary" onclick="openUpload()">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/><path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/></svg>Upload
        </button>
        <button class="btn" onclick="newFile()">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><polyline points="14 2 14 8 20 8"/><line x1="12" y1="18" x2="12" y2="12"/><line x1="9" y1="15" x2="15" y2="15"/></svg>New File
        </button>
        <button class="btn" onclick="newFolder()">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/><line x1="12" y1="11" x2="12" y2="17"/><line x1="9" y1="14" x2="15" y2="14"/></svg>New Folder
        </button>
        <button class="btn" id="btnRename" style="display:none" onclick="renameSelected()">Rename</button>
        <button class="btn" id="btnCut" style="display:none" onclick="cutSelected()">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="6" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><line x1="20" y1="4" x2="8.12" y2="15.88"/><line x1="14.47" y1="14.48" x2="20" y2="20"/><line x1="8.12" y1="8.12" x2="12" y2="12"/></svg>Cut
        </button>
        <button class="btn" id="btnCopy" style="display:none" onclick="copySelected()">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>Copy
        </button>
        <button class="btn ok" id="btnPaste" style="display:none" onclick="pasteHere()">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1"/></svg>Paste
        </button>
        <button class="btn warn" id="btnHide" style="display:none" onclick="toggleHideSelected()">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>Hide/Unhide
        </button>
        <button class="btn ok" id="btnExtract" style="display:none" onclick="extractSelected()">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="8 17 12 21 16 17"/><line x1="12" y1="15" x2="12" y2="3"/></svg>Extract
        </button>
        <button class="btn danger" id="btnDel" style="display:none" onclick="deleteSelected()">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/></svg>Delete
        </button>
        <span class="sel-info" id="selInfo" style="display:none"></span>
        <button class="btn" style="margin-left:auto" onclick="switchPanel('terminal')">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>Terminal
        </button>
      </div>
      <div class="fl-header" id="fl-header" style="display:grid">
        <div style="display:flex;align-items:center;justify-content:center">
          <input type="checkbox" id="selAllCb" onchange="toggleSelAll(this)" style="width:13px;height:13px;accent-color:var(--blue);cursor:pointer">
        </div>
        <div></div>
        <div>Name</div><div>Size</div><div>Owner</div><div>Modified</div><div>Actions</div>
      </div>
      <div class="fl-body" id="fileList"><div class="empty-state">Loading...</div></div>
    </div>

    <!-- TERMINAL -->
    <div class="panel" id="panel-terminal">
      <div class="term-tabbar" id="termTabBar">
        <div class="term-tabbar-tabs" id="termTabBarTabs">
          <div class="tnew" onclick="newTermTab()" title="New Terminal Tab">+</div>
        </div>
        <div class="term-killall" onclick="killAllTerms()" title="Kill all terminals">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/><path d="M10 11v6"/><path d="M14 11v6"/></svg>
          Kill All
        </div>
      </div>
      <div class="term-instances" id="termInstances"></div>
    </div>

    <!-- PROCESSES -->
    <div class="panel" id="panel-processes">
      <div class="topbar">
        <button class="nav-btn" onclick="loadProcs()" title="Refresh">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
        </button>
        <span style="font-size:.82rem;color:var(--tx1);margin-left:6px">Process Manager</span>
      </div>
      <div class="proc-toolbar">
        <input class="proc-search" id="procSearch" type="text" placeholder="Filter by name or PID..." oninput="filterProcs()">
        <div class="proc-count" id="procCount">— processes</div>
        <label style="font-size:.75rem;color:var(--tx2);display:flex;align-items:center;gap:5px">
          <input type="checkbox" id="autoRefProc" onchange="toggleProcRefresh()" checked> Auto-refresh
        </label>
      </div>
      <div class="proc-wrap">
        <table class="proc-table">
          <thead><tr>
            <th onclick="sortProc('pid')">PID</th>
            <th onclick="sortProc('name')">Name</th>
            <th onclick="sortProc('username')">User</th>
            <th onclick="sortProc('status')">Status</th>
            <th onclick="sortProc('cpu_percent')">CPU%</th>
            <th onclick="sortProc('memory_percent')">MEM%</th>
            <th>RSS</th><th>Started</th><th>Action</th>
          </tr></thead>
          <tbody id="procBody"></tbody>
        </table>
      </div>
    </div>

    <!-- MONITOR -->
    <div class="panel" id="panel-monitor">
      <div class="topbar">
        <span style="font-size:.82rem;color:var(--tx1)">System Monitor</span>
        <span id="monUpd" style="font-size:.72rem;color:var(--tx2);margin-left:10px"></span>
      </div>
      <div class="mon-grid">
        <div class="mon-card">
          <div class="mon-title"><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="16" height="16" rx="2"/><rect x="9" y="9" width="6" height="6"/><line x1="9" y1="1" x2="9" y2="4"/><line x1="15" y1="1" x2="15" y2="4"/><line x1="9" y1="20" x2="9" y2="23"/><line x1="15" y1="20" x2="15" y2="23"/><line x1="20" y1="9" x2="23" y2="9"/><line x1="20" y1="14" x2="23" y2="14"/><line x1="1" y1="9" x2="4" y2="9"/><line x1="1" y1="14" x2="4" y2="14"/></svg>CPU</div>
          <div class="mon-big" id="mCpuV">0%</div>
          <div class="mon-sub" id="mCpuSub">— cores — threads</div>
          <div class="gauge"><div class="gfill cpu" id="mCpuB" style="width:0%"></div></div>
          <div class="mon-det">
            <div class="mon-det-i"><div class="mon-lbl">Frequency</div><div class="mon-val" id="mCpuF">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">Load avg 1m</div><div class="mon-val" id="mLoad">—</div></div>
          </div>
          <div class="core-grid" id="mCores"></div>
        </div>
        <div class="mon-card">
          <div class="mon-title">Memory</div>
          <div class="mon-big" id="mMemV">0%</div>
          <div class="mon-sub" id="mMemSub">— of —</div>
          <div class="gauge"><div class="gfill mem" id="mMemB" style="width:0%"></div></div>
          <div class="mon-det">
            <div class="mon-det-i"><div class="mon-lbl">Used</div><div class="mon-val" id="mMemU">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">Available</div><div class="mon-val" id="mMemA">—</div></div>
          </div>
        </div>
        <div class="mon-card">
          <div class="mon-title">Swap</div>
          <div class="mon-big" id="mSwpV">0%</div>
          <div class="mon-sub" id="mSwpSub">—</div>
          <div class="gauge"><div class="gfill swp" id="mSwpB" style="width:0%"></div></div>
        </div>
        <div class="mon-card">
          <div class="mon-title">Disk</div>
          <div class="mon-big" id="mDskV">0%</div>
          <div class="mon-sub" id="mDskSub">—</div>
          <div class="gauge"><div class="gfill dsk" id="mDskB" style="width:0%"></div></div>
          <div class="mon-det">
            <div class="mon-det-i"><div class="mon-lbl">Used</div><div class="mon-val" id="mDskU">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">Free</div><div class="mon-val" id="mDskF">—</div></div>
          </div>
        </div>
        <div class="mon-card" style="grid-column:span 2">
          <div class="mon-title">System Info</div>
          <div class="mon-det" style="grid-template-columns:repeat(3,1fr);gap:11px">
            <div class="mon-det-i"><div class="mon-lbl">Hostname</div><div class="mon-val" id="mHost">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">Platform</div><div class="mon-val" id="mPlat">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">Arch</div><div class="mon-val" id="mArch">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">Uptime</div><div class="mon-val" id="mUp">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">Users</div><div class="mon-val" id="mUsers">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">IP</div><div class="mon-val" id="mIp">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">Net Sent</div><div class="mon-val" id="mNS">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">Net Recv</div><div class="mon-val" id="mNR">—</div></div>
            <div class="mon-det-i"><div class="mon-lbl">Boot Time</div><div class="mon-val" id="mBoot">—</div></div>
          </div>
        </div>
      </div>
    </div>

    <!-- SERVERS -->
    <div class="panel" id="panel-servers">
      <div class="topbar">
        <span style="font-size:.82rem;color:var(--tx1)">Server Manager</span>
        <span style="font-size:.72rem;color:var(--tx2);margin-left:8px" id="srvCount"></span>
        <div style="margin-left:auto">
          <button class="btn primary" onclick="openAddServer()">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>Add Server
          </button>
        </div>
      </div>
      <div class="srv-wrap">
        <div class="srv-grid" id="srvGrid"></div>
      </div>
    </div>

    <!-- SETTINGS -->
    <div class="panel" id="panel-settings">
      <div class="topbar"><span style="font-size:.82rem;color:var(--tx1)">Settings</span></div>
      <div class="sett-wrap">
        <div class="sett-info">
          <strong>Yhujin Cloudhost v2.0.0</strong><br>
          Base directory: <span id="sBaseDir">__BASEDIR__</span><br>
          Storage: <span id="sStorage">__STORAGE__</span> &bull; Max upload: 1 GB
        </div>
        <div class="sett-sec">
          <div class="sett-sec-title">File Manager</div>
          <div class="sett-row">
            <div class="sett-lbl"><strong>Show Hidden Files</strong><span>Display files starting with .</span></div>
            <div class="toggle" id="togHidden" onclick="togSett('showHidden',this)"></div>
          </div>
          <div class="sett-row">
            <div class="sett-lbl"><strong>Default View</strong><span>List or grid view on startup</span></div>
            <select class="sett-inp" id="settDefView" onchange="saveSett('defaultView',this.value)">
              <option value="list">List</option><option value="grid">Grid</option>
            </select>
          </div>
          <div class="sett-row">
            <div class="sett-lbl"><strong>Confirm Delete</strong><span>Prompt before deleting</span></div>
            <div class="toggle on" id="togConfirm" onclick="togSett('confirmDelete',this)"></div>
          </div>
        </div>
        <div class="sett-sec">
          <div class="sett-sec-title">Terminal</div>
          <div class="sett-row">
            <div class="sett-lbl"><strong>Font Size</strong><span>Terminal font size in pixels</span></div>
            <input class="sett-inp" type="number" id="settFS" value="14" min="8" max="24" onchange="applyFontSize(this.value)">
          </div>
          <div class="sett-row">
            <div class="sett-lbl"><strong>Scrollback Lines</strong><span>Terminal output history</span></div>
            <input class="sett-inp" type="number" id="settSB" value="5000" min="100" max="50000">
          </div>
        </div>
        <div class="sett-sec">
          <div class="sett-sec-title">Monitoring</div>
          <div class="sett-row">
            <div class="sett-lbl"><strong>Refresh Interval</strong><span>System info update (seconds)</span></div>
            <input class="sett-inp" type="number" id="settRefresh" value="3" min="1" max="60" onchange="applyRefresh(this.value)">
          </div>
        </div>
      </div>
    </div>

  </div>
</div>

<!-- Extract Modal -->
<div class="modal" id="extractModal">
  <div class="modal-box">
    <div class="modal-hd"><h3 id="extractTitle">Extract Archive</h3><button class="modal-cl" onclick="closeExtract()">&#x00D7;</button></div>
    <div class="modal-bd">
      <div style="font-size:.79rem;color:var(--tx1);margin-bottom:9px">Destination folder (leave empty to extract here):</div>
      <input class="sett-inp" style="width:100%" id="extractDest" type="text" placeholder="e.g. my_folder">
      <div class="extract-log" id="extractLog" style="display:none"></div>
    </div>
    <div class="modal-ft">
      <button class="btn" onclick="closeExtract()">Cancel</button>
      <button class="btn ok" id="extractBtn" onclick="doExtract()">Extract</button>
    </div>
  </div>
</div>

<!-- Upload Modal -->
<div class="modal" id="uploadModal">
  <div class="modal-box">
    <div class="modal-hd"><h3>Upload Files</h3><button class="modal-cl" onclick="closeUpload()">×</button></div>
    <div class="modal-bd">
      <div class="upload-area" id="uploadArea" onclick="document.getElementById('fileInput').click()">
        <input type="file" id="fileInput" multiple style="display:none" onchange="handleSel()">
        <svg width="30" height="30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" style="margin-bottom:9px;color:var(--tx2)"><polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/><path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/></svg>
        <div>Drag &amp; drop files or click to browse</div>
        <div style="color:var(--tx2);font-size:.74rem;margin-top:5px">Max 1 GB per file</div>
      </div>
      <div class="upload-list" id="uploadList"></div>
    </div>
    <div class="modal-ft">
      <button class="btn" onclick="closeUpload()">Cancel</button>
      <button class="btn primary" onclick="doUpload()">Upload</button>
    </div>
  </div>
</div>

<!-- Add Server Modal -->
<div class="modal" id="addServerModal">
  <div class="modal-box">
    <div class="modal-hd"><h3 id="addSrvTitle">Add Server</h3><button class="modal-cl" onclick="closeAddServer()">×</button></div>
    <div class="modal-bd">
      <div class="form-row"><label>Server Name *</label><input id="srvName" type="text" placeholder="My Production Server"></div>
      <div class="form-row"><label>Host / IP *</label><input id="srvHost" type="text" placeholder="192.168.1.100 or server.example.com"></div>
      <div class="form-row"><label>Port</label><input id="srvPort" type="number" placeholder="22" value="22"></div>
      <div class="form-row"><label>Username</label><input id="srvUser" type="text" placeholder="root"></div>
      <div class="form-row"><label>Description</label><input id="srvDesc" type="text" placeholder="Optional description"></div>
      <div class="form-row"><label>Tags (comma separated)</label><input id="srvTags" type="text" placeholder="production, web, nginx"></div>
    </div>
    <div class="modal-ft">
      <button class="btn" onclick="closeAddServer()">Cancel</button>
      <button class="btn primary" onclick="saveServer()">Save Server</button>
    </div>
  </div>
</div>

<!-- Toast -->
<div class="toast-box" id="toastBox"></div>

<!-- File Viewer -->
<div class="viewer-modal" id="viewerModal">
  <div class="viewer-topbar">
    <button class="viewer-nav-btn" id="vwrBack" onclick="viewerBack()" disabled>
      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="15 18 9 12 15 6"/></svg>
    </button>
    <button class="viewer-nav-btn" id="vwrFwd" onclick="viewerFwd()" disabled>
      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
    </button>
    <span class="viewer-title" id="vwrTitle">—</span>
    <span class="viewer-meta" id="vwrMeta"></span>
    <button class="btn ok" id="vwrEditBtn" style="margin-left:8px;flex-shrink:0;display:none" onclick="openEditorFromViewer()">Edit</button>
    <button class="btn primary" style="margin-left:4px;flex-shrink:0" onclick="vwrDownload()">Download</button>
    <button class="btn" style="margin-left:4px;flex-shrink:0" onclick="closeViewer()">Close</button>
  </div>
  <div class="viewer-body" id="vwrBody"></div>
</div>

<!-- File Editor Modal -->
<div class="viewer-modal" id="editorModal">
  <div class="viewer-topbar">
    <span class="viewer-title" id="edTitle">—</span>
    <span class="viewer-meta" id="edMeta"></span>
    <span id="edDirty" style="font-size:.72rem;color:var(--yellow);margin-left:8px;display:none">● unsaved</span>
    <div style="margin-left:auto;display:flex;gap:6px;flex-shrink:0">
      <select id="edTabSize" style="background:var(--bg2);border:1px solid var(--border);color:var(--tx1);font-family:inherit;font-size:.75rem;padding:3px 7px;border-radius:4px;outline:none" onchange="applyTabSize()">
        <option value="2">2 spaces</option>
        <option value="4" selected>4 spaces</option>
        <option value="8">8 spaces</option>
      </select>
      <button class="btn" onclick="edWordWrap()">Wrap</button>
      <button class="btn ok" onclick="saveFile()" id="edSaveBtn">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
        Save
      </button>
      <button class="btn" onclick="closeEditor()">Close</button>
    </div>
  </div>
  <div style="flex:1;display:flex;flex-direction:column;overflow:hidden;position:relative">
    <div id="edLineBar" style="position:absolute;bottom:10px;right:16px;font-size:.72rem;color:var(--tx2);z-index:10;pointer-events:none">Ln 1, Col 1</div>
    <textarea id="edTextarea" spellcheck="false" autocorrect="off" autocapitalize="off"
      style="flex:1;width:100%;height:100%;background:#0c0c10;color:#dde1e8;border:none;outline:none;
             padding:16px;font-family:'JetBrains Mono',monospace;font-size:.82rem;line-height:1.6;
             resize:none;tab-size:4;white-space:pre;overflow-wrap:normal;overflow-x:auto"
      oninput="edOnInput()" onkeydown="edKeyDown(event)" onclick="edUpdatePos()" onkeyup="edUpdatePos()"
    ></textarea>
  </div>
</div>

<!-- Terminal Context Menu -->
<div class="ctx-menu" id="termCtxMenu">
  <div class="ctx-item" onclick="ctxCopy()">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
    Copy
  </div>
  <div class="ctx-item" onclick="ctxPaste()">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1"/></svg>
    Paste
  </div>
  <div class="ctx-sep"></div>
  <div class="ctx-item" onclick="ctxCopyAll()">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
    Copy All Output
  </div>
  <div class="ctx-item" onclick="ctxClear()">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/></svg>
    Clear Terminal
  </div>
  <div class="ctx-sep"></div>
  <div class="ctx-item" onclick="ctxNewTab()">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
    New Tab
  </div>
</div>

<script>
// ── API FETCH (auto auth redirect) ───────────────────────────────────────────
function apiFetch(url, opts={}, _retry){
  return fetch(url, opts).then(r=>{
    if(r.status === 401){ window.location.href='/login'; throw new Error('Unauthorized'); }
    return r;
  });
}

// ── STATE ──
let curPanel='files', curPath='/', curView='list';
let allFiles=[], selFiles=[], pathHist=['/'], histIdx=0;
let procData=[], procSort='cpu_percent', procAsc=false;
let procTimer=null, monTimer=null;
let settings={showHidden:false,confirmDelete:true,defaultView:'list',termFontSize:14};
let upFiles=[];
let vwrHist=[], vwrIdx=-1, vwrCurPath=null;
let servers=[];

// ── COOKIE HELPERS ──
function setCookie(name,val,days=365){
  const d=new Date();d.setTime(d.getTime()+(days*86400000));
  document.cookie=name+'='+encodeURIComponent(val)+';expires='+d.toUTCString()+';path=/;SameSite=Lax';
}
function getCookie(name){
  const ca=document.cookie.split(';');
  for(let c of ca){c=c.trim();if(c.startsWith(name+'='))return decodeURIComponent(c.substring(name.length+1));}
  return null;
}
function saveToCookie(key,obj){try{setCookie(key,JSON.stringify(obj));}catch(e){}}
function loadFromCookie(key,def){try{const v=getCookie(key);return v?JSON.parse(v):def;}catch(e){return def;}}

function doLogout(){
  window.location.href='/logout';
}

// ── INIT ──
document.addEventListener('DOMContentLoaded',()=>{
  // Dev mode badge
  if(DEV_MODE){
    document.getElementById('devBadge').style.display='inline';
    document.title = '[DEV] Yhujin Cloudhost';
  }
  loadSetts();
  const lastPanel=getCookie('yhujin_panel')||'files';
  loadFiles('/');
  restoreTermSessions();
  updateMon();
  monTimer=setInterval(updateMon,3000);
  loadProcs();
  procTimer=setInterval(loadProcs,4000);
  const ua=document.getElementById('uploadArea');
  ua.addEventListener('dragover',e=>{e.preventDefault();ua.classList.add('dragover')});
  ua.addEventListener('dragleave',()=>ua.classList.remove('dragover'));
  ua.addEventListener('drop',e=>{e.preventDefault();ua.classList.remove('dragover');upFiles=Array.from(e.dataTransfer.files);renderUpList()});
  applySetts();
  loadServers();
  setTimeout(()=>switchPanel(lastPanel),50);
  document.addEventListener('click',()=>hideCtxMenu());
  // TIDAK close session saat halaman ditutup — biarkan proses tetap berjalan
});

// ── PANEL SWITCH ──
function switchPanel(p){
  curPanel=p;
  setCookie('yhujin_panel',p);
  document.querySelectorAll('.tab-item').forEach(t=>t.classList.toggle('active',t.dataset.panel===p));
  document.querySelectorAll('.panel').forEach(x=>x.classList.toggle('active',x.id==='panel-'+p));
  document.querySelectorAll('.sb-icon[id^="sb-"]').forEach(s=>s.classList.toggle('active',s.id==='sb-'+p));
  if(p==='terminal'){setTimeout(()=>{const t=activeTerm();if(t){t.fit.fit();t.term.focus();}},80);}
}

// ── FILES ──
function loadFiles(path){
  curPath=path;
  document.getElementById('fileList').innerHTML='<div class="empty-state">Loading...</div>';
  apiFetch('/api/files'+path)
    .then(r=>r.json()).then(data=>{
      if(data.error){toast(data.error,'error');return}
      allFiles=data;
      if(!settings.showHidden) allFiles=allFiles.filter(f=>!f.is_hidden);
      renderBC(path);
      histBtns();
      renderFiles(allFiles);
      selFiles=[];selToolbar();
    }).catch(()=>{
      document.getElementById('fileList').innerHTML='<div class="empty-state">Error loading files</div>';
    });
}
function goTo(sub){
  let np;
  if(sub==='..'){
    const parts=curPath.replace(/\/+$/,'').split('/').filter(Boolean);
    parts.pop(); np='/'+parts.join('/'); if(!np)np='/';
  } else {
    np=(curPath.endsWith('/')?curPath:curPath+'/')+sub;
  }
  pathHist=pathHist.slice(0,histIdx+1);
  pathHist.push(np);
  histIdx=pathHist.length-1;
  loadFiles(np);
}
function histBack(){if(histIdx>0){histIdx--;loadFiles(pathHist[histIdx]);}}
function histFwd(){if(histIdx<pathHist.length-1){histIdx++;loadFiles(pathHist[histIdx]);}}
function histBtns(){
  document.getElementById('btnBack').disabled=histIdx<=0;
  document.getElementById('btnFwd').disabled=histIdx>=pathHist.length-1;
}
function renderBC(path){
  const parts=path.split('/').filter(Boolean);
  let html='<span class="bc-item" onclick="bcNav(\'/\')">root</span>';
  let built='';
  parts.forEach((p,i)=>{
    built+='/'+p;
    const snap=built;
    const last=i===parts.length-1;
    html+='<span class="bc-sep">/</span><span class="bc-item'+(last?' last':'')+'" onclick="bcNav(\''+snap.replace(/'/g,"\\'")+'\')">'
          +escH(p)+'</span>';
  });
  document.getElementById('breadcrumb').innerHTML=html;
}
function bcNav(path){
  pathHist=pathHist.slice(0,histIdx+1);pathHist.push(path);histIdx=pathHist.length-1;loadFiles(path);
}
function refreshFiles(){loadFiles(curPath);}
function renderFiles(files){
  const hdr=document.getElementById('fl-header');
  if(!files||!files.length){hdr.style.display='none';document.getElementById('fileList').innerHTML='<div class="empty-state">Empty directory</div>';return;}
  hdr.style.display=curView==='list'?'grid':'none';
  document.getElementById('fileList').innerHTML=curView==='list'
    ?files.map(fiRow).join('')
    :'<div class="file-grid">'+files.map(fiGrid).join('')+'</div>';
}
function fiIcon(isDir,ext,big){
  const s=big?28:14;
  if(isDir)return` + "`" + `<svg width="${s}" height="${s}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>` + "`" + `;
  if(['.jpg','.jpeg','.png','.gif','.webp','.svg','.bmp'].includes(ext))
    return` + "`" + `<svg width="${s}" height="${s}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>` + "`" + `;
  if(['.py','.js','.ts','.html','.css','.json','.sh','.php','.rb','.go','.rs','.c','.cpp','.java'].includes(ext))
    return` + "`" + `<svg width="${s}" height="${s}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>` + "`" + `;
  if(['.zip','.tar','.gz','.bz2','.xz','.7z','.rar'].includes(ext))
    return` + "`" + `<svg width="${s}" height="${s}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>` + "`" + `;
  return` + "`" + `<svg width="${s}" height="${s}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><polyline points="13 2 13 9 20 9"/></svg>` + "`" + `;
}
function fiRow(f){
  const sel=selFiles.includes(f.name),hid=f.is_hidden;
  const dt=new Date(f.modified).toLocaleString('en-GB',{day:'2-digit',month:'short',hour:'2-digit',minute:'2-digit'});
  return` + "`" + `<div class="fi${sel?' sel':''}${hid?' hid':''}" onclick="rowClick(event,'${eq(f.name)}')" ondblclick="openFile('${eq(f.name)}')">
    <div class="fi-cb" onclick="event.stopPropagation()"><input type="checkbox" ${sel?'checked':''} onchange="cbChange(event,'${eq(f.name)}')"></div>
    <div class="fi-icon${f.is_dir?' d':''}">${fiIcon(f.is_dir,f.extension,false)}</div>
    <div class="fi-name${f.is_dir?' d':''}${hid?' hid':''}" title="${eh(f.name)}">${eh(f.name)}${hid?'<span style="font-size:.62rem;color:var(--tx2)"> [hidden]</span>':''}</div>
    <div class="fi-size">${f.is_dir?'—':f.size_formatted}</div>
    <div class="fi-own" title="${f.owner}:${f.group}" style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap">${f.owner}</div>
    <div class="fi-mod" title="${f.permissions}">${dt}</div>
    <div class="fi-act">
      ${f.is_dir
        ?` + "`" + `<button class="ib" onclick="event.stopPropagation();goTo('${eq(f.name)}')">Open</button>` + "`" + `
        :` + "`" + `<button class="ib" onclick="event.stopPropagation();openViewer('${eq(f.name)}')">View</button>${isTextFile(f.extension,f.name)?` + "`" + `<button class="ib ed" onclick="event.stopPropagation();openEditor('${eq(f.name)}')">Edit</button>` + "`" + `:''}
         <button class="ib" onclick="event.stopPropagation();dlFile('${eq(f.name)}')">DL</button>` + "`" + `
      }
      ${f.is_archive?` + "`" + `<button class="ib ex" onclick="event.stopPropagation();openExtractModal(['${eq(f.name)}'])">Ext</button>` + "`" + `:''}
      <button class="ib del" onclick="event.stopPropagation();delOne('${eq(f.name)}')">Del</button>
    </div>
  </div>` + "`" + `;
}
function fiGrid(f){
  const sel=selFiles.includes(f.name),hid=f.is_hidden;
  const dbl=f.is_dir?` + "`" + `goTo('${eq(f.name)}')` + "`" + `:` + "`" + `openViewer('${eq(f.name)}')` + "`" + `;
  return` + "`" + `<div class="fgi${sel?' sel':''}${hid?' hid':''}" onclick="rowClick(event,'${eq(f.name)}')" ondblclick="${dbl}">
    <div class="fgi-cb" onclick="event.stopPropagation()"><input type="checkbox" ${sel?'checked':''} onchange="cbChange(event,'${eq(f.name)}')"></div>
    <div class="fgi-icon${!f.is_dir?' f':''}">${fiIcon(f.is_dir,f.extension,true)}</div>
    <div class="fgi-name" title="${eh(f.name)}">${eh(f.name)}</div>
    <div class="fgi-size">${f.is_dir?'—':f.size_formatted}</div>
  </div>` + "`" + `;
}
function rowClick(e,name){
  if(e.target.tagName==='INPUT'||e.target.tagName==='BUTTON') return;
  if(e.ctrlKey||e.metaKey){
    selFiles.includes(name)?selFiles=selFiles.filter(x=>x!==name):selFiles.push(name);
  } else if(e.shiftKey){
    const idx=allFiles.findIndex(f=>f.name===name);
    const lastIdx=selFiles.length?allFiles.findIndex(f=>f.name===selFiles[selFiles.length-1]):-1;
    if(lastIdx>=0){
      const lo=Math.min(idx,lastIdx),hi=Math.max(idx,lastIdx);
      selFiles=allFiles.slice(lo,hi+1).map(f=>f.name);
    } else selFiles=[name];
  } else {
    selFiles=selFiles.includes(name)&&selFiles.length===1?[]:[name];
  }
  renderFiles(allFiles);selToolbar();
}
function cbChange(e,name){
  e.stopPropagation();
  if(e.target.checked){if(!selFiles.includes(name))selFiles.push(name);}
  else selFiles=selFiles.filter(x=>x!==name);
  renderFiles(allFiles);selToolbar();
}
function toggleSelAll(cb){
  selFiles=cb.checked?allFiles.map(f=>f.name):[];
  renderFiles(allFiles);selToolbar();
}
function selToolbar(){
  const n=selFiles.length,h=n>0;
  const info=document.getElementById('selInfo');
  if(info){info.style.display=h?'':'none';if(h)info.textContent=n+' selected';}
  document.getElementById('btnDel').style.display=h?'':'none';
  document.getElementById('btnRename').style.display=n===1?'':'none';
  document.getElementById('btnHide').style.display=h?'':'none';
  document.getElementById('btnCut').style.display=h?'':'none';
  document.getElementById('btnCopy').style.display=h?'':'none';
  const hasArc=selFiles.some(nm=>allFiles.find(f=>f.name===nm)?.is_archive);
  document.getElementById('btnExtract').style.display=hasArc?'':'none';
  document.getElementById('selAllCb').indeterminate=h&&n<allFiles.length;
  document.getElementById('selAllCb').checked=h&&n===allFiles.length;
  updatePasteBtn();
}

// ── CLIPBOARD ──
let clipboard = null; // {mode:'cut'|'copy', srcs:[], fromPath:''}

function updatePasteBtn(){
  document.getElementById('btnPaste').style.display=clipboard?'':'none';
}

function cutSelected(){
  if(!selFiles.length)return;
  clipboard={mode:'cut', srcs:[...selFiles], fromPath:curPath};
  toast('Cut '+selFiles.length+' item — navigasi ke folder tujuan lalu Paste','info');
  selFiles=[];renderFiles(allFiles);selToolbar();
}

function copySelected(){
  if(!selFiles.length)return;
  clipboard={mode:'copy', srcs:[...selFiles], fromPath:curPath};
  toast('Copy '+selFiles.length+' item — navigasi ke folder tujuan lalu Paste','info');
  selFiles=[];renderFiles(allFiles);selToolbar();
}

function pasteHere(){
  if(!clipboard)return;
  const srcPaths=clipboard.srcs.map(n=>{
    const base=clipboard.fromPath.endsWith('/')?clipboard.fromPath:clipboard.fromPath+'/';
    return base+n;
  });
  const endpoint=clipboard.mode==='cut'?'/api/move':'/api/copy';
  apiFetch(endpoint,{
    method:'POST',
    headers:{'Content-Type':'application/json'},
    body:JSON.stringify({src:srcPaths, dest:curPath})
  }).then(r=>r.json()).then(d=>{
    if(d.error){toast(d.error,'error');return;}
    toast(d.message||'Selesai','success');
    if(clipboard.mode==='cut') clipboard=null;
    updatePasteBtn();
    loadFiles(curPath);
  });
}
function openFile(name){
  const f=allFiles.find(x=>x.name===name);
  if(!f)return;
  if(f.is_dir){ goTo(name); }
  else { openViewer(name); }
}
function dlFile(name){
  const p=(curPath.endsWith('/')?curPath:curPath+'/')+name;
  const a=document.createElement('a');a.href='/api/download'+p;a.download=name;a.click();
}
function delOne(name){
  if(settings.confirmDelete&&!confirm('Delete '+name+'?'))return;
  const p=(curPath.endsWith('/')?curPath:curPath+'/')+name;
  apiFetch('/api/delete'+p,{method:'DELETE'}).then(r=>r.json()).then(d=>{
    if(d.error)toast(d.error,'error');else{toast('Deleted','success');loadFiles(curPath);}
  });
}
function deleteSelected(){
  if(!selFiles.length)return;
  if(settings.confirmDelete&&!confirm('Delete '+selFiles.length+' item(s)?'))return;
  Promise.all(selFiles.map(n=>{
    const p=(curPath.endsWith('/')?curPath:curPath+'/')+n;
    return apiFetch('/api/delete'+p,{method:'DELETE'});
  })).then(()=>{toast('Deleted','success');selFiles=[];loadFiles(curPath);});
}
function renameSelected(){
  if(selFiles.length!==1)return;
  const oldN=selFiles[0], newN=prompt('Rename to:',oldN);
  if(!newN||newN===oldN)return;
  apiFetch('/api/rename',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({path:curPath,old_name:oldN,new_name:newN})})
    .then(r=>r.json()).then(d=>{if(d.error)toast(d.error,'error');else{toast('Renamed','success');selFiles=[];loadFiles(curPath);}});
}
function newFile(){
  const n=prompt('File name:');if(!n)return;
  apiFetch('/api/create/file',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({path:curPath,name:n})})
    .then(r=>r.json()).then(d=>{if(d.error)toast(d.error,'error');else{toast('Created','success');loadFiles(curPath);}});
}
function newFolder(){
  const n=prompt('Folder name:');if(!n)return;
  apiFetch('/api/create/folder',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({path:curPath,name:n})})
    .then(r=>r.json()).then(d=>{if(d.error)toast(d.error,'error');else{toast('Created','success');loadFiles(curPath);}});
}
function setView(v,btn){
  curView=v;
  document.querySelectorAll('.vbtn').forEach(b=>b.classList.remove('active'));
  btn.classList.add('active');
  renderFiles(allFiles);
}
function filterFiles(){
  const q=document.getElementById('searchInput').value.toLowerCase();
  renderFiles(q?allFiles.filter(f=>f.name.toLowerCase().includes(q)):allFiles);
}

function isTextFile(ext,name){
  const exts=['.txt','.md','.log','.csv','.ini','.cfg','.conf','.env',
    '.py','.js','.ts','.jsx','.tsx','.html','.htm','.css','.scss','.less',
    '.json','.yaml','.yml','.toml','.xml','.sql','.sh','.bash','.zsh',
    '.php','.rb','.go','.rs','.c','.cpp','.h','.hpp','.java','.kt','.swift',
    '.r','.m','.lua','.pl','.pm','.awk','.sed','.vue','.svelte','.tf','.hcl',
    '.gitignore','.htaccess','.nginx','.dockerfile'];
  const noExt=['makefile','dockerfile','readme','license','changelog','authors','todo'];
  return exts.includes(ext)||noExt.includes((name||'').toLowerCase());
}

// ── FILE EDITOR ──
let edFilePath=null, edWordWrapOn=false, edOriginal='';

function openEditor(name){
  const fp=filePath(name);
  edFilePath=fp;
  document.getElementById('edTitle').textContent=name;
  document.getElementById('edMeta').textContent='';
  document.getElementById('edDirty').style.display='none';
  document.getElementById('editorModal').classList.add('active');
  const ta=document.getElementById('edTextarea');
  ta.value='Loading...';
  ta.disabled=true;
  apiFetch('/api/content'+fp).then(r=>r.json()).then(d=>{
    if(d.error){ta.value='Error: '+d.error;return;}
    edOriginal=d.content;
    ta.value=d.content;
    ta.disabled=false;
    ta.focus();
    document.getElementById('edMeta').textContent=d.lines+' lines · '+d.size;
    edUpdatePos();
  }).catch(e=>{ta.value='Error loading file';});
}

function openEditorFromViewer(){
  const name=document.getElementById('vwrTitle').textContent;
  closeViewer();
  openEditor(name);
}

function closeEditor(){
  if(document.getElementById('edDirty').style.display!=='none'){
    if(!confirm('Ada perubahan yang belum disimpan. Tutup tanpa menyimpan?'))return;
  }
  document.getElementById('editorModal').classList.remove('active');
  document.getElementById('edTextarea').value='';
  edFilePath=null;edOriginal='';
}

function saveFile(){
  if(!edFilePath)return;
  const content=document.getElementById('edTextarea').value;
  const btn=document.getElementById('edSaveBtn');
  btn.textContent='Saving...';btn.disabled=true;
  apiFetch('/api/save',{method:'POST',headers:{'Content-Type':'application/json'},
    body:JSON.stringify({path:edFilePath,content:content})
  }).then(r=>r.json()).then(d=>{
    btn.innerHTML=` + "`" + `<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg> Save` + "`" + `;
    btn.disabled=false;
    if(d.error){toast(d.error,'error');return;}
    edOriginal=content;
    document.getElementById('edDirty').style.display='none';
    toast('Saved: '+edFilePath.split('/').pop(),'success');
    loadFiles(curPath);
  }).catch(()=>{btn.textContent='Save';btn.disabled=false;toast('Save failed','error');});
}

function edOnInput(){
  const dirty=document.getElementById('edTextarea').value!==edOriginal;
  document.getElementById('edDirty').style.display=dirty?'':'none';
  edUpdatePos();
}

function edUpdatePos(){
  const ta=document.getElementById('edTextarea');
  const val=ta.value.substring(0,ta.selectionStart);
  const ln=val.split('\n').length;
  const col=val.split('\n').pop().length+1;
  document.getElementById('edLineBar').textContent=` + "`" + `Ln ${ln}, Col ${col}` + "`" + `;
}

function edKeyDown(e){
  if((e.ctrlKey||e.metaKey)&&e.key==='s'){e.preventDefault();saveFile();return;}
  if(e.key==='Tab'){
    e.preventDefault();
    const ta=e.target,s=ta.selectionStart,en=ta.selectionEnd;
    const spaces=' '.repeat(parseInt(document.getElementById('edTabSize').value)||4);
    ta.value=ta.value.substring(0,s)+spaces+ta.value.substring(en);
    ta.selectionStart=ta.selectionEnd=s+spaces.length;
    edOnInput();
  }
}

function edWordWrap(){
  edWordWrapOn=!edWordWrapOn;
  const ta=document.getElementById('edTextarea');
  ta.style.whiteSpace=edWordWrapOn?'pre-wrap':'pre';
  ta.style.overflowWrap=edWordWrapOn?'break-word':'normal';
  ta.style.overflowX=edWordWrapOn?'hidden':'auto';
}

function applyTabSize(){
  const v=document.getElementById('edTabSize').value;
  document.getElementById('edTextarea').style.tabSize=v;
}

document.addEventListener('keydown',e=>{
  if(e.key==='Escape'&&document.getElementById('editorModal').classList.contains('active'))closeEditor();
});

// ── FILE VIEWER ──
const TEXT_EXTS=['.txt','.md','.log','.csv','.ini','.cfg','.conf','.env',
  '.py','.js','.ts','.jsx','.tsx','.html','.htm','.css','.scss','.less',
  '.json','.yaml','.yml','.toml','.xml','.sql','.sh','.bash','.zsh',
  '.php','.rb','.go','.rs','.c','.cpp','.h','.hpp','.java','.kt','.swift',
  '.r','.m','.lua','.pl','.pm','.awk','.sed','.dockerfile','makefile',
  '.gitignore','.htaccess','.nginx','.vue','.svelte','.tf','.hcl'];
const IMG_EXTS=['.jpg','.jpeg','.png','.gif','.webp','.svg','.bmp','.ico','.avif'];
const VID_EXTS=['.mp4','.webm','.ogv','.mov','.avi','.mkv'];
const AUD_EXTS=['.mp3','.wav','.ogg','.flac','.aac','.m4a'];

function filePath(name){return(curPath.endsWith('/')?curPath:curPath+'/')+name;}
function openViewer(name){
  const fp=filePath(name);
  vwrHist=vwrHist.slice(0,vwrIdx+1);
  vwrHist.push({name,path:fp,dirPath:curPath});
  vwrIdx=vwrHist.length-1;
  _showViewer(name,fp);
}
function _showViewer(name,fp){
  vwrCurPath=fp;
  const ext=name.includes('.')?'.'+name.split('.').pop().toLowerCase():'';
  document.getElementById('vwrTitle').textContent=name;
  document.getElementById('vwrMeta').textContent='';
  document.getElementById('viewerModal').classList.add('active');
  document.getElementById('vwrEditBtn').style.display=isTextFile(ext,name)?'':'none';
  vwrNavBtns();
  const body=document.getElementById('vwrBody');
  body.innerHTML='<div style="flex:1;display:flex;align-items:center;justify-content:center;color:var(--tx2);font-size:.85rem">Loading...</div>';
  if(IMG_EXTS.includes(ext)){
    body.innerHTML=` + "`" + `<div class="viewer-img"><img src="/api/serve${fp}" alt="${eh(name)}" onload="document.getElementById('vwrMeta').textContent=this.naturalWidth+'×'+this.naturalHeight"></div>` + "`" + `;
  } else if(VID_EXTS.includes(ext)){
    body.innerHTML=` + "`" + `<div class="viewer-media"><video controls autoplay src="/api/serve${fp}" style="max-width:100%;max-height:calc(100vh - 120px)"></video></div>` + "`" + `;
  } else if(AUD_EXTS.includes(ext)){
    body.innerHTML=` + "`" + `<div class="viewer-media" style="flex-direction:column;gap:16px">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="var(--tx2)" stroke-width="1.5"><path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/></svg>
      <div style="color:var(--tx0);font-size:.9rem">${eh(name)}</div>
      <audio controls autoplay src="/api/serve${fp}" style="width:400px;max-width:90vw"></audio>
    </div>` + "`" + `;
  } else if(ext==='.pdf'){
    body.innerHTML=` + "`" + `<div class="viewer-pdf"><iframe src="/api/serve${fp}"></iframe></div>` + "`" + `;
  } else if(TEXT_EXTS.includes(ext)||ext===''||isLikelyText(name)){
    apiFetch('/api/content'+fp).then(r=>r.json()).then(d=>{
      if(d.error){body.innerHTML=` + "`" + `<div class="viewer-unsupported"><div>${eh(d.error)}</div></div>` + "`" + `;return;}
      document.getElementById('vwrMeta').textContent=d.lines+' lines · '+d.size;
      const escaped=d.content.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
      const lines=escaped.split('\n').map(l=>` + "`" + `<span class="ln">${l}</span>` + "`" + `).join('\n');
      body.innerHTML=` + "`" + `<pre class="viewer-code line-nums"><code class="line-nums">${lines}</code></pre>` + "`" + `;
    });
  } else {
    body.innerHTML=` + "`" + `<div class="viewer-unsupported">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><polyline points="13 2 13 9 20 9"/></svg>
      <div>Preview not available for this file type</div>
      <div style="font-size:.76rem;color:var(--tx2)">${eh(name)}</div>
      <button class="btn primary" onclick="dlFile('${eq(name)}')">Download File</button>
    </div>` + "`" + `;
  }
}
function isLikelyText(name){
  const noExt=['makefile','dockerfile','readme','license','changelog','authors','todo','gemfile','rakefile','procfile'];
  return noExt.includes(name.toLowerCase());
}
function closeViewer(){document.getElementById('viewerModal').classList.remove('active');document.getElementById('vwrBody').innerHTML='';vwrCurPath=null;}
function vwrDownload(){if(!vwrCurPath)return;const name=vwrCurPath.split('/').pop();const a=document.createElement('a');a.href='/api/download'+vwrCurPath;a.download=name;a.click();}
function viewerBack(){if(vwrIdx>0){vwrIdx--;const h=vwrHist[vwrIdx];if(h.dirPath!==curPath)curPath=h.dirPath;_showViewer(h.name,h.path);}}
function viewerFwd(){if(vwrIdx<vwrHist.length-1){vwrIdx++;const h=vwrHist[vwrIdx];if(h.dirPath!==curPath)curPath=h.dirPath;_showViewer(h.name,h.path);}}
function vwrNavBtns(){document.getElementById('vwrBack').disabled=vwrIdx<=0;document.getElementById('vwrFwd').disabled=vwrIdx>=vwrHist.length-1;}
document.addEventListener('keydown',e=>{if(e.key==='Escape'&&document.getElementById('viewerModal').classList.contains('active'))closeViewer();});

// ── HIDE/UNHIDE ──
function toggleHideOne(name){
  apiFetch('/api/hide',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({path:curPath,name:name})})
    .then(r=>r.json()).then(d=>{d.error?toast(d.error,'error'):(toast(d.message,'success'),selFiles=[],loadFiles(curPath));});
}
function toggleHideSelected(){
  if(!selFiles.length)return;
  Promise.all(selFiles.map(n=>apiFetch('/api/hide',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({path:curPath,name:n})}).then(r=>r.json())))
    .then(res=>{const ok=res.filter(r=>!r.error).length;toast('Toggled hidden on '+ok+' item(s)','success');selFiles=[];loadFiles(curPath);});
}

// ── EXTRACT ──
let extractFiles=[];
function openExtractModal(names){
  extractFiles=names||selFiles.filter(n=>allFiles.find(f=>f.name===n)?.is_archive);
  if(!extractFiles.length){toast('No archive selected','warn');return;}
  document.getElementById('extractTitle').textContent='Extract: '+extractFiles.join(', ');
  document.getElementById('extractDest').value='';
  const log=document.getElementById('extractLog');log.style.display='none';log.textContent='';
  document.getElementById('extractBtn').disabled=false;
  document.getElementById('extractModal').classList.add('active');
}
function closeExtract(){document.getElementById('extractModal').classList.remove('active');}
function extractSelected(){openExtractModal(null);}
function doExtract(){
  const dest=document.getElementById('extractDest').value.trim();
  const log=document.getElementById('extractLog');log.style.display='block';log.textContent='Extracting…\n';
  document.getElementById('extractBtn').disabled=true;
  Promise.all(extractFiles.map(name=>
    apiFetch('/api/extract',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({path:curPath,name:name,dest:dest||null})})
      .then(r=>r.json())
  )).then(results=>{
    results.forEach((r,i)=>{log.textContent+=extractFiles[i]+': '+(r.error?'ERROR - '+r.error:'OK → '+r.dest)+'\n';});
    log.textContent+='Done.';document.getElementById('extractBtn').disabled=false;loadFiles(curPath);
  });
}

// ── UPLOAD ──
function openUpload(){document.getElementById('uploadModal').classList.add('active');}
function closeUpload(){document.getElementById('uploadModal').classList.remove('active');upFiles=[];document.getElementById('uploadList').innerHTML='';}
function handleSel(){upFiles=Array.from(document.getElementById('fileInput').files);renderUpList();}
function renderUpList(){
  document.getElementById('uploadList').innerHTML=upFiles.map(f=>` + "`" + `<div class="upload-item"><span>${eh(f.name)}</span><span>${fmtBytes(f.size)}</span></div>` + "`" + `).join('');
}
function doUpload(){
  if(!upFiles.length)return;
  const fd=new FormData();
  upFiles.forEach(f=>fd.append('files',f));
  fd.append('path',curPath);
  toast('Uploading '+upFiles.length+' file(s)...','info');
  apiFetch('/api/upload',{method:'POST',body:fd}).then(r=>r.json()).then(d=>{closeUpload();toast(d.message||'Done','success');loadFiles(curPath);});
}

// ── TERMINAL MULTI-TAB ──

function restoreTermSessions(){
  // Cek session lama dari cookie, coba reconnect
  const savedSids=loadFromCookie('yhujin_terms',[]);
  if(!savedSids.length){newTermTab();return;}
  // Fetch daftar session aktif dari server
  apiFetch('/api/terminal/list').then(r=>r.json()).then(list=>{
    const activeSids=list.map(s=>s.session_id);
    const toRestore=savedSids.filter(sid=>activeSids.includes(sid));
    if(!toRestore.length){
      // Semua session sudah mati (server di-restart), buat baru
      saveToCookie('yhujin_terms',[]);
      newTermTab();
      return;
    }
    // Reconnect ke session yang masih hidup
    toRestore.forEach((sid,i)=>attachTermTab(sid,'Terminal '+(i+1)));
    // Simpan hanya yang masih hidup
    saveToCookie('yhujin_terms',toRestore);
    // Kalau ada session baru yang belum tersimpan, biarkan
  }).catch(()=>{newTermTab();});
}

function attachTermTab(sid,label){
  const id='t'+(++termCtr);
  if(!label)label='Terminal '+termCtr;
  const instEl=document.createElement('div');instEl.className='term-inst';instEl.id='ti-'+id;
  document.getElementById('termInstances').appendChild(instEl);
  const tabEl=document.createElement('div');tabEl.className='term-tab';tabEl.id='tt-'+id;tabEl.dataset.tid=id;
  tabEl.innerHTML='<span class="tdot"></span><span>'+label+'</span>'
    +'<span class="tcl" onclick="event.stopPropagation();closeTermTab(\''+id+'\')" title="Detach">&minus;</span>'
    +'<span class="tcl" style="color:var(--red)" onclick="event.stopPropagation();if(confirm(\'Matikan proses ini?\'))killTermTab(\''+id+'\')" title="Kill">&#x00D7;</span>';
  tabEl.addEventListener('click',()=>switchTermTab(id));
  const tnew=document.querySelector('.tnew');
  document.getElementById('termTabBarTabs').insertBefore(tabEl,tnew);

  const term=new Terminal({
    fontFamily:"'JetBrains Mono','Fira Code',monospace",
    fontSize:settings.termFontSize||14,lineHeight:1.3,theme:TERM_THEME,
    scrollback:5000,cursorBlink:true,allowTransparency:true,bellStyle:'none',
    convertEol:false,macOptionIsMeta:true,rightClickSelectsWord:false
  });
  const fit=new FitAddon.FitAddon();
  const links=new WebLinksAddon.WebLinksAddon();
  term.loadAddon(fit);term.loadAddon(links);term.open(instEl);fit.fit();

  const obj={id,label,term,fit,sid,poll:null,el:instEl,tabEl};
  termTabs.push(obj);
  switchTermTab(id);

  instEl.addEventListener('contextmenu',(e)=>{e.preventDefault();e.stopPropagation();showCtxMenu(e.clientX,e.clientY,obj);});

  // Ambil scrollback lalu mulai polling
  apiFetch('/api/terminal/attach/'+sid).then(r=>r.json()).then(d=>{
    tabEl.querySelector('.tdot').style.background=d.running?'var(--green)':'var(--red)';
    if(d.scrollback&&d.scrollback.length){
      try{
        const bin=atob(d.scrollback);
        const b=new Uint8Array(bin.length);
        for(let i=0;i<bin.length;i++)b[i]=bin.charCodeAt(i);
        term.write(b);
      }catch(e){}
    }
    obj.poll=setInterval(()=>pollTerm(obj),60);
    setTimeout(()=>{fit.fit();term.focus();},100);
  }).catch(()=>{
    term.write('\r\n\x1b[33mSession tidak ditemukan, membuat session baru...\x1b[0m\r\n');
    startNewSession(obj,tabEl);
  });

  term.onData(data=>{
    if(!obj.sid)return;
    apiFetch('/api/terminal/input/'+obj.sid,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({input:toB64(data)})});
  });
  window.addEventListener('resize',()=>{if(obj.id===activeTermId)fitTermTab(obj);});
}

function startNewSession(obj,tabEl){
  apiFetch('/api/terminal/start',{method:'POST'}).then(r=>r.json()).then(d=>{
    obj.sid=d.session_id;
    tabEl.querySelector('.tdot').style.background='var(--green)';
    // Simpan ke cookie
    const saved=loadFromCookie('yhujin_terms',[]);
    if(!saved.includes(d.session_id)){saved.push(d.session_id);saveToCookie('yhujin_terms',saved);}
    obj.poll=setInterval(()=>pollTerm(obj),60);
    if(obj.id===activeTermId)obj.term.focus();
  }).catch(()=>obj.term.write('\r\n\x1b[31mFailed to connect\x1b[0m\r\n'));
}

const TERM_THEME={
  background:'#0c0c10',foreground:'#dde1e8',cursor:'#4a9eff',cursorAccent:'#0c0c10',
  selectionBackground:'rgba(74,158,255,.3)',black:'#1a1a1f',red:'#f85149',green:'#3fb950',
  yellow:'#d29922',blue:'#4a9eff',magenta:'#bc8cff',cyan:'#39d0d8',white:'#dde1e8',
  brightBlack:'#5a6070',brightRed:'#ff7b72',brightGreen:'#56d364',brightYellow:'#e3b341',
  brightBlue:'#79c0ff',brightMagenta:'#d2a8ff',brightCyan:'#56d4dd',brightWhite:'#ffffff'
};
let termTabs=[],activeTermId=null,termCtr=0;
function activeTerm(){return termTabs.find(t=>t.id===activeTermId);}

function newTermTab(){
  const id='t'+(++termCtr),label='Terminal '+termCtr;
  const instEl=document.createElement('div');instEl.className='term-inst';instEl.id='ti-'+id;
  document.getElementById('termInstances').appendChild(instEl);
  const tabEl=document.createElement('div');tabEl.className='term-tab';tabEl.id='tt-'+id;tabEl.dataset.tid=id;
  tabEl.innerHTML='<span class="tdot"></span><span>'+label+'</span>'
    +'<span class="tcl" onclick="event.stopPropagation();closeTermTab(\''+id+'\')" title="Detach (proses tetap jalan)">&minus;</span>'
    +'<span class="tcl" style="color:var(--red)" onclick="event.stopPropagation();if(confirm(\'Matikan proses ini?\'))killTermTab(\''+id+'\')" title="Kill (matikan proses)">&#x00D7;</span>';
  tabEl.addEventListener('click',()=>switchTermTab(id));
  const tnew=document.querySelector('.tnew');
  document.getElementById('termTabBarTabs').insertBefore(tabEl,tnew);

  const term=new Terminal({
    fontFamily:"'JetBrains Mono','Fira Code',monospace",
    fontSize:settings.termFontSize||14,
    lineHeight:1.3,
    theme:TERM_THEME,
    scrollback:5000,
    cursorBlink:true,
    allowTransparency:true,
    bellStyle:'none',
    convertEol:false,
    macOptionIsMeta:true,
    rightClickSelectsWord:false
  });
  const fit=new FitAddon.FitAddon();
  const links=new WebLinksAddon.WebLinksAddon();
  term.loadAddon(fit);term.loadAddon(links);term.open(instEl);fit.fit();

  const obj={id,label,term,fit,sid:null,poll:null,el:instEl,tabEl};
  termTabs.push(obj);
  switchTermTab(id);

  instEl.addEventListener('contextmenu',(e)=>{
    e.preventDefault();e.stopPropagation();showCtxMenu(e.clientX,e.clientY,obj);
  });

  startNewSession(obj,tabEl);

  term.onData(data=>{
    if(!obj.sid)return;
    apiFetch('/api/terminal/input/'+obj.sid,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({input:toB64(data)})});
  });

  window.addEventListener('resize',()=>{if(obj.id===activeTermId)fitTermTab(obj);});
}

function switchTermTab(id){
  activeTermId=id;
  document.querySelectorAll('.term-tab').forEach(t=>t.classList.toggle('active',t.dataset.tid===id));
  document.querySelectorAll('.term-inst').forEach(t=>t.classList.toggle('active',t.id==='ti-'+id));
  const t=termTabs.find(x=>x.id===id);
  if(t)setTimeout(()=>{t.fit.fit();t.term.focus();},60);
}

function closeTermTab(id){
  const idx=termTabs.findIndex(t=>t.id===id);if(idx===-1)return;
  const t=termTabs[idx];
  clearInterval(t.poll);
  // Hapus session ID dari cookie tapi JANGAN kill proses backend
  // Proses tetap berjalan, bisa reconnect kapan saja
  const saved=loadFromCookie('yhujin_terms',[]).filter(s=>s!==t.sid);
  saveToCookie('yhujin_terms',saved);
  t.term.dispose();t.el.remove();t.tabEl.remove();
  termTabs.splice(idx,1);
  if(!termTabs.length){newTermTab();return;}
  switchTermTab(termTabs[Math.min(idx,termTabs.length-1)].id);
}

function killTermTab(id){
  // Benar-benar matikan proses — dipanggil eksplisit oleh user
  const idx=termTabs.findIndex(t=>t.id===id);if(idx===-1)return;
  const t=termTabs[idx];
  clearInterval(t.poll);
  if(t.sid){
    apiFetch('/api/terminal/close/'+t.sid,{method:'DELETE'}).catch(()=>{});
    const saved=loadFromCookie('yhujin_terms',[]).filter(s=>s!==t.sid);
    saveToCookie('yhujin_terms',saved);
  }
  t.term.dispose();t.el.remove();t.tabEl.remove();
  termTabs.splice(idx,1);
  if(!termTabs.length){newTermTab();return;}
  switchTermTab(termTabs[Math.min(idx,termTabs.length-1)].id);
}

function killAllTerms(){
  if(!confirm('Kill semua terminal dan buat tab baru?'))return;
  // Hentikan semua session backend
  const toKill=[...termTabs];
  toKill.forEach(t=>{
    clearInterval(t.poll);
    if(t.sid)apiFetch('/api/terminal/close/'+t.sid,{method:'DELETE'}).catch(()=>{});
    try{t.term.dispose();}catch(e){}
    t.el.remove();t.tabEl.remove();
  });
  termTabs=[];
  activeTermId=null;
  // Bersihkan cookie
  saveToCookie('yhujin_terms',[]);
  // Buat satu tab baru
  newTermTab();
  toast('Semua terminal dihentikan','info');
}

function pollTerm(obj){
  if(!obj.sid)return;
  apiFetch('/api/terminal/output/'+obj.sid).then(r=>r.json()).then(d=>{
    if(d.output&&d.output.length){
      try{const bin=atob(d.output);const b=new Uint8Array(bin.length);for(let i=0;i<bin.length;i++)b[i]=bin.charCodeAt(i);obj.term.write(b);}
      catch(e){obj.term.write(d.output);}
    }
  }).catch(()=>{});
}

function fitTermTab(obj){
  obj.fit.fit();
  if(obj.sid)apiFetch('/api/terminal/resize/'+obj.sid,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({cols:obj.term.cols,rows:obj.term.rows})});
}

function applyFontSize(v){settings.termFontSize=parseInt(v);saveSetts();termTabs.forEach(t=>{t.term.options.fontSize=parseInt(v);fitTermTab(t);});}

// ── TERMINAL CONTEXT MENU ──
let ctxTermObj=null;
function showCtxMenu(x,y,obj){
  ctxTermObj=obj;
  const m=document.getElementById('termCtxMenu');
  m.style.left=x+'px';m.style.top=y+'px';
  m.classList.add('active');
  // Adjust if off screen
  setTimeout(()=>{
    const r=m.getBoundingClientRect();
    if(r.right>window.innerWidth)m.style.left=(x-r.width)+'px';
    if(r.bottom>window.innerHeight)m.style.top=(y-r.height)+'px';
  },0);
}
function hideCtxMenu(){document.getElementById('termCtxMenu').classList.remove('active');}
async function ctxCopy(){
  hideCtxMenu();
  if(!ctxTermObj)return;
  const sel=ctxTermObj.term.getSelection();
  if(sel){try{await navigator.clipboard.writeText(sel);}catch(e){toast('Copy failed','error');}}}
async function ctxPaste(){
  hideCtxMenu();
  if(!ctxTermObj)return;
  try{
    const text=await navigator.clipboard.readText();
    if(text&&ctxTermObj.sid){
      apiFetch('/api/terminal/input/'+ctxTermObj.sid,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({input:toB64(text)})});
    }
  }catch(e){toast('Paste failed — check clipboard permissions','error');}
}
function ctxCopyAll(){
  hideCtxMenu();
  if(!ctxTermObj)return;
  // Serialize entire buffer as text
  const rows=[];
  const buf=ctxTermObj.term.buffer.active;
  for(let i=0;i<buf.length;i++){
    const line=buf.getLine(i);
    if(line)rows.push(line.translateToString(true));
  }
  navigator.clipboard.writeText(rows.join('\n')).then(()=>toast('Copied all output','success')).catch(()=>toast('Copy failed','error'));
}
function ctxClear(){hideCtxMenu();if(ctxTermObj)ctxTermObj.term.clear();}
function ctxNewTab(){hideCtxMenu();newTermTab();}

// ── PROCESSES ──
function loadProcs(){
  apiFetch('/api/processes').then(r=>r.json()).then(data=>{procData=data;renderProcs();});
}
function renderProcs(){
  let d=[...procData];
  const q=document.getElementById('procSearch').value.toLowerCase();
  if(q)d=d.filter(p=>p.name.toLowerCase().includes(q)||String(p.pid).includes(q));
  d.sort((a,b)=>{
    let va=a[procSort],vb=b[procSort];
    if(typeof va==='string'){va=va.toLowerCase();vb=vb.toLowerCase();}
    return procAsc?(va>vb?1:-1):(va<vb?1:-1);
  });
  document.getElementById('procCount').textContent=d.length+' processes';
  document.getElementById('procBody').innerHTML=d.map(p=>` + "`" + `
    <tr>
      <td class="pid">${p.pid}</td>
      <td class="pname" title="${eh(p.cmdline)}">${eh(p.name)}</td>
      <td>${eh(p.username)}</td>
      <td class="${p.status==='running'?'prun':'pslp'}">${p.status}</td>
      <td class="pcpu">${p.cpu_percent}%</td>
      <td class="pmem">${p.memory_percent}%</td>
      <td>${p.memory_rss}</td>
      <td>${p.started}</td>
      <td><button class="kill-btn" onclick="killProc(${p.pid})">Kill</button></td>
    </tr>` + "`" + `).join('');
}
function filterProcs(){renderProcs();}
function sortProc(k){procSort===k?procAsc=!procAsc:(procSort=k,procAsc=false);renderProcs();}
function toggleProcRefresh(){
  const c=document.getElementById('autoRefProc').checked;
  if(c)procTimer=setInterval(loadProcs,4000);else{clearInterval(procTimer);procTimer=null;}
}
function killProc(pid){
  if(!confirm('Kill process '+pid+'?'))return;
  apiFetch('/api/process/kill/'+pid,{method:'POST'}).then(r=>r.json()).then(d=>{
    toast(d.message||d.error,d.error?'error':'success');loadProcs();
  });
}

// ── MONITOR ──
function updateMon(){
  apiFetch('/api/system-info').then(r=>r.json()).then(d=>{
    if(d.error)return;
    set('mCpuV',d.cpu_percent+'%');sw('mCpuB',d.cpu_percent);
    set('mCpuSub',d.cpu_cores+' cores / '+d.cpu_threads+' threads');
    set('mCpuF',d.cpu_freq_current?d.cpu_freq_current+' MHz':'—');
    set('mLoad',d.load_avg[0].toFixed(2));
    if(d.cpu_per_core&&d.cpu_per_core.length){
      document.getElementById('mCores').innerHTML=d.cpu_per_core.map((v,i)=>
        ` + "`" + `<div class="core-i"><div class="core-lbl">C${i}</div><div class="core-bar"><div class="core-fill" style="width:${v}%"></div></div><div style="font-size:.65rem;color:var(--tx2);text-align:right">${v}%</div></div>` + "`" + `
      ).join('');
    }
    set('mMemV',d.memory.percent+'%');sw('mMemB',d.memory.percent);
    set('mMemSub',d.memory.used+' of '+d.memory.total);
    set('mMemU',d.memory.used);set('mMemA',d.memory.available);
    set('mSwpV',d.swap.percent+'%');sw('mSwpB',d.swap.percent);
    set('mSwpSub',d.swap.used+' of '+d.swap.total);
    set('mDskV',d.disk.percent+'%');sw('mDskB',d.disk.percent);
    set('mDskSub',d.disk.used+' of '+d.disk.total);
    set('mDskU',d.disk.used);set('mDskF',d.disk.free);
    set('mHost',d.hostname);set('mPlat',d.platform+' '+d.platform_release);
    set('mArch',d.architecture);set('mUp',fmtUp(d.uptime));
    set('mUsers',d.users);set('mIp',d.ip_addresses[0]||'—');
    set('mNS',d.net_bytes_sent);set('mNR',d.net_bytes_recv);
    set('mBoot',d.boot_time?new Date(d.boot_time).toLocaleString():'—');
    set('monUpd','Updated '+new Date().toLocaleTimeString());
  });
}
function set(id,v){const e=document.getElementById(id);if(e)e.textContent=v;}
function sw(id,v){const e=document.getElementById(id);if(e)e.style.width=v+'%';}

// ── SERVERS ──
function loadServers(){
  servers=loadFromCookie('yhujin_servers',[]);
  renderServers();
}
function saveServers(){
  saveToCookie('yhujin_servers',servers);
}
function renderServers(){
  const grid=document.getElementById('srvGrid');
  const count=document.getElementById('srvCount');
  count.textContent=servers.length+' server(s)';
  if(!servers.length){
    grid.innerHTML='<div class="srv-empty" style="grid-column:1/-1"><svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" style="margin-bottom:12px;opacity:.3"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg><br>No servers added yet.<br><span style="font-size:.76rem">Click "Add Server" to add your first server.</span></div>';
    return;
  }
  grid.innerHTML=servers.map((s,i)=>{
    const statusClass=s.status||'disconnected';
    const dotClass=statusClass==='connected'?'green':(statusClass==='error'?'red':'gray');
    const statusLabel=statusClass==='connected'?'Connected':(statusClass==='error'?'Error':'Not connected');
    const tagsHtml=s.tags?s.tags.split(',').map(t=>t.trim()).filter(Boolean).map(t=>` + "`" + `<span style="background:var(--bg2);border:1px solid var(--border);border-radius:10px;padding:1px 7px;font-size:.67rem;color:var(--tx2)">${eh(t)}</span>` + "`" + `).join(''):'';
    return` + "`" + `<div class="srv-card ${statusClass}">
      <button class="srv-del-btn" onclick="deleteServer(${i})" title="Remove server">&#x00D7;</button>
      <div class="srv-name">${eh(s.name)}</div>
      <div class="srv-host">${eh(s.user?s.user+'@':'')}${eh(s.host)}:${eh(String(s.port||22))}</div>
      ${s.desc?` + "`" + `<div style="font-size:.74rem;color:var(--tx2);margin-bottom:6px">${eh(s.desc)}</div>` + "`" + `:''}
      ${tagsHtml?` + "`" + `<div class="srv-tags">${tagsHtml}</div>` + "`" + `:''}
      <div style="display:flex;align-items:center;gap:8px;margin-bottom:8px">
        <span class="srv-status"><span class="srv-dot ${dotClass}"></span>${statusLabel}</span>
        <span style="font-size:.72rem;color:var(--tx2)">${s.lastSeen?'Last seen: '+new Date(s.lastSeen).toLocaleString():''}</span>
      </div>
      <div class="srv-actions">
        <button class="btn" onclick="openSrvTerminal(${i})" title="Open SSH terminal">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>SSH
        </button>
        <button class="btn" onclick="pingSrv(${i})">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>Ping
        </button>
        <button class="btn" onclick="editServer(${i})">Edit</button>
        <button class="btn" onclick="copySshCmd(${i})">Copy SSH</button>
      </div>
    </div>` + "`" + `;
  }).join('');
}
function openAddServer(editIdx){
  document.getElementById('addSrvTitle').textContent=editIdx!=null?'Edit Server':'Add Server';
  const s=editIdx!=null?servers[editIdx]:{};
  document.getElementById('srvName').value=s.name||'';
  document.getElementById('srvHost').value=s.host||'';
  document.getElementById('srvPort').value=s.port||22;
  document.getElementById('srvUser').value=s.user||'';
  document.getElementById('srvDesc').value=s.desc||'';
  document.getElementById('srvTags').value=s.tags||'';
  document.getElementById('addServerModal')._editIdx=editIdx;
  document.getElementById('addServerModal').classList.add('active');
  setTimeout(()=>document.getElementById('srvName').focus(),100);
}
function closeAddServer(){document.getElementById('addServerModal').classList.remove('active');}
function saveServer(){
  const name=document.getElementById('srvName').value.trim();
  const host=document.getElementById('srvHost').value.trim();
  if(!name){toast('Server name is required','error');return;}
  if(!host){toast('Host is required','error');return;}
  const s={
    name,host,
    port:parseInt(document.getElementById('srvPort').value)||22,
    user:document.getElementById('srvUser').value.trim(),
    desc:document.getElementById('srvDesc').value.trim(),
    tags:document.getElementById('srvTags').value.trim(),
    status:'disconnected',
    added:new Date().toISOString()
  };
  const editIdx=document.getElementById('addServerModal')._editIdx;
  if(editIdx!=null){servers[editIdx]=Object.assign(servers[editIdx],s);}
  else{servers.push(s);}
  saveServers();renderServers();closeAddServer();
  toast(editIdx!=null?'Server updated':'Server added','success');
}
function deleteServer(i){
  if(!confirm('Remove server "'+servers[i].name+'"?'))return;
  servers.splice(i,1);saveServers();renderServers();toast('Server removed','info');
}
function editServer(i){openAddServer(i);}
function pingSrv(i){
  const s=servers[i];
  toast('Pinging '+s.host+'...','info');
  apiFetch('/api/ping',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({host:s.host})})
    .then(r=>r.json()).then(d=>{
      servers[i].status=d.reachable?'connected':'error';
      servers[i].lastSeen=d.reachable?new Date().toISOString():servers[i].lastSeen;
      saveServers();renderServers();
      toast(s.name+(d.reachable?' is reachable ('+d.latency+'ms)':' is unreachable'),d.reachable?'success':'error');
    }).catch(()=>{toast('Ping failed','error');});
}
function openSrvTerminal(i){
  const s=servers[i];
  switchPanel('terminal');
  setTimeout(()=>{
    newTermTab();
    setTimeout(()=>{
      const t=activeTerm();
      if(t&&t.sid){
        const user=s.user||'root';
        const cmd='ssh '+user+'@'+s.host+' -p '+s.port+'\n';
        apiFetch('/api/terminal/input/'+t.sid,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({input:toB64(cmd)})});
      }
    },600);
  },200);
}
function copySshCmd(i){
  const s=servers[i];
  const cmd='ssh '+(s.user?s.user+'@':'')+s.host+' -p '+s.port;
  navigator.clipboard.writeText(cmd).then(()=>toast('SSH command copied','success')).catch(()=>toast('Copy failed','error'));
}

// ── SETTINGS ──
function loadSetts(){
  const saved=loadFromCookie('yhujin_sett',{});
  Object.assign(settings,saved);
}
function saveSetts(){saveToCookie('yhujin_sett',settings);}
function applySetts(){
  document.getElementById('togHidden').classList.toggle('on',settings.showHidden);
  document.getElementById('togConfirm').classList.toggle('on',!!settings.confirmDelete);
  if(document.getElementById('settDefView'))document.getElementById('settDefView').value=settings.defaultView||'list';
  if(document.getElementById('settFS'))document.getElementById('settFS').value=settings.termFontSize||14;
}
function togSett(k,el){
  el.classList.toggle('on');settings[k]=el.classList.contains('on');saveSetts();
  if(k==='showHidden')loadFiles(curPath);
}
function saveSett(k,v){settings[k]=v;saveSetts();}
function applyRefresh(v){clearInterval(monTimer);monTimer=setInterval(updateMon,parseInt(v)*1000);}

// ── UTILS ──
function toB64(str){
  // Properly encode string (including unicode) to base64
  const bytes = new TextEncoder().encode(str);
  let binary = '';
  bytes.forEach(b => binary += String.fromCharCode(b));
  return btoa(binary);
}
function fmtUp(s){
  const d=Math.floor(s/86400),h=Math.floor((s%86400)/3600),m=Math.floor((s%3600)/60);
  if(d>0)return d+'d '+h+'h '+m+'m';if(h>0)return h+'h '+m+'m';return m+'m';
}
function fmtBytes(b){
  if(!b)return'0 B';const k=1024,sz=['B','KB','MB','GB','TB'];const i=Math.floor(Math.log(b)/Math.log(k));
  return parseFloat((b/Math.pow(k,i)).toFixed(1))+' '+sz[i];
}
function eh(s){return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');}
function eq(s){return String(s).replace(/\\/g,'\\\\').replace(/'/g,"\\'")}
function escH(s){return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');}
function toast(msg,type='info'){
  const c=document.getElementById('toastBox');
  const t=document.createElement('div');t.className='toast '+type;t.textContent=msg;
  c.appendChild(t);setTimeout(()=>t.remove(),3500);
}
</script>
</body>
</html>`
