package screensaver

import (
	"image"
	"image/color"
	"math/rand"
	"pscreen/bridge/modules"
	"pscreen/config"
	"pscreen/utils"
)

type Particle struct {
	velocity utils.CoordsFloat
	position utils.CoordsFloat
	lifetime int
	front    bool
}

type State struct {
	Particles []Particle
}

var CurrentState State

func updateState(state State) State {
	newParticles := []Particle{}
	frontParticles := 0
	for _, particle := range CurrentState.Particles {
		if (particle.lifetime > 0) && (particle.position.X >= 0) && (particle.position.Y >= 0) && (particle.position.X <= float64(config.Config.CanvasRenderDimensions.X)) && (particle.position.Y <= float64(config.Config.CanvasRenderDimensions.Y)) {
			newParticles = append(newParticles, Particle{velocity: particle.velocity, position: particle.position, lifetime: particle.lifetime - 1, front: false})
			if particle.front {
				newParticles = append(newParticles, Particle{velocity: utils.CoordsFloat{X: particle.velocity.X, Y: particle.velocity.Y + config.Config.ScreensaverParticleGravity}, position: utils.CoordsFloat{X: particle.position.X + particle.velocity.X, Y: particle.position.Y + particle.velocity.Y}, lifetime: config.Config.ScreensaverParticleLifetime, front: true})
				frontParticles++
			}
		}
	}

	for i := 0; i < config.Config.ScreensaverParticlesToSpawnEachTime; i++ {
		if frontParticles < config.Config.ScreensaverMaxParticles {
			newParticles = append(newParticles, Particle{velocity: utils.CoordsFloat{X: utils.RandFloat64Around0() * config.Config.ScreensaverParticleSpeed.X, Y: utils.RandFloat64Around0() * config.Config.ScreensaverParticleSpeed.Y}, position: utils.CoordsFloat{X: rand.Float64() * float64(config.Config.CanvasRenderDimensions.X), Y: rand.Float64() * float64(config.Config.CanvasRenderDimensions.X)}, lifetime: config.Config.ScreensaverParticleLifetime, front: true})
			frontParticles++
		}
	}
	return State{Particles: newParticles}
}

var ScreensaverModule modules.Module = modules.Module{RenderFunction: func(im *image.RGBA) *image.RGBA {
	//fmt.Println("Rendering screensaver!")
	CurrentState = updateState(CurrentState)
	for _, particle := range CurrentState.Particles {
		im.Set(int(particle.position.X), int(particle.position.Y), color.RGBA{255, 255, 255, 255})
	}
	return im
}}
