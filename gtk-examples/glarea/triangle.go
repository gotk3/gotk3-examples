package main

import (
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	triVerts  = 3 // The number of vertices in a triangle.
	floatSize = 4 // The size of a float32 in bytes.

	positionAttribute = 1
	colorAttribute    = 1

	positionElements = 3 // The number of floats describing a position(x, y, z).
	colorElements    = 3 // The number of floats describing a color(r, g, b).

	/* Each vertex in the the triangle has two attributes position and color.
	 * Offsets enable the right geometry and colors rendered.
	 * vertex [ position[x, y, z], color[r, g, b] ]
	 * float index [ 0 1 2 3 4 5 ]
	 * float value [ x y z r g b ]
	 * offfset:      ^     ^
	 *   postion ____|     |
	 *   color   __________|
	 */
	positionOffset = 0
	colorOffset    = positionElements * floatSize
)

// Window dimesnsions
var winWidth, winHeight int

type Triangle struct {
	attributes int
	data       []float32
	angle      float32
}

func (t *Triangle) Draw() {
	// Load the shader program into the rendering pipeline.
	gl.UseProgram(program)

	t.mvp()

	// Bind to the data in the buffer
	gl.BindVertexArray(vao)

	// Render the data
	gl.DrawArrays(gl.TRIANGLES, 0, int32(triVerts))

	// Done with the buffer and program so unbind them
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

func (t *Triangle) Data() unsafe.Pointer {
	// Return the address of the array containing all of the vertex data.
	return gl.Ptr(t.data)
}

func (t *Triangle) Stride() int32 {
	// Return the total number of bytes of data that describes each vertex.
	return int32(triVerts * t.attributes * floatSize)
}

func (t *Triangle) Size() int {
	// Return size of the data in number of bytes.
	return len(t.data) * floatSize
}

func (t *Triangle) SetAngle(a float32) {
	t.angle = a
}

func (t *Triangle) mvp() {
	// Get 4x4 identity matrix for the model's transformations
	model := mgl32.Ident4()

	// Apply the change in angle to the model's set of transformations
	model = mgl32.HomogRotate3DY(t.angle)

	// Set the handle to point to the address of the model matrix.
	gl.UniformMatrix4fv(modelIndex, 1, false, &model[0])

	// Get 4x4 projection matrix with a 60 degree field of view, an aspect ratio
	// of the window dimensions, near clipping plane, and a far clipping plane.
	projection := mgl32.Perspective(
		mgl32.DegToRad(40.0), float32(winWidth/winHeight), 0.1, -1.0,
	)
	// Set the handle to point to the address of the projection matrix.
	gl.UniformMatrix4fv(projectionIndex, 1, false, &projection[0])

	// Get 4x4 view matrix with an eye position, target position,
	// and the up direction with a positive bias in the y-axis.
	// Right-handed coordinate system.
	view := mgl32.LookAtV(
		mgl32.Vec3{0, 1, 2}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0},
	)
	// Set the handle to point to the address of the view matrix.
	gl.UniformMatrix4fv(viewIndex, 1, false, &view[0])
}
