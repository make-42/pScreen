package bridge

import (
	"pscreen/bridge/comms"
	"pscreen/bridge/modules"
	"pscreen/bridge/modules/blank"
	"pscreen/bridge/modules/clock"
	"pscreen/bridge/modules/discord"
	"pscreen/bridge/modules/dvd"
	"pscreen/bridge/modules/media"
	"pscreen/bridge/modules/monitor"
	"pscreen/bridge/modules/notifications"
	"pscreen/bridge/modules/osu"
	"pscreen/bridge/modules/pong"
	"pscreen/bridge/modules/qrcode"
	"pscreen/bridge/modules/screensaver"
	"pscreen/bridge/modules/visualizer"
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
	TimeOfSwitch       int64
	LastMod            modules.Module
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
			BridgeData.LastMod = ReturnCurrentModule()
			lastModuleID := constants.ClockModuleID
			if len(BridgeData.LoadedModules) > 0 {
				lastModuleID = BridgeData.LoadedModules[BridgeData.CurrentModule]
			}
			BridgeData.CurrentModule = (BridgeData.CurrentModule + 1) % len(BridgeData.LoadedModules)
			if BridgeData.LoadedModules[BridgeData.CurrentModule] != lastModuleID {
				BridgeData.TimeOfSwitch = time.Now().UTC().UnixMilli()
			}
		}
		frameBytes := renderer.RenderFrame(ReturnCurrentModule(), BridgeData.LastMod, BridgeData.TimeOfSwitch)
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
		case constants.DVDModuleID:
			return dvd.DVDModule
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
		case constants.VisualizerModuleID:
			return visualizer.VisualizerModule
		case constants.WeatherModuleID:
			return weather.WeatherModule
		}
	}
	return clock.ClockModule
}
