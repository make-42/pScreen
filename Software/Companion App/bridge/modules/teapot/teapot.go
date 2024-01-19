package cube

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"pscreen/bridge/modules"
	"pscreen/bridge/renderer"
	"pscreen/config"
	"time"

	"github.com/make-42/cpu3d/camera"
	"github.com/make-42/cpu3d/render"
	"github.com/make-42/cpu3d/transform"
	"github.com/make-42/cpu3d/utils"
)

var startTime = time.Now()

/*
	var pointCloudOrigin = []utils.SpaceCoords{
		{X: 0, Y: 0, Z: 0},
		{X: 0, Y: 0, Z: 1},
		{X: 0, Y: 1, Z: 0},
		{X: 0, Y: 1, Z: 1},
		{X: 1, Y: 0, Z: 0},
		{X: 1, Y: 0, Z: 1},
		{X: 1, Y: 1, Z: 0},
		{X: 1, Y: 1, Z: 1}}

	var pointCloudCentered = transform.PointCloudTranslateWorldCoords(&pointCloudOrigin, &utils.SpaceCoords{
		X: -0.5,
		Y: -0.5,
		Z: -0.5,
	})
*/
var TeapotModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	draw.Draw(im, im.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)
	worldCamera := camera.Camera{
		FocalLength: 0.06, //m
		SensorSize: utils.PlaneCoords{
			X: float64(config.Config.CanvasRenderDimensions.X) * 10e-4, //m
			Y: float64(config.Config.CanvasRenderDimensions.Y) * 10e-4, //m
		},
	}
	displayResolution := utils.IntPair{X: config.Config.CanvasRenderDimensions.X, Y: config.Config.CanvasRenderDimensions.Y}

	i := float64(time.Since(startTime).Milliseconds()) / 100
	/*
		pointCloudRotated := transform.PointCloudRotateWorldCoords(&pointCloudCentered, &utils.RotationCoords{
			X: float64(i) / 8,
			Y: float64(i) / 16,
			Z: float64(i) / 32,
		})

		pointCloud := transform.PointCloudTranslateWorldCoords(&pointCloudRotated, &utils.SpaceCoords{
			X: 0,
			Y: 0,
			Z: 2,
		})
		//screenCoords := camera.PointCloudWorldCoordsToSensorCoords(&pointCloud, &worldCamera)
		//render.RenderScreenPoints(4, displayResolution, &screenCoords, "output-pointcloud.png")

		edges := []utils.SpaceEdge{
			{A: pointCloud[0], B: pointCloud[1]},
			{A: pointCloud[0], B: pointCloud[2]},
			{A: pointCloud[0], B: pointCloud[4]},
			{A: pointCloud[7], B: pointCloud[3]},
			{A: pointCloud[7], B: pointCloud[5]},
			{A: pointCloud[7], B: pointCloud[6]},

			{A: pointCloud[1], B: pointCloud[5]},
			{A: pointCloud[1], B: pointCloud[3]},
			{A: pointCloud[2], B: pointCloud[6]},
			{A: pointCloud[2], B: pointCloud[3]},
			{A: pointCloud[4], B: pointCloud[6]},
			{A: pointCloud[4], B: pointCloud[5]},
		}

		lines := camera.EdgesWorldCoordsToSensorCoords(&edges, &worldCamera)
	*/
	edgesTeapotCentered := transform.EdgesTranslateWorldCoords(&renderer.TeapotModel, &utils.SpaceCoords{X: 0, Y: 0, Z: -3})
	edgesTeapotRotatedFirst := transform.EdgesRotateWorldCoords(&edgesTeapotCentered, &utils.RotationCoords{
		X: math.Pi * i / 100,
		Y: 0,
		Z: math.Pi / 4 * math.Sin(i/10),
	})
	edgesTeapotRotated := transform.EdgesRotateWorldCoords(&edgesTeapotRotatedFirst, &utils.RotationCoords{X: 0, Y: 0, Z: -math.Pi / 4})

	edgesTeapotScaled := transform.EdgesScaleWorldCoords(&edgesTeapotRotated, 0.2)
	edges := transform.EdgesTranslateWorldCoords(&edgesTeapotScaled, &utils.SpaceCoords{
		X: 0,
		Y: 0,
		Z: 2 + math.Sin(i/10),
	})
	lines := camera.EdgesWorldCoordsToSensorCoords(&edges, &worldCamera)
	render.RenderScreenLinesToImage(im, color.RGBA{255, 255, 255, 255}, displayResolution, &lines)
	return im
}}
