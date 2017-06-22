package visibility

import (
	"github.com/Frostman/aptomi/pkg/slinga"
	"sort"
)

type item struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

// TODO: change UserId -> UserId (and don't break UI...)
type detail struct {
	UserId       string
	Users        []*item
	Services     []*item
	Dependencies []*item
	Views        []*item
	Summary      slinga.ServiceUsageStateSummary
}

func NewDetails(userID string, globalUsers slinga.GlobalUsers, state slinga.ServiceUsageState) detail {
	summary := state.GetSummary()
	summary.Users = len(globalUsers.Users)
	r := detail{
		userID,
		make([]*item, 0),
		make([]*item, 0),
		make([]*item, 0),
		make([]*item, 0),
		summary,
	}

	// Users
	userIds := make([]string, 0)
	for userID := range globalUsers.Users {
		userIds = append(userIds, userID)
	}

	sort.Strings(userIds)

	if len(userIds) > 1 {
		r.Users = append([]*item{{"all", "All"}}, r.Users...)
	}
	for _, userID := range userIds {
		r.Users = append(r.Users, &item{userID, globalUsers.Users[userID].Name})
	}

	// Dependencies
	depIds := make([]string, 0)
	deps := state.Dependencies.DependenciesByID
	for depID, dep := range deps {
		if dep.UserID != userID {
			continue
		}

		depIds = append(depIds, depID)
	}

	sort.Strings(depIds)

	if len(depIds) > 1 {
		r.Dependencies = append([]*item{{"all", "All"}}, r.Dependencies...)
	}
	for _, depID := range depIds {
		r.Dependencies = append(r.Dependencies, &item{depID, deps[depID].ID})
	}

	// Services
	svcIds := make([]string, 0)
	for svcID, svc := range state.Policy.Services {
		if svc.Owner != userID {
			continue
		}
		svcIds = append(svcIds, svcID)
	}

	sort.Strings(svcIds)

	for _, svcID := range svcIds {
		r.Services = append(r.Services, &item{svcID, state.Policy.Services[svcID].Name})
	}

	if len(r.Dependencies) > 0 {
		r.Views = append(r.Views, &item{"consumer", "Service Consumer View"})
	}
	if len(r.Services) > 0 {
		r.Views = append(r.Views, &item{"service", "Service Owner View"})
	}
	if globalUsers.Users[userID].Labels["global_ops"] == "true" {
		r.Views = append(r.Views, &item{"globalops", "Global IT/Ops View"})
	}

	return r
}
