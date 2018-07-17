package controller

import (
	"github.com/labstack/echo"
	"net/http"
	"server/model"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"sort"
)

type GithubController struct{}

func (controller GithubController) Init(e *echo.Echo) {
	e.GET("/api/github", controller.GetUserInfoFromRepo)
}

func (controller GithubController) GetUserInfoFromRepo(ctx echo.Context) error {
	query := new(struct{
		Owner string  `query:"owner"`
		Repo  string  `query:"repo"`
		Size  int     `query:"size"`
		Sort  string  `query:"sort"`
	})

	query.Owner = "rails"
	query.Repo = "rails"
	query.Size = 100
	query.Sort = "PushEvent"

	if err := ctx.Bind(query); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	events := controller.GetEventsFromRepo(query.Owner, query.Repo, query.Size)
	info := controller.GetUserInfoFromEvents(events)

	users := controller.GetUsersFromInfo(info)
	controller.SortUsersByEvent(users, query.Sort)

	return ctx.JSON(http.StatusOK, users)
}

func (GithubController) GetEventsFromRepo(owner, repo string, size int) []model.GithubEvent {
	url := "https://api.github.com/repos/%s/%s/events?per_page=%d"
	resp, _ := http.Get(fmt.Sprintf(url, owner, repo, size))
	body, _ := ioutil.ReadAll(resp.Body)

	events := make([]model.GithubEvent, 0)
	json.Unmarshal(body, &events)

	return events
}

func (GithubController) GetUserInfoFromEvents(events []model.GithubEvent) map[string]map[string]int {
	info := make(map[string]map[string]int)

	for _, event := range events {
		author := event.Actor.Name

		if _, exist := info[author]; !exist {
			info[author] = make(map[string]int)
		}

		info[author][event.Type] += 1
	}

	return info
}

func (GithubController) GetUsersFromInfo(info map[string]map[string]int) []map[string]interface{} {
	users := make([]map[string]interface{}, 0)

	for user, events := range info {
		users = append(users, map[string]interface{}{
			"login": user,
			"events": events,
		})
	}

	return users
}

func (GithubController) SortUsersByEvent(users []map[string]interface{}, key string) {
	sort.Slice(users, func(l, r int) bool {
		left := users[l]["events"].(map[string]int)
		right := users[r]["events"].(map[string]int)

		return left[key] > right[key]
	})
}