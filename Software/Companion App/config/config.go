package config

import "pscreenapp/utils"

const AppVersion = "α 0.0.1"

const UpdateUIEveryXMilliseconds = 500
const PaddingIndentAmount = 1

const DelayBetweenModules = 5
const RenderDeviceScreenEveryXMilliseconds = 0

var CanvasRenderDimensions = utils.Coords{X: 256, Y: 64}

var SerialPortToUse string = ""

const SerialPortBaudRate = 115200
const RGBXMit = false

const I18nLanguage = "en_US"

const RotateScreen180Deg = false
const CircularScreenLayout = false
const UseWallpaper = true

const DebugSaveScreen = false

// Weather module
const OpenWeatherMapAPIKey = ""
const Lat = 0.0
const Long = 0.0
const UpdateWeatherEveryXMilliseconds = 240000
const WindIndicatorRadius = 6

// Monitor module
const CPUUsageSamplingMilliseconds = 1000
const CPUUsageBarMargin = 1

var CPUUsageBarDimensions = utils.Coords{X: 4, Y: 32}

// Media module
const MediaProgressBarHeight = 8
const MediaProgressBarWaveScale = 0.1
const MediaProgressBarIndicatorRadius = 4

// Pong module
const PongBallVelocity = 2.0
const PongBallRadius = 2.0
const PongPaddleWidth = 2.0
const PongPaddleInvisibleWidth = 10.0
const PongPaddleDistFromEdge = 10.0
const PongPaddleLength = 20.0
const PongPaddleVelocity = 0.2
const PongPaddleP = 0.4
const PongPaddleI = 0.1
const PongPaddleD = 0.2
const PongTimestepsPerFrame = 4

// QRCode module
const QRCodeContent = "WIFI:S:SSID;T:WPA;P:PASSWORD;H:false;;"
const QRCodeTitle = "SSID"
const QRCodeDescription = "PASSWORD"

// Screensaver module
const ScreensaverMaxParticles = 10
const ScreensaverParticlesToSpawnEachTime = 5
const ScreensaverParticleLifetime = 10
const ScreensaverParticleGravity = 0.1

var ScreensaverParticleSpeed = utils.CoordsFloat{X: 2, Y: 2}

// Notifications module
const UseNotificationsModule = true
const NotificationsDisplayForXMilliseconds = 7000
const NotificationsInvertForXMilliseconds = 1000
const NotificationsInvertDisplayEveryXFrames = 2
const NotificationsSystemSendsDoubleNotificationMessages = true

// Discord module
const DiscordAuthToken = ""
const DiscordGuildID = ""
const DiscordUserID = ""
