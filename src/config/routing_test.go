package config

import (
	"reflect"
	"testing"
)

var (
	exampleRoute = Route{
		CanID:     "100",
		Topic:     "example/topic",
		Direction: BIDIRECTIONAL,
		Converter: nil,
	}
	exampleRoute2 = Route{
		CanID:     "101",
		Topic:     "example/topic2",
		Direction: CAN2MQTT,
		Converter: nil,
	}
	exampleRoute3 = Route{
		CanID:     "102",
		Topic:     "example/topic3",
		Direction: MQTT2CAN,
		Converter: nil,
	}
)

func TestRouting_AddRoute(t *testing.T) {
	type args struct {
		route Route
	}
	tests := []struct {
		name    string
		routing *Routing
		args    args
		wantErr bool
	}{
		{
			name:    "Add route",
			routing: NewRouting(),
			args:    args{route: exampleRoute},
			wantErr: false,
		},
		{
			name: "Add route already exists",
			routing: &Routing{
				routeMap: map[string]Route{exampleRoute.GetHash(): exampleRoute},
				canMap:   map[string]string{exampleRoute.CanID: exampleRoute.GetHash()},
				mqttMap:  map[string]string{exampleRoute.Topic: exampleRoute.GetHash()},
			},
			args:    args{route: exampleRoute},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.routing
			if err := r.AddRoute(tt.args.route); (err != nil) != tt.wantErr {
				t.Errorf("Routing.AddRoute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, ok := r.routeMap[tt.args.route.GetHash()]; !ok {
				t.Error("Routing.routeMap hash not exists")
			}
		})
	}
}

func TestRouting_RemoveRoute(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		routing *Routing
		args    args
		wantErr bool
	}{
		{
			name:    "Remove not existing route",
			routing: NewRouting(),
			args:    args{hash: exampleRoute.GetHash()},
			wantErr: true,
		},
		{
			name: "Remove existing route",
			routing: &Routing{
				routeMap: map[string]Route{exampleRoute.GetHash(): exampleRoute},
				canMap:   map[string]string{exampleRoute.CanID: exampleRoute.GetHash()},
				mqttMap:  map[string]string{exampleRoute.Topic: exampleRoute.GetHash()},
			},
			args:    args{hash: exampleRoute.GetHash()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.routing
			if err := r.RemoveRoute(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Routing.RemoveRoute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRouting_RemoveRouteByCanId(t *testing.T) {
	type args struct {
		canId string
	}
	tests := []struct {
		name    string
		routing *Routing
		args    args
		wantErr bool
	}{
		{
			name:    "Hash not in canMap",
			routing: NewRouting(),
			args:    args{canId: exampleRoute.CanID},
			wantErr: true,
		},
		{
			name: "Remove Route",
			routing: &Routing{
				routeMap: map[string]Route{exampleRoute.GetHash(): exampleRoute},
				canMap:   map[string]string{exampleRoute.CanID: exampleRoute.GetHash()},
				mqttMap:  map[string]string{exampleRoute.Topic: exampleRoute.GetHash()},
			},
			args:    args{canId: exampleRoute.CanID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.routing
			if err := r.RemoveRouteByCanId(tt.args.canId); (err != nil) != tt.wantErr {
				t.Errorf("Routing.RemoveRouteByCanId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRouting_RemoveRouteByMqttTopic(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		routing *Routing
		args    args
		wantErr bool
	}{
		{
			name:    "Hash not in canMap",
			routing: NewRouting(),
			args:    args{topic: exampleRoute.Topic},
			wantErr: true,
		},
		{
			name: "Remove Route",
			routing: &Routing{
				routeMap: map[string]Route{exampleRoute.GetHash(): exampleRoute},
				canMap:   map[string]string{exampleRoute.CanID: exampleRoute.GetHash()},
				mqttMap:  map[string]string{exampleRoute.Topic: exampleRoute.GetHash()},
			},
			args:    args{topic: exampleRoute.Topic},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.routing
			if err := r.RemoveRouteByMqttTopic(tt.args.topic); (err != nil) != tt.wantErr {
				t.Errorf("Routing.RemoveRouteByMqttTopic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRouting_AddRoutes(t *testing.T) {
	type args struct {
		routes []Route
	}
	tests := []struct {
		name    string
		routing *Routing
		args    args
		wantErr bool
	}{
		{
			name:    "Add new routes",
			routing: NewRouting(),
			args: args{
				routes: []Route{
					exampleRoute, exampleRoute2,
				},
			},
			wantErr: false,
		},
		{
			name: "Add existing routes",
			routing: &Routing{
				routeMap: map[string]Route{exampleRoute.GetHash(): exampleRoute},
				canMap:   map[string]string{exampleRoute.CanID: exampleRoute.GetHash()},
				mqttMap:  map[string]string{exampleRoute.Topic: exampleRoute.GetHash()},
			},
			args: args{
				routes: []Route{
					exampleRoute, exampleRoute2,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.routing
			if err := r.AddRoutes(tt.args.routes); (err != nil) != tt.wantErr {
				t.Errorf("Routing.AddRoutes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRouting_UpdateRoutes(t *testing.T) {
	type args struct {
		routes []Route
	}
	tests := []struct {
		name    string
		routing *Routing
		args    args
		wantErr bool
	}{
		{
			name:    "Add new routes",
			routing: NewRouting(),
			args: args{
				routes: []Route{
					exampleRoute, exampleRoute2,
				},
			},
			wantErr: false,
		},
		{
			name: "Add existing routes",
			routing: &Routing{
				routeMap: map[string]Route{exampleRoute.GetHash(): exampleRoute},
				canMap:   map[string]string{exampleRoute.CanID: exampleRoute.GetHash()},
				mqttMap:  map[string]string{exampleRoute.Topic: exampleRoute.GetHash()},
			},
			args: args{
				routes: []Route{
					exampleRoute, exampleRoute2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.routing
			if err := r.UpdateRoutes(tt.args.routes); (err != nil) != tt.wantErr {
				t.Errorf("Routing.UpdateRoutes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRouting_CompareRoutes(t *testing.T) {
	type args struct {
		routes []Route
	}
	tests := []struct {
		name    string
		routing *Routing
		args    args
		want    []Route
		want1   []Route
		wantErr bool
	}{
		{
			name: "Compare routes",
			routing: &Routing{
				routeMap: map[string]Route{exampleRoute.GetHash(): exampleRoute, exampleRoute2.GetHash(): exampleRoute2},
				canMap:   map[string]string{exampleRoute.CanID: exampleRoute.GetHash(), exampleRoute2.CanID: exampleRoute2.GetHash()},
				mqttMap:  map[string]string{exampleRoute.Topic: exampleRoute.GetHash(), exampleRoute2.Topic: exampleRoute2.GetHash()},
			},
			args:    args{routes: []Route{exampleRoute, exampleRoute3}},
			want:    []Route{exampleRoute3},
			want1:   []Route{exampleRoute2},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.routing
			got, got1, err := r.CompareRoutes(tt.args.routes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Routing.CompareRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Routing.CompareRoutes() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Routing.CompareRoutes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
