package handler

import "net/http"

const homePage = `<!doctype html>
<html lang="zh-CN"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>bell-foundry 工艺台</title><link rel="stylesheet" href="/assets/style.css"></head>
<body><main><header><p class="eyebrow">BELL FOUNDRY / PROCESS DESK</p><h1>铸造工艺台</h1><p class="lede">从配料、烘型到调音，把每一炉的证据留在现场。</p></header>
<section class="cards"><article><span>数据库</span><strong id="db">读取中</strong></article><article><span>铸型</span><strong id="molds">-</strong></article><article><span>待处理告警</span><strong id="alerts">-</strong></article></section>
<section class="panel"><h2>现场入口</h2><p>使用 API 记录炉次、浇注和冷却采样；页面只展示当前工艺快照。</p><button id="refresh">刷新状态</button><pre id="output"></pre></section>
<script src="/assets/app.js"></script></main></body></html>`

const appJS = `async function refresh(){const h=await fetch('/healthz').then(r=>r.json());document.querySelector('#db').textContent=h.database||'未知';document.querySelector('#molds').textContent=h.molds??'-';document.querySelector('#output').textContent=JSON.stringify(h,null,2)}document.querySelector('#refresh').addEventListener('click',refresh);refresh();`
const appCSS = `:root{color-scheme:dark;font-family:ui-sans-serif,system-ui,sans-serif;background:#171512;color:#f1ead8}body{margin:0;background:radial-gradient(circle at 80% 0,#57442b,transparent 38%),#171512}main{max-width:900px;margin:0 auto;padding:48px 24px}.eyebrow{color:#d7a85f;letter-spacing:.18em;font-size:12px}h1{font-family:Georgia,serif;font-size:clamp(42px,8vw,80px);margin:8px 0}.lede{color:#b9b09f;font-size:18px}.cards{display:grid;grid-template-columns:repeat(3,1fr);gap:12px;margin:42px 0}.cards article,.panel{border:1px solid #66553a;background:#211d18;padding:20px}.cards span{display:block;color:#9d927f;font-size:12px}.cards strong{display:block;font-size:30px;margin-top:12px;color:#e3b96f}.panel{min-height:180px}button{border:0;background:#d59d4f;color:#171512;padding:10px 18px;font-weight:700;cursor:pointer}pre{white-space:pre-wrap;color:#c5b99e}@media(max-width:600px){.cards{grid-template-columns:1fr}.cards article{display:flex;justify-content:space-between;align-items:center}}
`

func (rt *Router) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(homePage))
		return
	}
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if r.URL.Path == "/assets/app.js" {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write([]byte(appJS))
		return
	}
	if r.URL.Path == "/assets/style.css" {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write([]byte(appCSS))
		return
	}
	http.NotFound(w, r)
}

func (rt *Router) handleHealth(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, rt.lab.Health(r.Context()))
}
