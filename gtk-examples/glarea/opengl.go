package main

import (
	"errors"
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	MILLI_SEC_PER_SEC = float32(1000.0)
	DEG_PER_TWO_SEC   = MILLI_SEC_PER_SEC * 360.0
)

/*
 * Global variables
 */

var (
	triangle = Triangle{
		data: []float32{
			+0.0, +0.5, +0.0, // Top
			+1.0, +0.0, +0.0, // Red

			+0.5, -0.5, +0.0, // Bottom Right
			+0.0, +1.0, +0.0, // Green

			-0.5, -0.5, +0.0, // Bottom Left
			+0.0, +0.0, +1.0, // Blue
		},
		attributes: positionAttribute + colorAttribute,
		angle:      0.0,
	}

	modelIndex, viewIndex, projectionIndex int32
	vao, positionIndex, colorIndex         uint32
	epoch                                  int64
)

/*
 * Set up the OpenGL state
 */

func realize(glarea *gtk.GLArea) {
	log.Println("realize")

	// Make the the GLArea's GdkGLContext the current OpenGL context.
	glarea.MakeCurrent()

	// Initialize OpenGL.
	err := gl.Init()
	errorCheck(err)

	// Initialize shaders.
	err = initShaders()
	errorCheck(err)

	// Initialize buffer.
	vao, err = initBuffer()
	errorCheck(err)

	// Enable depth test.
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// Bind our callback function to the glarea, the callback with be called 60
	// times in one second.
	glarea.AddTickCallback(update, uintptr(0))
	// Set the epoch.
	epoch = glarea.GetFrameClock().GetFrameTime()

	// Log out OpenGL version and window dimensions.
	log.Printf("opengl version %s", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Printf("window - width: %d height: %d", winWidth, winHeight)
}

/*
 * Clean up the OpenGL state
 */

func unrealize(glarea *gtk.GLArea) {
	log.Println("unrealize")

	// Make sure the area that is being cleaned up is the current context
	// otherwise opengl api calls with panic.
	glarea.MakeCurrent()

	// Clean up all of the allocated data.
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteProgram(program)
}

/*
 * Render the OpenGL state
 */

func render(glarea *gtk.GLArea) bool {
	// log.Println("render")

	// Enable attribute index 0 as being used.
	gl.EnableVertexAttribArray(0)

	// Set background color (rgba).
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	triangle.Draw()

	// Flush the contents of the pipeline.
	gl.Flush()

	return true
}

/*
 * Helper functions
 */

func initBuffer() (vao uint32, e error) {
	// Allocate, assign, and bind a single Vertex Array Object.
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Allocate and assign Vertex Buffer Object.
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	// Bind Vertex Buffer Object as being the active buffer and storing vertex
	// attributes(coordinates and color data).
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// Copy the vertex data from array to our buffer.
	gl.BufferData(
		gl.ARRAY_BUFFER, triangle.Size(),
		triangle.Data(), gl.STATIC_DRAW,
	)

	// Enable and set the position attribute with the local data.
	gl.EnableVertexAttribArray(positionIndex)
	gl.VertexAttribPointer(positionIndex, 3, gl.FLOAT, false, triangle.Stride(),
		gl.PtrOffset(positionOffset))

	// Enable and set the color attribute with the local data.
	gl.EnableVertexAttribArray(colorIndex)
	gl.VertexAttribPointer(colorIndex, 3, gl.FLOAT, false, triangle.Stride(),
		gl.PtrOffset(colorOffset))

	// Finished loading data unbind it vao.
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	// No longer need the vbo.
	gl.DeleteBuffers(1, &vbo)

	if vao == 0 {
		e = errors.New("Error initializing buffers.")
	}
	return
}

func initShaders() (err error) {
	// Configure the vertex and fragment shaders.
	program, err = newProgram(vertexShader, fragmentShader)

	// Get the location of the "position" and "color" attributes.
	positionIndex = uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	colorIndex = uint32(gl.GetAttribLocation(program, gl.Str("color\x00")))

	// Get the location of the "model", "view", "projection" uniforms.
	modelIndex = gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewIndex = gl.GetUniformLocation(program, gl.Str("view\x00"))
	projectionIndex = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	return
}

func update(widget *gtk.Widget, frameClock *gdk.FrameClock, userData uintptr) bool {
	// Calculate the delta time.
	delta := float32(frameClock.GetFrameTime() - epoch)

	// Set the angle of rotation of the triangle based on the delta.
	triangle.SetAngle(delta / DEG_PER_TWO_SEC)

	// Log out the delta time and the angle.
	// log.Printf("update delta %v angle %v", delta, triangle.angle)

	// Queue up the re-rendering of the GLArea.
	widget.QueueDraw()
	return true
}
