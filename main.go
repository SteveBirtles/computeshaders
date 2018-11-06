package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	_ "image/png"
	"time"
)

const windowWidth = 1280
const windowHeight = 720

var (
	frames            = 0
	second            = time.Tick(time.Second)
	frameLength       float64
	windowTitlePrefix = "OpenGL Compute Shader"
	window            *glfw.Window
)

func main() {

	initiateOpenGL()

	prepareBuffers()
	prepareTextures()

	prepareShaders()

	gl.ClearColor(1, 1, 1, 1)

	projection := mgl32.Perspective(mgl32.DegToRad(50.0), float32(windowWidth/windowHeight), 1.0, 100.0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	for !window.ShouldClose() {

		frameStart := time.Now()

		processInputs()
		renderWorld()

		glfw.PollEvents()
		frames++
		select {
		case <-second:
			window.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			frames = 0
		default:
		}
		frameLength = time.Since(frameStart).Seconds()

	}

	glfw.Terminate()
}
