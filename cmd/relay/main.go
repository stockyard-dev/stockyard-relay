package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-relay/internal/server";"github.com/stockyard-dev/stockyard-relay/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="8620"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./relay-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("relay: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Relay — Self-hosted webhook relay\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Receive:    http://localhost:%s/hook/{channel}\n  Data:       %s\n  ─────────────────────────────────\n\n",port,port,port,dataDir)
log.Printf("relay: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
