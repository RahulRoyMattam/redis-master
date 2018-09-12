package slave

//SetCommand entry accepted in the POST body by agera
type SetCommand struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Expire int    `json:"expire"`
}
