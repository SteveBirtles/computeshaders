package main

import "github.com/go-gl/mathgl/mgl32"

var (
	nParticles     mgl32.Vec3
	totalParticles uint32

	time, deltaT, speed, angle float32
	particlesVao               uint32
	bhVao, bhBuf               uint32
	bh1, bh2                   mgl32.Vec4
)

func prepareBuffers() {

	initPos := make([]float32, 0)
	initVel := make([]float32, totalParticles*4)

	var p mgl32.Vec4

	dx := 2.0 / (nParticles.X() - 1)
	dy := 2.0 / (nParticles.Y() - 1)
	dz := 2.0 / (nParticles.Z() - 1)

	transf := mgl32.Translate3D(mgl32.Mat4(1.0), mgl32.Vec3{-1, -1, -1})

	for i := 0; i < nParticles.X(); i++ {
		for j := 0; j < nParticles.Y(); j++ {
			for k := 0; k < nParticles.Z(); k++ {
				p = mgl32.Vec4{dx * i, dy * j, dz * k, 1.0}
				p = transf * p
				initPos = append(initPos, p.X())
				initPos = append(initPos, p.Y())
				initPos = append(initPos, p.Z())
				initPos = append(initPos, p.W())
			}
		}
	}
}
