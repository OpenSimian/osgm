package mgm

import (
  "github.com/gorilla/websocket"
  "fmt"
  "net/http"
  "encoding/json"
  "../../go_modules/simian"
  "github.com/satori/go.uuid"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type ClientWebsocketHandler struct {
  ClientMgr ClientManager
}

func (ch ClientWebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Println("New connection on ws")
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    fmt.Println(err)
    return
  }
  c := ch.ClientMgr.NewClient(ws)
  c.process()
}



type ClientAuthHandler struct {
  ClientMgr ClientManager
}

func (ch ClientAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  
  type clientAuthRequest struct {
    Username string
    Password string
  }
  
  var t clientAuthRequest   
  err := decoder.Decode(&t)
  if err != nil {
    fmt.Println("Invalid auth request")
    return
  }
  
  type clientAuthResponse struct {
    Uuid uuid.UUID
    Message string
    Success bool
  }
  
  sim, _ := simian.Instance()
  guid,err := sim.Auth(t.Username,t.Password);
  if err != nil {
    response := clientAuthResponse{uuid.UUID{}, err.Error(), false}
    js, err := json.Marshal(response)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    w.Header().Set("Content-Type", "application/jons")
    w.Write(js)
    return
  }
  
  token := uuid.NewV4()
  response := clientAuthResponse{guid, token.String(), true}
  js, err := json.Marshal(response)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  
  ch.ClientMgr.addAuthenticatedUser(clientAuth{guid, token, r.RemoteAddr})
  
  w.Header().Set("Content-Type", "application/jons")
  w.Write(js)
}