package bluesky

type AuthResponse struct {
	AccessJwt string `json:"accessJwt"`
	Did       string `json:"did"`
}

type PostRecord struct {
	Type      string `json:"$type"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type CreateRecordRequest struct {
	Repo       string     `json:"repo"`
	Collection string     `json:"collection"`
	Record     PostRecord `json:"record"`
}

type FeedResponse struct {
	Feed []struct {
		Record struct {
			Text      string `json:"text"`
			CreatedAt string `json:"createdAt"`
		} `json:"record"`
	} `json:"feed"`
}
