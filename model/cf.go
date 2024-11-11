package model

type DNSRecord struct {
	Content   string `json:"content"`
	Name      string `json:"name"`
	Proxied   bool   `json:"proxied"`
	Type      string `json:"type"`
	Comment   string `json:"comment"`
	CreatedOn string `json:"created_on"`
	Id        string `json:"id"`
	Locked    bool   `json:"locked"`
	Meta      struct {
		AutoAdded bool   `json:"auto_added"`
		Source    string `json:"source"`
	}
	ModifiedOn string `json:"modified_on"`
	Proxiable  bool   `json:"proxiable"`
	Tags       []string
	TTL        int    `json:"ttl"`
	ZoneId     string `json:"zone_id"`
	ZoneName   string `json:"zone_name"`
}

type PostDNSRecord struct {
	Content string   `json:"content"`
	Name    string   `json:"name"`
	Proxied bool     `json:"proxied"`
	Type    string   `json:"type"`
	TTL     int      `json:"ttl"`
	Comment string   `json:"comment"`
	Tags    []string `json:"tags"`
}

type GetDnsRecordsResp struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []string      `json:"messages"`
	Result   []DNSRecord   `json:"result"`
}

type PostDnsRecordsReps struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []string      `json:"messages"`
	Result   DNSRecord     `json:"result"`
}

type IPChangeEvent struct {
	OldIP string
	NewIP string
	Type  string // A AAAA
}
