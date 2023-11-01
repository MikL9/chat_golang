package ws

type ChatCreated struct {
	ID      int     `json:"id"`
	IsGroup bool    `json:"is_group"`
	Logo    int     `json:"logo"`
	Name    string  `json:"name"`
	Users   []*User `json:"users"`
}
