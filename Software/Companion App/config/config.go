package config

import (
	"errors"
	"os"
	"path/filepath"
	"pscreen/constants"
	"pscreen/utils"

	"github.com/kirsle/configdir"
	"gopkg.in/yaml.v3"
)

type config struct {
	UpdateUIEveryXMilliseconds int
	UIPaddingIndentAmount      int

	LoadedModules                        []string
	ModulePersistance                    bool
	ChangeModuleEveryXMilliseconds       int
	RenderDeviceScreenEveryXMilliseconds int
	TransitionMilliseconds               float64
	TransitionBlurringSigma              float64

	CanvasRenderDimensions utils.Coords
	RotateScreen180Deg     bool
	CircularScreenLayout   bool
	UseWallpaper           bool

	AutoStartXMit           bool
	DefaultPortSerialNumber string
	SerialPortBaudRate      int
	RGBXMit                 bool

	I18nLanguage string

	DebugSaveScreen bool

	// Weather module
	OpenWeatherMapAPIKey            string
	Lat                             float64
	Long                            float64
	UpdateWeatherEveryXMilliseconds int64
	WindIndicatorRadius             float64

	// Monitor module
	CPUUsageSamplingMilliseconds int
	CPUUsageBarMargin            int

	CPUUsageBarDimensions utils.Coords

	// Media module
	MediaProgressBarHeight          int
	MediaProgressBarWaveScale       float64
	MediaProgressBarIndicatorRadius int

	// osu! module
	OsuTrackedKeys           []uint16 // see https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
	OsuTrackedKeysLabels     []string
	OsuTrackedKeysDimensions utils.Coords

	// Pong module
	PongBallVelocity         float64
	PongBallRadius           float64
	PongPaddleWidth          float64
	PongPaddleInvisibleWidth float64
	PongPaddleDistFromEdge   float64
	PongPaddleLength         float64
	PongPaddleVelocity       float64
	PongPaddleP              float64
	PongPaddleI              float64
	PongPaddleD              float64
	PongTimestepsPerFrame    int

	// QRCode module
	QRCodeContent     string
	QRCodeTitle       string
	QRCodeDescription string

	// Screensaver module
	ScreensaverMaxParticles             int
	ScreensaverParticlesToSpawnEachTime int
	ScreensaverParticleLifetime         int
	ScreensaverParticleGravity          float64

	ScreensaverParticleSpeed utils.CoordsFloat

	// Notifications module
	UseNotificationsModule                             bool
	NotificationsDisplayForXMilliseconds               int64
	NotificationsInvertForXMilliseconds                int64
	NotificationsInvertDisplayEveryXFrames             int
	NotificationsSystemSendsDoubleNotificationMessages bool

	// Discord module
	DiscordAuthToken string
	DiscordGuildID   string
	DiscordUserID    string

	// DVD module
	DVDLogoVelocity float64

	// Visualizer module
	VisualizerUseMicrophone              bool
	VisualizerInputDelayMillis           int
	VisualizerSampleRate                 int
	VisualizerSampleBufferSize           int
	VisualizerCumulativeSampleBufferSize int
	VisualizerScaleSmoothing             float64 // 0-1
	VisualizerScale                      float64 // 0-1
	VisualizerOscilloscopeScale          float64
	VisualizerMinScale                   float64
	VisualizerFFTCutoff                  float64
	VisualizerFFTScalingBase             float64
	VisualizerFFTScalingStartValue       float64
	VisualizerFFTBarWidth                int
	VisualizerFFTBarSpacing              int
	VisualizerFFTSmoothing               float64 // 0-1
	VisualizerShowFFT                    bool
}

var DefaultConfig = config{
	UpdateUIEveryXMilliseconds: 500,
	UIPaddingIndentAmount:      1,

	LoadedModules:                        []string{},
	ModulePersistance:                    true,
	ChangeModuleEveryXMilliseconds:       5000,
	RenderDeviceScreenEveryXMilliseconds: 0,
	TransitionMilliseconds:               400,
	TransitionBlurringSigma:              10,

	CanvasRenderDimensions: utils.Coords{X: 256, Y: 64},
	RotateScreen180Deg:     false,
	CircularScreenLayout:   false,
	UseWallpaper:           true,

	AutoStartXMit:           false,
	DefaultPortSerialNumber: "",
	SerialPortBaudRate:      921600,
	RGBXMit:                 false,

	I18nLanguage: "en_US",

	DebugSaveScreen: false,

	// Weather module
	OpenWeatherMapAPIKey:            "",
	Lat:                             0.0,
	Long:                            0.0,
	UpdateWeatherEveryXMilliseconds: 240000,
	WindIndicatorRadius:             6,

	// Monitor module
	CPUUsageSamplingMilliseconds: 1000,
	CPUUsageBarMargin:            1,

	CPUUsageBarDimensions: utils.Coords{X: 4, Y: 32},

	// Media module
	MediaProgressBarHeight:          8,
	MediaProgressBarWaveScale:       0.1,
	MediaProgressBarIndicatorRadius: 4,

	// osu! module
	OsuTrackedKeys:           []uint16{0x57, 0x58},
	OsuTrackedKeysLabels:     []string{"w", "x"},
	OsuTrackedKeysDimensions: utils.Coords{X: 32, Y: 48},

	// Pong module
	PongBallVelocity:         2.0,
	PongBallRadius:           2.0,
	PongPaddleWidth:          2.0,
	PongPaddleInvisibleWidth: 10.0,
	PongPaddleDistFromEdge:   10.0,
	PongPaddleLength:         20.0,
	PongPaddleVelocity:       0.2,
	PongPaddleP:              0.4,
	PongPaddleI:              0.1,
	PongPaddleD:              0.2,
	PongTimestepsPerFrame:    4,

	// QRCode module
	QRCodeContent:     "WIFI:S:SSID;T:WPA;P:PASSWORD;H:false;;",
	QRCodeTitle:       "SSID",
	QRCodeDescription: "PASSWORD",

	// Screensaver module
	ScreensaverMaxParticles:             10,
	ScreensaverParticlesToSpawnEachTime: 5,
	ScreensaverParticleLifetime:         10,
	ScreensaverParticleGravity:          0.1,

	ScreensaverParticleSpeed: utils.CoordsFloat{X: 2, Y: 2},

	// Notifications module
	UseNotificationsModule:                             true,
	NotificationsDisplayForXMilliseconds:               7000,
	NotificationsInvertForXMilliseconds:                1000,
	NotificationsInvertDisplayEveryXFrames:             2,
	NotificationsSystemSendsDoubleNotificationMessages: true,

	// Discord module
	DiscordAuthToken: "TOKEN",
	DiscordGuildID:   "XXXXXXXXXXXXXXXXXX",
	DiscordUserID:    "XXXXXXXXXXXXXXXXXX",

	// DVD module
	DVDLogoVelocity: 1.0,

	// Visualizer module
	VisualizerUseMicrophone:              false,
	VisualizerInputDelayMillis:           0,
	VisualizerSampleRate:                 192000,
	VisualizerSampleBufferSize:           1024,
	VisualizerCumulativeSampleBufferSize: 32768,
	VisualizerScaleSmoothing:             0.99, // 0-1
	VisualizerScale:                      0.7,  // 0-1
	VisualizerOscilloscopeScale:          4.0,
	VisualizerMinScale:                   10,
	VisualizerFFTCutoff:                  0.014,
	VisualizerFFTScalingBase:             1.06,
	VisualizerFFTScalingStartValue:       20.0,
	VisualizerFFTBarWidth:                4,
	VisualizerFFTBarSpacing:              1,
	VisualizerFFTSmoothing:               0.1, // 0-1
	VisualizerShowFFT:                    true,
}

var Config config

func ModuleConfigNamesToIDs(moduleNames []string) []int {
	IDs := []int{}
	for _, moduleName := range moduleNames {
		if x, found := constants.ModuleNames.Get(moduleName); found {
			IDs = append(IDs, x)
		} else {
			utils.CheckError(errors.New("config: unknown module name"))
		}
	}
	return IDs
}

func ModuleIDsToConfigNames(moduleIDs []int) []string {
	moduleNames := []string{}
	for _, moduleID := range moduleIDs {
		if x, found := constants.ModuleNames.GetInverse(moduleID); found {
			moduleNames = append(moduleNames, x)
		} else {
			utils.CheckError(errors.New("config: unknown module id"))
		}
	}
	return moduleNames
}

func getConfigPath() string {
	configPath := configdir.LocalConfig("ontake", "pscreen")
	configFilePath := filepath.Join(configPath, "config.yml")
	utils.CheckError(configdir.MakePath(configPath))
	return configFilePath
}

func SaveConfig() {
	configFilePath := getConfigPath()
	outConfig, err := yaml.Marshal(Config)
	utils.CheckError(err)
	err = os.WriteFile(configFilePath, outConfig, 0644)
	utils.CheckError(err)
}

func ParseConfig() {
	configFilePath := getConfigPath()
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		Config = DefaultConfig
	} else {
		utils.CheckError(err)
		inConfig, err := os.ReadFile(configFilePath)
		utils.CheckError(err)
		err = yaml.Unmarshal(inConfig, &Config)
		utils.CheckError(err)
	}
	SaveConfig()
}
