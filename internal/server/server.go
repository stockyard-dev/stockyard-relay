package server
import ("encoding/json";"net/http";"github.com/stockyard-dev/stockyard-relay/internal/store")
type Server struct{db *store.DB;limits Limits;mux *http.ServeMux}
func New(db *store.DB,tier string)*Server{s:=&Server{db:db,limits:LimitsFor(tier),mux:http.NewServeMux()};s.routes();return s}
func(s *Server)ListenAndServe(addr string)error{return(&http.Server{Addr:addr,Handler:s.mux}).ListenAndServe()}
func(s *Server)routes(){
    s.mux.HandleFunc("GET /health",s.handleHealth)
    s.mux.HandleFunc("GET /api/version",s.handleVersion)
    s.mux.HandleFunc("GET /api/limits",s.handleLimits)
    s.mux.HandleFunc("GET /api/stats",s.handleStats)
    s.mux.HandleFunc("GET /api/events",s.handleListEvents)
    s.mux.HandleFunc("POST /api/events",s.handleCreateEvent)
    s.mux.HandleFunc("DELETE /api/events/{id}",s.handleDeleteEvent)
    s.mux.HandleFunc("GET /api/subscriptions",s.handleListSubscriptions)
    s.mux.HandleFunc("POST /api/subscriptions",s.handleCreateSubscription)
    s.mux.HandleFunc("DELETE /api/subscriptions/{id}",s.handleDeleteSubscription)
    s.mux.HandleFunc("POST /api/fire/{name}",s.handleFire)
    s.mux.HandleFunc("GET /api/deliveries",s.handleListDeliveries)
    s.mux.HandleFunc("GET /",s.handleUI)
}
func(s *Server)handleHealth(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]string{"status":"ok","service":"stockyard-relay"})}  
func(s *Server)handleVersion(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]string{"version":"0.1.0","service":"stockyard-relay"})}  
func(s *Server)handleLimits(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]interface{}{"tier":s.limits.Tier,"description":s.limits.Description,"is_pro":s.limits.IsPro()})}
func writeJSON(w http.ResponseWriter,status int,v interface{}){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);json.NewEncoder(w).Encode(v)}
func writeError(w http.ResponseWriter,status int,msg string){writeJSON(w,status,map[string]string{"error":msg})}
func(s *Server)handleUI(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};w.Header().Set("Content-Type","text/html");w.Write(dashboardHTML)}
