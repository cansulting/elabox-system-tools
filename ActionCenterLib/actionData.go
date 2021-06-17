package actioncenter

type ActionData struct {
	Action string `json:"action"`
	State  string `json:"state"`
	Data   string `json:"data"`
}

type SubscriptionData struct {
	Action string `json:"action"`
	AppId  string `json:"appId"`
}
