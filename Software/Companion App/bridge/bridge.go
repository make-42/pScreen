package bridge

import (
	"pscreen/bridge/comms"
	"pscreen/bridge/modules"
	"pscreen/bridge/modules/blank"
	"pscreen/bridge/modules/clock"
	"pscreen/bridge/modules/discord"
	"pscreen/bridge/modules/media"
	"pscreen/bridge/modules/monitor"
	"pscreen/bridge/modules/notifications"
	"pscreen/bridge/modules/osu"
	"pscreen/bridge/modules/pong"
	"pscreen/bridge/modules/qrcode"
	"pscreen/bridge/modules/screensaver"
	"pscreen/bridge/modules/weather"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"pscreen/constants"
	"runtime"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type bridgeData struct {
	LoadedModules      []int
	ModuleDisplayStart int64
	CurrentModule      int
	DetectedPorts      []*enumerator.PortDetails
	CommsReady         bool
}

var BridgeData = bridgeData{
	LoadedModules: []int{},
	CommsReady:    false,
}
var Port serial.Port
var FrameDeltaTime int64

func BridgeStartXMit() {
	Port = comms.EstablishComms()
	BridgeData.CommsReady = true
}

func BridgeEnumSerialDevices() {
	BridgeData.DetectedPorts = comms.EnumSerialDevices()
}

func BridgeMainThread() {
	BridgeData.LoadedModules = config.ModuleConfigNamesToIDs(config.Config.LoadedModules)
	if runtime.GOOS == "linux" {
		if (config.Config.UseNotificationsModule) && (!notifications.CurrentModuleState.ReceivingNotifications) {
			go notifications.ListenForNotifications()
		}
	}
	lastFrameT := time.Now().UnixNano()
	for {
		if time.Now().UTC().UnixMilli()-BridgeData.ModuleDisplayStart > int64(config.Config.ChangeModuleEveryXMilliseconds) {
			BridgeData.ModuleDisplayStart = time.Now().UTC().UnixMilli()
			if BridgeData.CurrentModule < len(BridgeData.LoadedModules)-1 {
				BridgeData.CurrentModule++
			} else {
				BridgeData.CurrentModule = 0
			}
		}
		frameBytes := renderer.RenderFrame(ReturnCurrentModule())
		if BridgeData.CommsReady {
			comms.SendBytes(Port, frameBytes)
		}
		time.Sleep(time.Millisecond*time.Duration(config.Config.RenderDeviceScreenEveryXMilliseconds) - time.Nanosecond*time.Duration(time.Now().UnixNano()-lastFrameT))
		FrameDeltaTime = time.Now().UnixNano() - lastFrameT
		lastFrameT = time.Now().UnixNano()
	}
}

func ReturnCurrentModule() modules.Module {
	if notifications.CurrentModuleState.DisplayingNotification {
		return notifications.NotificationsModule
	}
	if len(BridgeData.LoadedModules) > 0 {
		if BridgeData.CurrentModule > len(BridgeData.LoadedModules)-1 {
			BridgeData.CurrentModule = 0
		}
		switch BridgeData.LoadedModules[BridgeData.CurrentModule] {
		case constants.BlankModuleID:
			return blank.BlankModule
		case constants.ClockModuleID:
			return clock.ClockModule
		case constants.DiscordModuleID:
			return discord.DiscordModule
		case constants.MediaModuleID:
			return media.MediaModule
		case constants.MonitorModuleID:
			return monitor.MonitorModule
		case constants.OsuModuleID:
			return osu.OsuModule
		case constants.PongModuleID:
			return pong.PongModule
		case constants.QRCodeModuleID:
			return qrcode.QRCodeModule
		case constants.ScreensaverModuleID:
			return screensaver.ScreensaverModule
		case constants.WeatherModuleID:
			return weather.WeatherModule
		}
	}
	return clock.ClockModule
}
