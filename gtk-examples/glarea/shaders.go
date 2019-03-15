package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

/* Shaders from
 * https://www.khronos.org/opengl/wiki/Tutorial2:_VAOs,_VBOs,_Vertex_and_Fragment_Shaders_(C_/_SDL)#tutorial2.vert
 * https://www.khronos.org/opengl/wiki/Tutorial2:_VAOs,_VBOs,_Vertex_and_Fragment_Shaders_(C_/_SDL)#tutorial2.frag
 */

const (
	vertexShader = `
#version 330

// vertex position and color data
in vec3 position;
in vec3 color;

// mvp matrices 
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

// output the out_color variable to the next shader in the chain
out vec3 f_color;

void main(void) {
  gl_Position = projection * view * model * vec4(position, 1.0);

  // gl_Position = vec4(position, 1.0);

  // pass the color through unmodified
  f_color = color;
}
` + "\x00"

	fragmentShader = `
#version 330

// It was expressed that some drivers required this next line
// to function properly
precision highp float;

in vec3 f_color;

void main(void) {
    // Pass through original color with full opacity.
    gl_FragColor = vec4(f_color,1.0);
}
` + "\x00"
)

// Shader program handle
var program uint32

/* Shader compiler functions from
 * https://github.com/go-gl/example/blob/master/gl41core-cube/cube.go
 */

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	gpuProgram := gl.CreateProgram()

	gl.AttachShader(gpuProgram, vertexShader)
	gl.AttachShader(gpuProgram, fragmentShader)
	gl.LinkProgram(gpuProgram)

	var status int32
	gl.GetProgramiv(gpuProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(gpuProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(gpuProgram, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return gpuProgram, nil
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

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
