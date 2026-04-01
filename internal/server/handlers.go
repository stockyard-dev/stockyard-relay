package server
import("bytes";"crypto/hmac";"crypto/rand";"crypto/sha256";"encoding/hex";"encoding/json";"fmt";"io";"net/http";"strconv";"time";"github.com/stockyard-dev/stockyard-relay/internal/store")
func(s *Server)handleListEvents(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListEvents();if list==nil{list=[]store.Event{}};writeJSON(w,200,list)}
func(s *Server)handleCreateEvent(w http.ResponseWriter,r *http.Request){
    if !s.limits.IsPro(){n,_:=s.db.CountEvents();if n>=5{writeError(w,403,"free tier: 5 events max");return}}
    var e store.Event;json.NewDecoder(r.Body).Decode(&e)
    if e.Name==""{writeError(w,400,"name required");return}
    if err:=s.db.CreateEvent(&e);err!=nil{writeError(w,500,err.Error());return}
    writeJSON(w,201,e)}
func(s *Server)handleDeleteEvent(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteEvent(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListSubscriptions(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListSubscriptions();if list==nil{list=[]store.Subscription{}};writeJSON(w,200,list)}
func(s *Server)handleCreateSubscription(w http.ResponseWriter,r *http.Request){
    var sub store.Subscription;json.NewDecoder(r.Body).Decode(&sub)
    if sub.CallbackURL==""{writeError(w,400,"callback_url required");return}
    if sub.Secret==""{b:=make([]byte,16);rand.Read(b);sub.Secret=hex.EncodeToString(b)}
    if err:=s.db.CreateSubscription(&sub);err!=nil{writeError(w,500,err.Error());return}
    writeJSON(w,201,sub)}
func(s *Server)handleDeleteSubscription(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteSubscription(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleFire(w http.ResponseWriter,r *http.Request){
    name:=r.PathValue("name")
    event,_:=s.db.GetEventByName(name);if event==nil{writeError(w,404,"event not found");return}
    body,_:=io.ReadAll(io.LimitReader(r.Body,1<<20));if !json.Valid(body){body=[]byte("{}")}
    subs,_:=s.db.GetSubsForEvent(event.ID)
    for _,sub:=range subs{d:=&store.Delivery{SubscriptionID:sub.ID,Payload:string(body)};s.db.CreateDelivery(d);go deliver(s.db,d,sub,body)}
    writeJSON(w,200,map[string]interface{}{"event":name,"deliveries_queued":len(subs)})}
func deliver(db *store.DB,d *store.Delivery,sub store.Subscription,payload[]byte){
    client:=&http.Client{Timeout:10*time.Second}
    for attempt:=1;attempt<=3;attempt++{
        d.Attempts=attempt
        req,err:=http.NewRequest("POST",sub.CallbackURL,bytes.NewReader(payload));if err!=nil{d.Status="failed";d.ErrorMsg=err.Error();db.UpdateDelivery(d);return}
        req.Header.Set("Content-Type","application/json")
        if sub.Secret!=""{mac:=hmac.New(sha256.New,[]byte(sub.Secret));mac.Write(payload);req.Header.Set("X-Relay-Signature","sha256="+hex.EncodeToString(mac.Sum(nil)))}
        resp,err:=client.Do(req);if err!=nil{d.ErrorMsg=err.Error();if attempt<3{time.Sleep(time.Duration(attempt*2)*time.Second);continue};d.Status="failed";db.UpdateDelivery(d);return}
        resp.Body.Close();d.ResponseCode=resp.StatusCode
        if resp.StatusCode>=200&&resp.StatusCode<300{now:=time.Now();d.DeliveredAt=&now;d.Status="delivered";db.UpdateDelivery(d);return}
        d.ErrorMsg=fmt.Sprintf("HTTP %d",resp.StatusCode);if attempt<3{time.Sleep(time.Duration(attempt*2)*time.Second);continue}
        d.Status="failed";db.UpdateDelivery(d);return}}
func(s *Server)handleListDeliveries(w http.ResponseWriter,r *http.Request){limit:=50;if l:=r.URL.Query().Get("limit");l!=""{if n,err:=strconv.Atoi(l);err==nil{limit=n}};list,_:=s.db.ListDeliveries(limit);if list==nil{list=[]store.Delivery{}};writeJSON(w,200,list)}
func(s *Server)handleStats(w http.ResponseWriter,r *http.Request){ev,_:=s.db.CountEvents();del,_:=s.db.CountDeliveries();writeJSON(w,200,map[string]interface{}{"events":ev,"deliveries":del})}
