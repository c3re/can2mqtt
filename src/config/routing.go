package config

import (
	"errors"
	"fmt"

	"github.com/jaster-prj/can2mqtt/common"
)

type Routing struct {
	routeMap map[string]Route
	canMap   map[string]string
	mqttMap  map[string]string
}

func NewRouting() *Routing {
	return &Routing{
		routeMap: map[string]Route{},
		canMap:   map[string]string{},
		mqttMap:  map[string]string{},
	}
}

func (r *Routing) AddRoute(route Route) error {
	if _, ok := r.routeMap[route.GetHash()]; ok {
		return fmt.Errorf("route with Hash %s already exists", route.GetHash())
	}
	r.routeMap[route.GetHash()] = route
	r.canMap[route.CanID] = route.GetHash()
	r.mqttMap[route.Topic] = route.GetHash()
	return nil
}

func (r *Routing) RemoveRoute(hash string) error {
	if _, ok := r.routeMap[hash]; !ok {
		return errors.New("route not in configuration")
	}
	delete(r.canMap, r.routeMap[hash].CanID)
	delete(r.mqttMap, r.routeMap[hash].CanID)
	delete(r.routeMap, hash)
	return nil
}

func (r *Routing) RemoveRouteByCanId(canId string) error {
	hash, ok := r.canMap[canId]
	if !ok {
		return errors.New("route not in configuration")
	}
	return r.RemoveRoute(hash)
}

func (r *Routing) RemoveRouteByMqttTopic(topic string) error {
	hash, ok := r.mqttMap[topic]
	if !ok {
		return errors.New("route not in configuration")
	}
	return r.RemoveRoute(hash)
}

func (r *Routing) GetRouteByCanId(canId string) (*Route, error) {
	hash, ok := r.canMap[canId]
	if !ok {
		return nil, errors.New("route not in configuration")
	}
	return common.POINTER(r.routeMap[hash]), nil
}

func (r *Routing) GetRouteByMqttTopic(topic string) (*Route, error) {
	hash, ok := r.mqttMap[topic]
	if !ok {
		return nil, errors.New("route not in configuration")
	}
	return common.POINTER(r.routeMap[hash]), nil
}

func (r *Routing) AddRoutes(routes []Route) error {
	var ret error
	for _, route := range routes {
		if err := r.AddRoute(route); err != nil {
			ret = errors.Join(ret, fmt.Errorf("failed adding Route for %s: %v", route.GetHash(), err))
		}
	}
	return ret
}

func (r *Routing) UpdateRoutes(routes []Route) error {
	r.routeMap = map[string]Route{}
	r.canMap = map[string]string{}
	r.mqttMap = map[string]string{}
	return r.AddRoutes(routes)
}

func (r *Routing) CompareRoutes(routes []Route) ([]Route, []Route, error) {
	addRoutes := []Route{}
	delRoutesMap := r.routeMap
	delRoutes := []Route{}
	//unchanged := []Route{}
	for _, route := range routes {
		if _, ok := delRoutesMap[route.GetHash()]; ok {
			delete(delRoutesMap, route.GetHash())
			//unchanged = append(unchanged, route)
		} else {
			addRoutes = append(addRoutes, route)
		}
	}
	for _, route := range delRoutesMap {
		delRoutes = append(delRoutes, route)
	}
	return addRoutes, delRoutes, nil
}
