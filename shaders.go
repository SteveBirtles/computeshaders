package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	_ "image/png"
)

const terminator = "\x00"

var (
	shaderRenderProgram  uint32
	shaderComputeProgram uint32
)

const vertexShader = `#version 430

layout (location = 0) in vec4 VertexPosition;

out vec3 Position;

uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

void main()
{
    Position = (ModelViewMatrix * VertexPosition).xyz;
    gl_Position = MVP * VertexPosition;
}`

const fragmentShader = `#version 430

in vec3 Position;

uniform vec4 Color;

layout( location = 0 ) out vec4 FragColor;

void main() {
  FragColor = Color;
}`

const computeShader = `#version 430

layout( local_size_x = 1000 ) in;

uniform float Gravity1 = 1000.0;
uniform vec3 BlackHolePos1 = vec3(5,0,0);

uniform float Gravity2 = 1000.0;
uniform vec3 BlackHolePos2 = vec3(-5,0,0);

uniform float ParticleMass = 0.1;
uniform float ParticleInvMass = 1.0 / 0.1;
uniform float DeltaT = 0.0005;
uniform float MaxDist = 45.0;

layout(std430, binding=0) buffer Pos {
  vec4 Position[];
};
layout(std430, binding=1) buffer Vel {
  vec4 Velocity[];
};

void main() {
  uint idx = gl_GlobalInvocationID.x;

  vec3 p = Position[idx].xyz;

  // Force from black hole #1
  vec3 d = BlackHolePos1 - p;
  float dist = length(d);
  vec3 force = (Gravity1 / dist) * normalize(d);

  // Force from black hole #2
  d = BlackHolePos2 - p;
  dist = length(d);
  force += (Gravity2 / dist) * normalize(d);

  // Reset particles that get too far from the attractors
  if( dist > MaxDist ) {
    Position[idx] = vec4(0,0,0,1);
  } else {
    // Apply simple Euler integrator
    vec3 a = force * ParticleInvMass;
    Position[idx] = vec4(
        p + Velocity[idx].xyz * DeltaT + 0.5 * a * DeltaT * DeltaT, 1.0);
    Velocity[idx] = vec4( Velocity[idx].xyz + a * DeltaT, 0.0);
  }
}`

func prepareShaders() {

	var err error

	shaderRenderProgram, err = newShaderRenderProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	shaderComputeProgram, err = newShaderComputeProgram(computeShader)
	if err != nil {
		panic(err)
	}

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 5000.0)
	projectionUniform := gl.GetUniformLocation(shaderRenderProgram, gl.Str("projection"+terminator))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(shaderRenderProgram, gl.Str("model"+terminator))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	tex := gl.GetUniformLocation(shaderRenderProgram, gl.Str("tex"+terminator))
	gl.Uniform1i(tex, 0)

	vertexIn := uint32(gl.GetAttribLocation(shaderRenderProgram, gl.Str("vertexIn"+terminator)))
	gl.EnableVertexAttribArray(vertexIn)
	gl.VertexAttribPointer(vertexIn, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))

	texCoordIn := uint32(gl.GetAttribLocation(shaderRenderProgram, gl.Str("texCoordIn"+terminator)))
	gl.EnableVertexAttribArray(texCoordIn)
	gl.VertexAttribPointer(texCoordIn, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))

	colourIn := uint32(gl.GetAttribLocation(shaderRenderProgram, gl.Str("colourIn"+terminator)))
	gl.EnableVertexAttribArray(colourIn)
	gl.VertexAttribPointer(colourIn, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(5*4))

	gl.BindFragDataLocation(shaderRenderProgram, 0, gl.Str("colourOut"+terminator))

}
