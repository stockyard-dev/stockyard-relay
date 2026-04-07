package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Relay</title><link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet"><style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}.main{padding:1.5rem;max-width:960px;margin:0 auto}.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}.toolbar{display:flex;gap:.5rem;margin-bottom:1rem}.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.search:focus{outline:none;border-color:var(--leather)}.item{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}.item:hover{border-color:var(--leather)}.item-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}.item-title{font-size:.85rem;font-weight:700}.item-sub{font-size:.7rem;color:var(--cd);margin-top:.1rem}.item-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}.item-notes{font-size:.65rem;color:var(--cm);margin-top:.3rem;font-style:italic;padding:.3rem .5rem;border-left:2px solid var(--bg3);display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}.item-actions{display:flex;gap:.3rem;flex-shrink:0}.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm)}.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}.btn-sm{font-size:.55rem;padding:.2rem .4rem}.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw;max-height:90vh;overflow-y:auto}.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> RELAY</h1><button class="btn btn-p" onclick="openForm()">+ Add</button></div>
<div class="main"><div class="stats" id="stats"></div><div class="toolbar"><input class="search" id="search" placeholder="Search..." oninput="render()"></div><div id="list"></div></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/channels').then(function(r){return r.json()});items=r.channels||[];document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+items.length+'</div><div class="st-l">Total</div></div><div class="st"><div class="st-v">-</div><div class="st-l">&nbsp;</div></div><div class="st"><div class="st-v">-</div><div class="st-l">&nbsp;</div></div>';render();}
function render(){var q=(document.getElementById('search').value||'').toLowerCase();var f=items;
if(q)f=f.filter(function(i){return (i.name||'').toLowerCase().includes(q)||(i.slug||'').toLowerCase().includes(q)||(i.targets||'').toLowerCase().includes(q)||(i.delivery_count||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No items found.</div>';return;}
var h='';f.forEach(function(i){
h+='<div class="item"><div class="item-top"><div style="flex:1"><div class="item-title">'+esc(i.name)+'</div>';
if(i.slug)h+='<div class="item-sub">'+esc(i.slug)+'</div>';
h+='</div><div class="item-actions"><button class="btn btn-sm" onclick="openEdit(\''+i.id+'\')">Edit</button><button class="btn btn-sm" onclick="del(\''+i.id+'\')" style="color:var(--red)">&#10005;</button></div></div>';
h+='<div class="item-meta">';
if(i.targets)h+='<span>'+esc(i.targets)+'</span>';
if(i.delivery_count)h+='<span>'+esc(i.delivery_count)+'</span>';
if(i.last_webhook)h+='<span>'+esc(i.last_webhook)+'</span>';
if(i.channel_id)h+='<span>'+esc(i.channel_id)+'</span>';
if(i.method)h+='<span>'+esc(i.method)+'</span>';
if(i.headers)h+='<span>'+esc(i.headers)+'</span>';
h+='<span>'+ft(i.created_at)+'</span></div>';
if(i.body)h+='<div class="item-notes">'+esc((i.body||'').substring(0,200))+'</div>';
h+='</div>';});document.getElementById('list').innerHTML=h;}
async function del(id){if(!confirm('Delete?'))return;await fetch(A+'/channels/'+id,{method:'DELETE'});load();}
function formHTML(item){var i=item||{name:"",slug:"",targets:"",delivery_count:"",last_webhook:"",channel_id:"",method:"",headers:"",body:"",source_ip:"",status:"",targets_hit:"",channels:"",deliveries:""};var isEdit=!!item;
var h='<h2>'+(isEdit?'EDIT':'NEW')+' RELAY</h2>';
h+='<div class="fr"><label>Name *</label><input id="f-name" value="'+esc(i.name)+'"></div>';
h+='<div class="fr"><label>Slug</label><input id="f-slug" value="'+esc(i.slug)+'"></div>';
h+='<div class="fr"><label>Targets</label><input id="f-targets" value="'+esc(i.targets)+'"></div>';
h+='<div class="fr"><label>Delivery Count</label><input id="f-delivery-count" value="'+esc(i.delivery_count)+'"></div>';
h+='<div class="fr"><label>Last Webhook</label><input id="f-last-webhook" value="'+esc(i.last_webhook)+'"></div>';
h+='<div class="fr"><label>Channel I D</label><input id="f-channel-id" value="'+esc(i.channel_id)+'"></div>';
h+='<div class="fr"><label>Method</label><input id="f-method" value="'+esc(i.method)+'"></div>';
h+='<div class="fr"><label>Headers</label><input id="f-headers" value="'+esc(i.headers)+'"></div>';
h+='<div class="fr"><label>Body</label><textarea id="f-body" rows="3">'+esc(i.body)+'</textarea></div>';
h+='<div class="fr"><label>Source I P</label><input id="f-source-ip" value="'+esc(i.source_ip)+'"></div>';
h+='<div class="fr"><label>Status</label><input id="f-status" value="'+esc(i.status)+'"></div>';
h+='<div class="fr"><label>Targets Hit</label><input id="f-targets-hit" value="'+esc(i.targets_hit)+'"></div>';
h+='<div class="fr"><label>Channels</label><input id="f-channels" value="'+esc(i.channels)+'"></div>';
h+='<div class="fr"><label>Deliveries</label><input id="f-deliveries" value="'+esc(i.deliveries)+'"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button></div>';return h;}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var x=null;for(var j=0;j<items.length;j++){if(items[j].id===id){x=items[j];break;}}if(!x)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(x);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var body={name:document.getElementById('f-name').value.trim(),slug:document.getElementById('f-slug').value.trim(),targets:document.getElementById('f-targets').value.trim(),delivery_count:document.getElementById('f-delivery-count').value.trim(),last_webhook:document.getElementById('f-last-webhook').value.trim(),channel_id:document.getElementById('f-channel-id').value.trim(),method:document.getElementById('f-method').value.trim(),headers:document.getElementById('f-headers').value.trim(),body:document.getElementById('f-body').value.trim(),source_ip:document.getElementById('f-source-ip').value.trim(),status:document.getElementById('f-status').value.trim(),targets_hit:document.getElementById('f-targets-hit').value.trim(),channels:document.getElementById('f-channels').value.trim(),deliveries:document.getElementById('f-deliveries').value.trim()};
if(editId){await fetch(A+'/channels/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/channels',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}closeModal();load();}
function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});load();
</script><script>
(function(){
  fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
    if(!cfg||typeof cfg!=='object')return;
    if(cfg.dashboard_title){
      document.title=cfg.dashboard_title;
      var h1=document.querySelector('h1');
      if(h1){
        var inner=h1.innerHTML;
        var firstSpan=inner.match(/<span[^>]*>[^<]*<\/span>/);
        if(firstSpan){h1.innerHTML=firstSpan[0]+' '+cfg.dashboard_title}
        else{h1.textContent=cfg.dashboard_title}
      }
    }
  }).catch(function(){});
})();
</script>
</body></html>`
