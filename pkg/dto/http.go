package dto

type Lookup struct {
	Key      string `json:"key"`
	Location string `json:"location"`
}

type LookupBatch struct {
	Keys     []string `json:"keys"`
	Location string   `json:"location"`
}
