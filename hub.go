package main

type Message struct {
	msg []byte
	from string
}

type Hub struct {
	clients map[*Client]bool
	broadcast chan *Message
	register chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast: make(chan *Message),
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for{
		select{
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _,ok := h.clients[client]; ok{
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			var msg_sent []byte
			for client := range h.clients {
				if message.from != client.conn.RemoteAddr().String() {
					msg_sent = append([]byte("From "+message.from+"\t"), message.msg...)
				}else{
					msg_sent = append([]byte("From self\t"), message.msg...)
				}
				select {
				case client.send <- msg_sent:
				default:
					close(client.send)
					delete(h.clients, client)
				}

			}
		}
	}
}
