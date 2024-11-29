package config

import (
	"encoding/json"
	"log/slog"
)

type Updater struct {
	routing   *Routing
	callbacks []func([]Route, []Route)
}

func NewUpdate() *Updater {
	return &Updater{
		routing:   nil,
		callbacks: []func([]Route, []Route){},
	}
}

func (u *Updater) WithRouting(routing *Routing) *Updater {
	u.routing = routing
	return u
}

func (u *Updater) RegisterCallback(callback func([]Route, []Route)) {
	u.callbacks = append(u.callbacks, callback)
}

func (u *Updater) ConfigUpdate(config []byte) {
	var routings []Route
	if err := json.Unmarshal(config, &routings); err != nil {
		slog.Error("Unmarshal config error", "config", string(config), "error", err)
		return
	}
	defer u.routing.UpdateRoutes(routings)
	addRoute, delRoute, err := u.routing.CompareRoutes(routings)
	if err != nil {
		slog.Error("ComperRoutes error", "error", err)
		return
	}
	u.Inform(addRoute, delRoute)
}

func (u *Updater) Inform(addRoute []Route, delRoute []Route) {
	for _, callback := range u.callbacks {
		callback(addRoute, delRoute)
	}
}
