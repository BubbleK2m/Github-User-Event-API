package model

type GithubEvent struct {
	Type  string      `json:"type"`
	Actor GithubActor `json:"actor"`
}

type GithubActor struct {
	Name string `json:"login"`
}
