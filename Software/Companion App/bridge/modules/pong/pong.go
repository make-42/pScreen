package pong

import (
	"image"
	"math"
	"math/rand"
	"pscreenapp/bridge/modules"
	"pscreenapp/bridge/renderer"
	"pscreenapp/config"
	"pscreenapp/utils"
	"time"

	"git.sr.ht/~sbinet/gg"
	"go.einride.tech/pid"
)

type Paddle struct {
	Position utils.CoordsFloat
	PID      pid.Controller
}

type Ball struct {
	Position utils.CoordsFloat
	Velocity utils.CoordsFloat
}

type State struct {
	Ball    Ball
	Paddles [2]Paddle
	Break   bool
}

var CurrentState State
var IsInit = false

func predictPaddleY(state State, paddleIndex int) float64 {
	for !(state.Ball.Position.X < config.Config.PongPaddleDistFromEdge+config.Config.PongBallRadius || state.Ball.Position.X > float64(config.Config.CanvasRenderDimensions.X)-config.Config.PongPaddleDistFromEdge-config.Config.PongBallRadius) && !state.Break {
		state = updateState(state, false, false, paddleIndex)
	}
	return state.Ball.Position.Y
}

func updateState(state State, predictPaddleMovement bool, resetStateOnGameOver bool, paddleIndex int) State {
	state.Ball.Position.X += state.Ball.Velocity.X
	state.Ball.Position.Y += state.Ball.Velocity.Y
	if state.Ball.Position.Y < config.Config.PongBallRadius || state.Ball.Position.Y > float64(config.Config.CanvasRenderDimensions.Y)-config.Config.PongBallRadius {
		state.Ball.Velocity.Y = -state.Ball.Velocity.Y
	}
	if predictPaddleMovement {
		for i, paddle := range CurrentState.Paddles {
			if (i == 0 && state.Ball.Velocity.X < 0) || (i == 1 && state.Ball.Velocity.X > 0) {
				paddlePredictedY := predictPaddleY(state, i)
				paddlePredictedY = math.Max(math.Min(float64(config.Config.CanvasRenderDimensions.Y)-config.Config.PongPaddleLength/2, paddlePredictedY), config.Config.PongPaddleLength/2)
				state.Paddles[i].PID.Update(pid.ControllerInput{
					ReferenceSignal:  paddlePredictedY,
					ActualSignal:     paddle.Position.Y,
					SamplingInterval: time.Millisecond * time.Duration(1000*4/config.Config.PongTimestepsPerFrame),
				})
				state.Paddles[i].Position.Y = state.Paddles[i].Position.Y + 4/float64(config.Config.PongTimestepsPerFrame)*config.Config.PongPaddleVelocity*state.Paddles[i].PID.State.ControlSignal
				var collisionPaddlePosition utils.CoordsFloat
				if i == 0 {
					collisionPaddlePosition = utils.CoordsFloat{X: paddle.Position.X - config.Config.PongPaddleInvisibleWidth - config.Config.PongPaddleWidth/2, Y: paddle.Position.Y - config.Config.PongPaddleLength/2}
				} else {
					collisionPaddlePosition = utils.CoordsFloat{X: paddle.Position.X - config.Config.PongPaddleWidth/2, Y: paddle.Position.Y - config.Config.PongPaddleLength/2}
				}
				if utils.IsPointInRectangle(collisionPaddlePosition, config.Config.PongPaddleWidth+config.Config.PongPaddleInvisibleWidth, config.Config.PongPaddleLength, state.Ball.Position, config.Config.PongBallRadius) {
					state.Ball.Velocity.X = -state.Ball.Velocity.X
					state.Ball.Position.X = math.Max(math.Min(state.Ball.Position.X, float64(config.Config.CanvasRenderDimensions.X)-config.Config.PongPaddleDistFromEdge-config.Config.PongPaddleWidth/2), config.Config.PongPaddleDistFromEdge+config.Config.PongPaddleWidth/2)
				}
			}
		}
	} else {
		bounce := false
		if state.Ball.Position.X > float64(config.Config.CanvasRenderDimensions.X)-config.Config.PongPaddleDistFromEdge-config.Config.PongPaddleWidth/2 {
			bounce = true
			if paddleIndex == 1 {
				state.Break = true
			}
		}
		if state.Ball.Position.X < config.Config.PongPaddleDistFromEdge+config.Config.PongPaddleWidth/2 {
			bounce = true
			if paddleIndex == 0 {
				state.Break = true
			}
		}
		if bounce {
			state.Ball.Velocity.X = -state.Ball.Velocity.X
			state.Ball.Position.X = math.Max(math.Min(state.Ball.Position.X, float64(config.Config.CanvasRenderDimensions.X)-config.Config.PongPaddleDistFromEdge-config.Config.PongPaddleWidth/2), config.Config.PongPaddleDistFromEdge+config.Config.PongPaddleWidth/2)
		}
	}
	if resetStateOnGameOver {
		if state.Ball.Position.X < 0 || state.Ball.Position.X > float64(config.Config.CanvasRenderDimensions.X) {
			return resetState()
		}
	}
	return state
}

func resetState() State {
	//angle := utils.RandFloat64Around0() * 2 * math.Pi
	angle := 2 * math.Pi / 360 * 40
	return State{
		Ball: Ball{
			Position: utils.CoordsFloat{
				X: float64(config.Config.CanvasRenderDimensions.X) / 2,
				Y: float64(config.Config.CanvasRenderDimensions.Y) / 2,
			},
			Velocity: utils.CoordsFloat{
				X: math.Cos(angle) * config.Config.PongBallVelocity * 4 / float64(config.Config.PongTimestepsPerFrame),
				Y: -math.Sin(angle) * config.Config.PongBallVelocity * 4 / float64(config.Config.PongTimestepsPerFrame),
			},
		},
		Paddles: [2]Paddle{
			{
				Position: utils.CoordsFloat{
					X: config.Config.PongPaddleDistFromEdge,
					Y: rand.Float64()*(float64(config.Config.CanvasRenderDimensions.Y)-config.Config.PongPaddleLength) + config.Config.PongPaddleLength/2,
				},
				PID: pid.Controller{
					Config: pid.ControllerConfig{
						ProportionalGain: config.Config.PongPaddleP,
						IntegralGain:     config.Config.PongPaddleI,
						DerivativeGain:   config.Config.PongPaddleD,
					},
				},
			},
			{
				Position: utils.CoordsFloat{
					X: float64(config.Config.CanvasRenderDimensions.X) - config.Config.PongPaddleDistFromEdge,
					Y: rand.Float64()*(float64(config.Config.CanvasRenderDimensions.Y)-config.Config.PongPaddleLength) + config.Config.PongPaddleLength/2,
				},
				PID: pid.Controller{
					Config: pid.ControllerConfig{
						ProportionalGain: config.Config.PongPaddleP,
						IntegralGain:     config.Config.PongPaddleI,
						DerivativeGain:   config.Config.PongPaddleD,
					},
				},
			},
		},
		Break: false,
	}
}

var PongModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	if !IsInit {
		CurrentState = resetState()
		IsInit = true
	}
	for i := 0; i < config.Config.PongTimestepsPerFrame; i++ {
		CurrentState = updateState(CurrentState, true, true, -1)
	}
	dc := gg.NewContextForRGBA(im)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(CurrentState.Ball.Position.X-config.Config.PongBallRadius, CurrentState.Ball.Position.Y-config.Config.PongBallRadius, config.Config.PongBallRadius*2, config.Config.PongBallRadius*2)
	dc.Fill()
	for _, paddle := range CurrentState.Paddles {
		dc.DrawRectangle(paddle.Position.X-config.Config.PongPaddleWidth/2, paddle.Position.Y-config.Config.PongPaddleLength/2, config.Config.PongPaddleWidth, config.Config.PongPaddleLength)
		dc.Fill()
	}
	return renderer.RemoveAntiAliasing(im)
}}
