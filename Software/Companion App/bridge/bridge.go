package bridge

import (
	"pscreenapp/bridge/comms"
	"pscreenapp/bridge/modules"
	"pscreenapp/bridge/modules/blank"
	"pscreenapp/bridge/modules/clock"
	"pscreenapp/bridge/modules/media"
	"pscreenapp/bridge/modules/monitor"
	"pscreenapp/bridge/modules/notifications"
	"pscreenapp/bridge/modules/pong"
	"pscreenapp/bridge/modules/screensaver"
	"pscreenapp/bridge/modules/weather"
	"pscreenapp/bridge/renderer"
	"pscreenapp/config"
	"pscreenapp/constants"
	"pscreenapp/utils"
	"runtime"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type bridgeData struct {
	LoadedModules       []int
	DelayBetweenModules int
	ModuleDisplayStart  int64
	CurrentModule       int
	DetectedPorts       []*enumerator.PortDetails
	CommsReady          bool
}

var BridgeData = bridgeData{LoadedModules: []int{constants.MediaModuleID}, DelayBetweenModules: config.DelayBetweenModules, CommsReady: false}
var Port serial.Port

func BridgeStartXMit() {
	Port = comms.EstablishComms()
	BridgeData.CommsReady = true
}

func BridgeEnumSerialDevices() {
	var err error
	BridgeData.DetectedPorts, err = enumerator.GetDetailedPortsList()
	utils.CheckError(err)
}

func BridgeMainThread() {
	if runtime.GOOS == "linux" {
		if (config.UseNotificationsModule) && (!notifications.CurrentModuleState.ReceivingNotifications) {
			go notifications.ListenForNotifications()
		}
	}
	for {
		if time.Now().UTC().UnixMilli()-BridgeData.ModuleDisplayStart > int64(BridgeData.DelayBetweenModules)*1000 {
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
		time.Sleep(time.Second * config.RenderDeviceScreenEveryXMilliseconds / 1000)
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
		case constants.MediaModuleID:
			return media.MediaModule
		case constants.MonitorModuleID:
			return monitor.MonitorModule
		case constants.PongModuleID:
			return pong.PongModule
		case constants.ScreensaverModuleID:
			return screensaver.ScreensaverModule
		case constants.WeatherModuleID:
			return weather.WeatherModule
		}
	}
	return clock.ClockModule
}
