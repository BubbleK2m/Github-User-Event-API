package model

type GithubActor struct {
	Name string `json:"login"`
}

type GithubEvent struct {
	Type string       `json:"type"`
	Actor GithubActor `json:"actor"`
}
