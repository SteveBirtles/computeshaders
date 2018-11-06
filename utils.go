package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
	"strings"
)

func init() {

	runtime.LockOSThread()

}

func initiateOpenGL() {

	var err error
	if err = glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err = glfw.CreateWindow(windowWidth, windowHeight, windowTitlePrefix, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	window.SetCursorPos(windowWidth/2, windowHeight/2)
	window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

}

func newShaderComputeProgram(computeShaderSource string) (unit32, error) {

	computeShader, err := compileShader(computeShaderSource+terminator, gl.COMPUTE_SHADER)
	if err != nil {
		return 0, err
	}

	shaderRenderProgram = gl.CreateProgram()

	gl.AttachShader(shaderRenderProgram, computeShader)
	gl.LinkProgram(shaderRenderProgram)

	var status int32
	gl.GetProgramiv(shaderRenderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderRenderProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat(terminator, int(logLength+1))
		gl.GetProgramInfoLog(shaderRenderProgram, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link compute shader program: %v", log)
	}

	gl.DeleteShader(computeShader)

	return shaderRenderProgram, nil

}

func newShaderRenderProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {

	vertexShader, err := compileShader(vertexShaderSource+terminator, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource+terminator, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	shaderRenderProgram = gl.CreateProgram()

	gl.AttachShader(shaderRenderProgram, vertexShader)
	gl.AttachShader(shaderRenderProgram, fragmentShader)
	gl.LinkProgram(shaderRenderProgram)

	var status int32
	gl.GetProgramiv(shaderRenderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderRenderProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat(terminator, int(logLength+1))
		gl.GetProgramInfoLog(shaderRenderProgram, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link render shader program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shaderRenderProgram, nil

}

func compileShader(source string, shaderType uint32) (uint32, error) {

	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat(terminator, int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil

}
