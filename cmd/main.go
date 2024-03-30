package main

import (
	"fmt"

	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tenlisboa/lisbengine/pkg/core"
)

const windowWidth = 800
const windowHeight = 600

var camera *core.Camera

func init() {
	runtime.LockOSThread()

	camera = core.NewCamera(mgl32.Vec3{0.0, 0.0, 3.0}, mgl32.Vec3{0.0, 1.0, 0.0})
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	shader := core.NewShader("assets/shaders/shader.vert", "assets/shaders/shader.frag")

	camera.EnableFirstPersonControl()
	camera.SetLastMousePos(windowWidth/float32(2.0), windowHeight/float32(2.0))

	texture1 := core.NewTexture(0)
	texture1.SetWrapX(core.Repeat)
	texture1.SetWrapY(core.Repeat)
	texture1.SetMinFilter(core.LinearMipmapLinear)
	texture1.SetMagFilter(core.Linear)
	texture1.LoadImage("assets/images/container.jpg")

	texture2 := core.NewTexture(1)
	texture2.SetWrapX(core.Repeat)
	texture2.SetWrapY(core.Repeat)
	texture2.SetMinFilter(core.LinearMipmapLinear)
	texture2.SetMagFilter(core.Linear)
	texture2.LoadImage("assets/images/awesomeface.png")

	shader.Use()
	shader.SetInt("texture1", 0)
	shader.SetInt("texture2", 1)

	// Configure the vertex data
	var vao uint32
	var vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	defer func() {
		gl.DeleteVertexArrays(1, &vao)
		gl.DeleteBuffers(1, &vbo)
	}()

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	cubePositions := []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5}}

	angle := 0.0
	lastTick := glfw.GetTime()
	delta := 0.0

	for !window.ShouldClose() {
		processInput(window, delta)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		delta = time - lastTick
		lastTick = time

		texture1.Use()
		texture2.Use()
		shader.Use()

		angle += delta
		view := camera.GetViewMatrix()
		view = view.Mul4(mgl32.Translate3D(0, 0, -3))
		projection := mgl32.Perspective(mgl32.DegToRad(camera.Fov), float32(windowWidth/windowHeight), 0.1, 100)

		shader.SetMatrix("view", view)
		shader.SetMatrix("projection", projection)

		gl.BindVertexArray(vao)

		for i, pos := range cubePositions {
			model := mgl32.Ident4()
			model = model.Mul4(mgl32.Translate3D(pos[0], pos[1], pos[2]))
			angle := float32(20.0 * (i + 1))
			model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(angle), mgl32.Vec3{1, 0.3, 0}))
			shader.SetMatrix("model", model)

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		gl.BindVertexArray(0)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func isPressing(key glfw.Key) bool {
	return glfw.GetCurrentContext().GetKey(key) == glfw.Press
}

func processInput(window *glfw.Window, delta float64) {
	if isPressing(glfw.KeyEscape) {
		window.SetShouldClose(true)
		glfw.Terminate()
	}
	if isPressing(glfw.KeyW) {
		fmt.Println("Foward: ", delta)
		camera.Move(core.Forward, float32(delta))
	}
	if isPressing(glfw.KeyS) {
		camera.Move(core.Backward, float32(delta))
	}
	if isPressing(glfw.KeyA) {
		camera.Move(core.Left, float32(delta))
	}
	if isPressing(glfw.KeyD) {
		camera.Move(core.Right, float32(delta))
	}
}

var cubeVertices = []float32{
	-0.5, -0.5, -0.5, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	-0.5, 0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0,

	-0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 1.0,
	-0.5, 0.5, 0.5, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,

	-0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, -0.5, 1.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,
	-0.5, 0.5, 0.5, 1.0, 0.0,

	0.5, 0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,

	-0.5, 0.5, -0.5, 0.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0}
