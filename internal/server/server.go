package server
import ("encoding/json";"fmt";"io";"log";"net/http";"strings";"github.com/stockyard-dev/stockyard-relay/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux}
func New(db *store.DB)*Server{s:=&Server{db:db,mux:http.NewServeMux()}
s.mux.HandleFunc("GET /api/channels",s.listChannels);s.mux.HandleFunc("POST /api/channels",s.createChannel);s.mux.HandleFunc("GET /api/channels/{id}",s.getChannel);s.mux.HandleFunc("DELETE /api/channels/{id}",s.deleteChannel)
s.mux.HandleFunc("GET /api/channels/{id}/deliveries",s.listDeliveries)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){
if strings.HasPrefix(r.URL.Path,"/hook/"){s.receiveHook(w,r);return};s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)listChannels(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"channels":oe(s.db.ListChannels())})}
func(s *Server)createChannel(w http.ResponseWriter,r *http.Request){var c store.Channel;json.NewDecoder(r.Body).Decode(&c);if c.Name==""||c.Slug==""{we(w,400,"name and slug required");return};c.Enabled=true;s.db.CreateChannel(&c);wj(w,201,s.db.GetChannel(c.ID))}
func(s *Server)getChannel(w http.ResponseWriter,r *http.Request){c:=s.db.GetChannel(r.PathValue("id"));if c==nil{we(w,404,"not found");return};wj(w,200,c)}
func(s *Server)deleteChannel(w http.ResponseWriter,r *http.Request){s.db.DeleteChannel(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)listDeliveries(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"deliveries":oe(s.db.ListDeliveries(r.PathValue("id"),50))})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"status":"ok","service":"relay","channels":st.Channels})}
func(s *Server)receiveHook(w http.ResponseWriter,r *http.Request){
slug:=strings.TrimPrefix(r.URL.Path,"/hook/");ch:=s.db.GetBySlug(slug)
if ch==nil{we(w,404,"channel not found");return}
if !ch.Enabled{we(w,503,"channel disabled");return}
body,_:=io.ReadAll(r.Body);hdrs:=fmt.Sprintf("%v",r.Header)
del:=&store.Delivery{ChannelID:ch.ID,Method:r.Method,Headers:hdrs,Body:string(body),SourceIP:r.RemoteAddr,Status:"received",TargetsHit:0}
s.db.RecordDelivery(del)
wj(w,200,map[string]any{"received":true,"delivery_id":del.ID})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
