package screensaver

import (
	"image"
	"image/color"
	"math/rand"
	"pscreenapp/bridge/modules"
	"pscreenapp/config"
	"pscreenapp/utils"
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
	for _, particle := range CurrentState.Particles {
		if (particle.lifetime > 0) && (particle.position.X >= 0) && (particle.position.Y >= 0) && (particle.position.X <= float64(config.CanvasRenderDimensions.X)) && (particle.position.Y <= float64(config.CanvasRenderDimensions.Y)) {
			newParticles = append(newParticles, Particle{velocity: particle.velocity, position: particle.position, lifetime: particle.lifetime - 1, front: false})
			if particle.front {
				newParticles = append(newParticles, Particle{velocity: particle.velocity, position: utils.CoordsFloat{X: particle.position.X + particle.velocity.X, Y: particle.position.Y + particle.velocity.Y}, lifetime: config.ParticleLifetime, front: true})
			}
		}
	}

	for i := 0; i < config.ScreensaverParticlesToSpawnEachTime; i++ {
		if len(newParticles) < config.ScreensaverMaxParticles {
			newParticles = append(newParticles, Particle{velocity: utils.CoordsFloat{X: (rand.Float64() - 0.5) * 2, Y: rand.Float64()}, position: utils.CoordsFloat{X: rand.Float64() * float64(config.CanvasRenderDimensions.X), Y: 0}, lifetime: config.ParticleLifetime, front: true})
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
