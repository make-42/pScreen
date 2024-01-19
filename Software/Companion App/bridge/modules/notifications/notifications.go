package notifications

import (
	"image"
	"pscreen/bridge/modules"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"pscreen/utils"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/make-42/go-notifications"

	"github.com/ztrue/tracerr"
)

type ModuleState struct {
	ReceivingNotifications   bool
	DisplayingNotification   bool
	Inverted                 bool
	InvertLifetime           int
	NotificationReceivedTime int64
	CurrentNotification      notifications.Notification
}

var CurrentModuleState = ModuleState{
	ReceivingNotifications:   false,
	DisplayingNotification:   false,
	Inverted:                 false,
	InvertLifetime:           config.Config.NotificationsInvertDisplayEveryXFrames,
	NotificationReceivedTime: 0,
	CurrentNotification:      notifications.Notification{},
}

func ListenForNotifications() {
	CurrentModuleState.ReceivingNotifications = true
	notificationReceiver, err := notifications.NewNotificationReceiver(config.Config.NotificationsSystemSendsDoubleNotificationMessages)
	utils.CheckError(tracerr.Wrap(err))
	channel := notificationReceiver.GetChannel()
	for v := range channel {
		utils.CheckError(tracerr.Wrap(v.Error))
		CurrentModuleState.CurrentNotification = v
		CurrentModuleState.DisplayingNotification = true
		CurrentModuleState.Inverted = true
		CurrentModuleState.InvertLifetime = config.Config.NotificationsInvertDisplayEveryXFrames
		CurrentModuleState.NotificationReceivedTime = time.Now().UTC().UnixMilli()
	}
	notificationReceiver.Close()
}

func updateModuleState(moduleState ModuleState) ModuleState {
	if (time.Now().UTC().UnixMilli() - moduleState.NotificationReceivedTime) > config.Config.NotificationsDisplayForXMilliseconds {
		moduleState.DisplayingNotification = false
	}
	if (time.Now().UTC().UnixMilli() - moduleState.NotificationReceivedTime) < config.Config.NotificationsInvertForXMilliseconds {
		if moduleState.InvertLifetime <= 0 {
			moduleState.Inverted = !moduleState.Inverted
			moduleState.InvertLifetime = config.Config.NotificationsInvertDisplayEveryXFrames
		}
	} else {
		moduleState.Inverted = false
	}
	moduleState.InvertLifetime--
	return moduleState
}

var NotificationsModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(renderer.TinyFont)
	dc.DrawStringAnchored(CurrentModuleState.CurrentNotification.Body.Summary, 4, 0, 0, 1)
	dc.SetFontFace(renderer.TinyFont)
	dc.DrawStringWrapped(CurrentModuleState.CurrentNotification.Body.Body, 4, 12, 0, 0, float64(config.Config.CanvasRenderDimensions.X)-8, 1, gg.AlignLeft)
	dc.SetFontFace(renderer.TinyFont)
	dc.DrawStringAnchored(CurrentModuleState.CurrentNotification.Body.ApplicationName, float64(config.Config.CanvasRenderDimensions.X)-4, float64(config.Config.CanvasRenderDimensions.Y)-4, 1, 0)
	CurrentModuleState = updateModuleState(CurrentModuleState)
	if CurrentModuleState.Inverted {
		return renderer.InvertImage(renderer.RemoveAntiAliasing(im))
	}
	return renderer.RemoveAntiAliasing(im)
}}
