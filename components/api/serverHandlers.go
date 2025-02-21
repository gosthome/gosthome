package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	ehp "github.com/gosthome/gosthome/components/api/esphomeproto"
	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/gosthome/gosthome/components/button"
	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/bus"
	"github.com/gosthome/gosthome/core/entity"
)

var (
	defaultHandlers = safeMessageHandlers{m: map[int]AnyMessageHandler{}}
)

func dH(mh MessageHandler) byte {
	defaultHandlers.locked(func(m map[int]AnyMessageHandler) {
		m[mh.Type] = mh.Handler
	})
	return 0
}

func WithAuth(mh MessageHandler) MessageHandler {
	oh := mh.Handler
	mh.Handler = func(ctx context.Context, c *Connection, m ehp.EsphomeMessageTyper) ([]ehp.EsphomeMessageTyper, error) {
		if !c.authenticated {
			return nil, errors.New("unauthenticated access")
		}
		return oh(ctx, c, m)
	}
	return mh
}

const (
	apiVersionMajor = 1
	apiVersionMinor = 10
)

var (
	_ = dH(Handler[ehp.HelloRequest](func(ctx context.Context, c *Connection, msg *ehp.HelloRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Info("Client connected", "clientApiVersionMajor", msg.ApiVersionMajor, "clientApiVersionMinor", msg.ApiVersionMinor, "clientInfo", msg.ClientInfo)
		c.clientInfo = msg.ClientInfo
		return []ehp.EsphomeMessageTyper{
			&ehp.HelloResponse{
				ApiVersionMajor: apiVersionMajor,
				ApiVersionMinor: apiVersionMinor,

				// Name:       c.server.config.Gosthome.Name,
				ServerInfo: fmt.Sprintf("gosthome %s based on aioesphomeapi %s", core.Version(), ehp.ESPHOME_VERSION),
			},
		}, nil
	}))
	_ = dH(Handler[ehp.ConnectRequest](func(ctx context.Context, c *Connection, msg *ehp.ConnectRequest) ([]ehp.EsphomeMessageTyper, error) {
		pass := c.server.config.Password
		valid := !pass.Valid() || pass.Check(msg.Password)
		slog.Info("Connect request", "valid", valid)
		if valid {
			c.authenticated = true
		}
		c.canAuthenticate = false
		return []ehp.EsphomeMessageTyper{
			&ehp.ConnectResponse{
				InvalidPassword: !valid,
			},
		}, nil
	}))
	_ = dH(Handler[ehp.DisconnectRequest](func(ctx context.Context, c *Connection, msg *ehp.DisconnectRequest) ([]ehp.EsphomeMessageTyper, error) {
		// option (needs_setup_connection) = false;
		// option (needs_authentication) = false;
		return []ehp.EsphomeMessageTyper{&ehp.DisconnectResponse{}}, frameshakers.ErrCloseConnection
	}))
	_ = dH(Handler[ehp.PingRequest](func(ctx context.Context, c *Connection, msg *ehp.PingRequest) ([]ehp.EsphomeMessageTyper, error) {

		// option (needs_setup_connection) = false;
		// option (needs_authentication) = false;
		return []ehp.EsphomeMessageTyper{&ehp.PingResponse{}}, nil
	}))
	_ = dH(Handler[ehp.DeviceInfoRequest](func(ctx context.Context, c *Connection, msg *ehp.DeviceInfoRequest) ([]ehp.EsphomeMessageTyper, error) {
		cfg := core.GetNode(ctx).Config
		serverCfg := c.server.config
		return []ehp.EsphomeMessageTyper{&ehp.DeviceInfoResponse{
			UsesPassword:    serverCfg.Password.Valid(),
			Name:            cfg.Gosthome.Name,
			FriendlyName:    cfg.Gosthome.FriendlyName,
			SuggestedArea:   cfg.Gosthome.Area,
			MacAddress:      cfg.Gosthome.MAC.String(),
			EsphomeVersion:  ehp.ESPHOME_VERSION,
			CompilationTime: "2022",
			Manufacturer:    "gosthome",
			Model:           runtime.GOOS + "/" + runtime.GOARCH,
			HasDeepSleep:    false,
			ProjectName:     cfg.Gosthome.Project.Name,
			ProjectVersion:  cfg.Gosthome.Project.Version,
			// WebserverPort:               0,
			LegacyBluetoothProxyVersion: 0,
			BluetoothProxyFeatureFlags:  0,
			LegacyVoiceAssistantVersion: 0,
			VoiceAssistantFeatureFlags:  0,
		}}, nil
	}))
	_ = dH(Handler[ehp.GetTimeRequest](func(ctx context.Context, c *Connection, msg *ehp.GetTimeRequest) ([]ehp.EsphomeMessageTyper, error) {
		return []ehp.EsphomeMessageTyper{&ehp.GetTimeResponse{
			EpochSeconds: uint32(time.Now().Unix()),
		}}, nil
	}))
	_ = dH(WithAuth(Handler[ehp.ListEntitiesRequest](func(ctx context.Context, c *Connection, msg *ehp.ListEntitiesRequest) ([]ehp.EsphomeMessageTyper, error) {
		ret := []ehp.EsphomeMessageTyper{}
		for _, ent := range entity.IterateRegistry(core.GetNode(ctx).Registry) {
			if ent.Internal() {
				continue
			}
			switch typed := ent.(type) {
			case entity.BinarySensor:
				ret = append(ret, &ehp.ListEntitiesBinarySensorResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon:        typed.Icon(),
					DeviceClass: typed.DeviceClass(),
				})
			case entity.Cover:
				ret = append(ret, &ehp.ListEntitiesCoverResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon:        typed.Icon(),
					DeviceClass: typed.DeviceClass(),
				})
			case entity.Fan:
				ret = append(ret, &ehp.ListEntitiesFanResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Light:
				ret = append(ret, &ehp.ListEntitiesLightResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Sensor:
				ret = append(ret, &ehp.ListEntitiesSensorResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       typed.DeviceClass(),
					UnitOfMeasurement: typed.UnitOfMeasurement(),
				})
			case entity.Switch:
				ret = append(ret, &ehp.ListEntitiesSwitchResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon:        typed.Icon(),
					DeviceClass: typed.DeviceClass(),
				})
			case entity.Button:
				ret = append(ret, &ehp.ListEntitiesButtonResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon:        typed.Icon(),
					DeviceClass: typed.DeviceClass(),
				})
			case entity.TextSensor:
				ret = append(ret, &ehp.ListEntitiesTextSensorResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon:        typed.Icon(),
					DeviceClass: typed.DeviceClass(),
				})
			case entity.Service:
				ret = append(ret, &ehp.ListEntitiesServicesResponse{
					Key:  typed.HashID(),
					Name: typed.Name(),
					Args: nil, // TODO: fill
				})
			case entity.Camera:
				ret = append(ret, &ehp.ListEntitiesCameraResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Climate:
				ret = append(ret, &ehp.ListEntitiesClimateResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Number:
				ret = append(ret, &ehp.ListEntitiesNumberResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon:              typed.Icon(),
					DeviceClass:       typed.DeviceClass(),
					UnitOfMeasurement: typed.UnitOfMeasurement(),
				})
			case entity.Date:
				ret = append(ret, &ehp.ListEntitiesDateResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Time:
				ret = append(ret, &ehp.ListEntitiesTimeResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Datetime:
				ret = append(ret, &ehp.ListEntitiesDateTimeResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Text:
				ret = append(ret, &ehp.ListEntitiesTextResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Select:
				ret = append(ret, &ehp.ListEntitiesSelectResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Lock:
				ret = append(ret, &ehp.ListEntitiesLockResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Valve:
				ret = append(ret, &ehp.ListEntitiesValveResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon:        typed.Icon(),
					DeviceClass: typed.DeviceClass(),
				})
			case entity.MediaPlayer:
				ret = append(ret, &ehp.ListEntitiesMediaPlayerResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.AlarmControlPanel:
				ret = append(ret, &ehp.ListEntitiesAlarmControlPanelResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon: typed.Icon(),
				})
			case entity.Event:
				ret = append(ret, &ehp.ListEntitiesEventResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon:        typed.Icon(),
					DeviceClass: typed.DeviceClass(),
				})
			case entity.Update:
				ret = append(ret, &ehp.ListEntitiesUpdateResponse{
					ObjectId:          typed.ID(),
					Key:               typed.HashID(),
					DisabledByDefault: typed.DisabledByDefault(),
					EntityCategory:    enum[ehp.EntityCategory](typed.EntityCategory()),
					Name:              typed.Name(),
					// UniqueId:          c.server.DefaultUniqueId(t, typed),
					Icon:        typed.Icon(),
					DeviceClass: typed.DeviceClass(),
				})
			default:
			}
		}
		return append(ret, &ehp.ListEntitiesDoneResponse{}), nil
	})))
	_ = dH(Handler[ehp.SubscribeStatesRequest](func(ctx context.Context, c *Connection, msg *ehp.SubscribeStatesRequest) ([]ehp.EsphomeMessageTyper, error) {
		c.subscribed = true
		b := bus.Get(ctx)
		sub := b.HandleEvents(bus.EventHandler(func(t *bus.StateChangeEvent) {
			r := stateResponse(t.Key, t.NewState)
			if r != nil {
				slog.Debug("Sending state change", "key", t.Key, "state", t.NewState, "to", c.clientInfo)
				err := c.SendMessages([]ehp.EsphomeMessageTyper{r})
				if err != nil {
					slog.Error("Failed to send state update", "err", err)
				}
			}
		}))
		c.busEvents = append(c.busEvents, sub)
		ret := []ehp.EsphomeMessageTyper{}
		for _, ent := range entity.IterateRegistry(core.GetNode(ctx).Registry) {
			es := entityState(ent)
			if es == nil {
				continue
			}
			ret = append(ret, es)
		}

		return ret, nil
	}))
	_ = dH(Handler[ehp.SubscribeLogsRequest](func(ctx context.Context, c *Connection, msg *ehp.SubscribeLogsRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command subscribe_logs, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.SubscribeHomeassistantServicesRequest](func(ctx context.Context, c *Connection, msg *ehp.SubscribeHomeassistantServicesRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command subscribe_homeassistant_services, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.SubscribeHomeAssistantStatesRequest](func(ctx context.Context, c *Connection, msg *ehp.SubscribeHomeAssistantStatesRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command subscribe_home_assistant_states, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.ExecuteServiceRequest](func(ctx context.Context, c *Connection, msg *ehp.ExecuteServiceRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command execute_service, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.CoverCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.CoverCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command cover_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.FanCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.FanCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command fan_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.LightCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.LightCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command light_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.SwitchCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.SwitchCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command switch_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.CameraImageRequest](func(ctx context.Context, c *Connection, msg *ehp.CameraImageRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command camera_image, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.ClimateCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.ClimateCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command climate_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.NumberCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.NumberCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command number_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.SelectCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.SelectCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command select_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.TextCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.TextCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command text_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.SirenCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.SirenCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command siren_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.ButtonCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.ButtonCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		core.GetNode(ctx).Bus.CallService(&button.ButtonPress{
			Key: msg.Key,
		})
		return nil, nil
	}))
	_ = dH(Handler[ehp.LockCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.LockCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command lock_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.ValveCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.ValveCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command valve_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.MediaPlayerCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.MediaPlayerCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command media_player_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.DateCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.DateCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command date_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.TimeCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.TimeCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command time_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.DateTimeCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.DateTimeCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command datetime_command, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.SubscribeBluetoothLEAdvertisementsRequest](func(ctx context.Context, c *Connection, msg *ehp.SubscribeBluetoothLEAdvertisementsRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command subscribe_bluetooth_le_advertisements, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.BluetoothDeviceRequest](func(ctx context.Context, c *Connection, msg *ehp.BluetoothDeviceRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_device_request, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.BluetoothGATTGetServicesRequest](func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTGetServicesRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_get_services, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.BluetoothGATTReadRequest](func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTReadRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_read, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.BluetoothGATTWriteRequest](func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTWriteRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_write, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.BluetoothGATTReadDescriptorRequest](func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTReadDescriptorRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_read_descriptor, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.BluetoothGATTWriteDescriptorRequest](func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTWriteDescriptorRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_write_descriptor, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.BluetoothGATTNotifyRequest](func(ctx context.Context, c *Connection, msg *ehp.BluetoothGATTNotifyRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command bluetooth_gatt_notify, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.UnsubscribeBluetoothLEAdvertisementsRequest](func(ctx context.Context, c *Connection, msg *ehp.UnsubscribeBluetoothLEAdvertisementsRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command unsubscribe_bluetooth_le_advertisements, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.SubscribeVoiceAssistantRequest](func(ctx context.Context, c *Connection, msg *ehp.SubscribeVoiceAssistantRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command subscribe_voice_assistant, doing nothing")
		return nil, nil
	}))
	_ = dH(Handler[ehp.AlarmControlPanelCommandRequest](func(ctx context.Context, c *Connection, msg *ehp.AlarmControlPanelCommandRequest) ([]ehp.EsphomeMessageTyper, error) {
		slog.Warn("gosthome Node got command alarm_control_panel_command, doing nothing")
		return nil, nil
	}))
)
