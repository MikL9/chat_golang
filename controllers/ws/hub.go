package ws

type Hub struct {
	clients        Clients
	registerChan   chan *Client
	unregisterChan chan *Client
	broadcastChan  chan Message
}

func NewHub() *Hub {
	return &Hub{
		clients:        make(Clients),
		broadcastChan:  make(chan Message),
		registerChan:   make(chan *Client),
		unregisterChan: make(chan *Client),
	}
}

func (h *Hub) Register(cl *Client) {
	h.registerChan <- cl
}

func (h *Hub) Unregister(cl *Client) {
	h.unregisterChan <- cl
}

func (h *Hub) Broadcast(msg Message) {
	h.broadcastChan <- msg
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.registerChan:
			h.clients[cl] = struct{}{}

		case cl := <-h.unregisterChan:
			if _, ok := h.clients[cl]; !ok {
				continue
			}

			delete(h.clients, cl)
			close(cl.send)

		case msg := <-h.broadcastChan:
			strategy := msg.Strategy()
			filteredClients := strategy.FilterClients(h.clients)

			for cl := range filteredClients {
				select {
				case cl.send <- msg:
				default:
					delete(h.clients, cl)
					close(cl.send)
				}
			}
		}
	}
}
