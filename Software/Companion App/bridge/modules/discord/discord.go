package discord

import (
	"image"
	"pscreenapp/bridge/modules"
	"pscreenapp/bridge/renderer"
	"pscreenapp/bridge/renderer/elements"
	"pscreenapp/config"
	"pscreenapp/utils"
	"time"

	"git.sr.ht/~sbinet/gg"
	"github.com/bwmarrin/discordgo"
)

var discordSession *discordgo.Session
var firstFrame = true

var DiscordModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	if firstFrame {
		firstFrame = false
		var err error
		discordSession, err = discordgo.New("Bot " + config.DiscordAuthToken)
		utils.CheckError(err)
		discordSession.StateEnabled = true
		discordSession.State.TrackChannels = true
		discordSession.State.TrackMembers = true
		discordSession.State.TrackVoice = true
		discordSession.State.TrackPresences = true
		discordSession.Identify.Intents |= discordgo.IntentGuilds
		discordSession.Identify.Intents |= discordgo.IntentGuildPresences
		discordSession.Identify.Intents |= discordgo.IntentGuildMembers
		// Open a websocket connection to Discord and begin listening.
		err = discordSession.Open()
		utils.CheckError(err)
	}
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(renderer.SmallFont)
	foundPresence := false
	for _, o := range discordSession.State.Guilds {
		if o.ID == config.DiscordGuildID {
			for _, u := range o.Presences {
				if u.User.ID == config.DiscordUserID {
					foundPresence = true
					if len(u.Activities) > 0 {
						currentActivity := u.Activities[len(u.Activities)-1]
						dc.SetFontFace(renderer.TinyFont)
						textOffset := 0.0
						if currentActivity.Timestamps.EndTimestamp != 0 && currentActivity.Timestamps.StartTimestamp != 0 {
							now := time.Now()
							elements.DrawMediaProgressBar(dc, (float64(now.UnixMilli())-float64(currentActivity.Timestamps.StartTimestamp))/1000, (float64(currentActivity.Timestamps.EndTimestamp)-float64(currentActivity.Timestamps.StartTimestamp))/1000)
						} else {
							textOffset = 14
							dc.DrawStringAnchored(currentActivity.Name, 4, 2, 0, 1)
						}
						dc.DrawStringWrapped(currentActivity.State, 4, 2+textOffset, 0, 0, float64(config.CanvasRenderDimensions.X)-8, 1, gg.AlignLeft)
						dc.DrawStringWrapped(currentActivity.Details, 4, 16+textOffset, 0, 0, float64(config.CanvasRenderDimensions.X)-8, 1, gg.AlignLeft)
					}
				}
			}
		}
	}
	if !foundPresence {
		dc.DrawStringAnchored("No activities", 4, -4, 0, 1)
	}
	return renderer.AddWallpaperToFrame(renderer.RemoveAntiAliasing(im))
}}
