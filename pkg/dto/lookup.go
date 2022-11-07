package dto

type Lookup struct {
	Key      string `json:"key"`
	Location string `json:"location"`
}

type LookupKeyBatch struct {
	Keys     []string `json:"keys"`
	Location string   `json:"location"`
}

type LookupDomain struct {
	Domain   string `json:"domain"`
	Location string `json:"location"`
}
