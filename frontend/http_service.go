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
	<title>Cracker Score Page</title>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">

	<!-- Optional theme -->
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">
</head>
<body>

<div class="container">
<div class="starter-template" id="clist">
<h1>Here's a bunch of workers guessing codes</h1>
<script>

var ws;

function reconnectWebSocket() {
  if (ws != null && ws.readyState !== WebSocket.CLOSED)
    return;

  const wsProtocol = window.location.protocol === "https:" ? "wss" : "ws";
  ws = new WebSocket(wsProtocol + "://" + window.location.host + "/ws");
  ws.onmessage = function(event) {
    var message = JSON.parse(event.data);
	var objTo = document.getElementById('clist');
	var cracker_element = document.getElementById(message.cracker_id);
    if (cracker_element == null) {
	  cracker_element = document.createElement("div");
      cracker_element.setAttribute("id", message.cracker_id);
	  objTo.appendChild(cracker_element);
    }
    cracker_element.innerText = "cracker id: " + message.cracker_id + " correct guesses:" + message.score.toString();
  };
}

reconnectWebSocket();
setInterval(reconnectWebSocket, 1000);

</script>
</div>	
</div>
</body>
</html>`

	_, _ = w.Write([]byte(body))
}
