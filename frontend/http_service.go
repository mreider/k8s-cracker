package main

import (
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
)

type HTTPService struct {
	port          int
	wsConnections *WebsocketConnections
}

func NewHTTPService(port int, wsConnections *WebsocketConnections) *HTTPService {
	return &HTTPService{
		port:          port,
		wsConnections: wsConnections,
	}
}

func (s *HTTPService) Run() {
	http.HandleFunc("/ws", s.serveWebsocket)
	http.HandleFunc("/", s.serveIndex)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *HTTPService) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer s.wsConnections.closeConnection(conn)
	s.wsConnections.addConnection(conn)

	_, _, _ = conn.Reader(r.Context())
}

func (s *HTTPService) serveIndex(w http.ResponseWriter, _ *http.Request) {
	const body = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Cracker Scores</title>
</head>
<body>
<script>

var ws;

function reconnectWebSocket() {
  if (ws != null && ws.readyState !== WebSocket.CLOSED)
    return;

  const wsProtocol = window.location.protocol === "https:" ? "wss" : "ws";
  ws = new WebSocket(wsProtocol + "://" + window.location.host + "/ws");
  ws.onmessage = function(event) {
    var message = JSON.parse(event.data);
	var cracker_element = document.getElementById(message.cracker_id);
    if (cracker_element == null) {
      cracker_element = document.createElement("div");
      cracker_element.setAttribute("id", message.cracker_id);
      document.body.append(cracker_element);
    }
    cracker_element.innerText = message.cracker_id + ": " + message.score.toString();
  };
}

reconnectWebSocket();
setInterval(reconnectWebSocket, 1000);

</script>
</body>
</html>`

	_, _ = w.Write([]byte(body))
}
