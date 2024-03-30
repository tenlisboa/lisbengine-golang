package core

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type CameraMove uint

const (
	Forward CameraMove = iota
	Backward
	Left
	Right
)

type Camera struct {
	Position mgl32.Vec3
	Front    mgl32.Vec3
	Up       mgl32.Vec3
	Right    mgl32.Vec3
	WorldUp  mgl32.Vec3

	Yaw   float32
	Pitch float32

	MovementSpeed    float32
	MouseSensitivity float32
	Fov              float32

	firstMouseMove bool
	lastXMousePos  float32
	lastYMousePos  float32
}

func NewCamera(position, up mgl32.Vec3, yaw, pitch float32) *Camera {

	camera := &Camera{
		Position: position,
		WorldUp:  up,
		Yaw:      yaw,
		Pitch:    pitch,

		Fov: 45,

		firstMouseMove: true,
	}

	camera.updateCameraVectors()

	return camera
}

func (c *Camera) SetLastMousePos(lastx, lasty float32) {
	c.lastXMousePos = lastx
	c.lastYMousePos = lasty
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

func (c *Camera) Move(direction CameraMove, deltaTime float32) {
	velocity := c.MovementSpeed * deltaTime
	switch direction {
	case Forward:
		c.Position.Add(c.Front.Mul(velocity))
	case Backward:
		c.Position.Sub(c.Front.Mul(velocity))
	case Right:
		c.Position.Sub(c.Right.Mul(velocity))
	case Left:
		c.Position.Add(c.Right.Mul(velocity))
	default:
		log.Fatalf("invalid camera direction: %d\n", direction)
	}
}

func (c *Camera) Look(xoffset, yoffset float32, constrainPitch bool) {
	xoffset *= c.MouseSensitivity
	yoffset *= c.MouseSensitivity

	c.Yaw += xoffset
	c.Pitch += yoffset

	fmt.Println(c.Yaw, c.Pitch)

	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}

	c.updateCameraVectors()
}

func (c *Camera) Zoom(yoffset float32) {
	c.Fov -= yoffset
	if c.Fov < 1.0 {
		c.Fov = 1.0
	}
	if c.Fov > 45.0 {
		c.Fov = 45.0
	}
}

func (c *Camera) mouseCallback(w *glfw.Window, xpos, ypos float64) {
	if c.firstMouseMove {
		c.SetLastMousePos(float32(xpos), float32(ypos))
		c.firstMouseMove = false
	}

	xoffset := float32(xpos) - c.lastXMousePos
	yoffset := c.lastYMousePos - float32(ypos)

	c.SetLastMousePos(float32(xpos), float32(ypos))

	c.Look(xoffset, yoffset, true)
}

func (c *Camera) EnableFirstPersonControl() {
	glfw.GetCurrentContext().SetCursorPosCallback(c.mouseCallback)
}

func (c *Camera) updateCameraVectors() {
	c.Front = mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(c.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))),
	}.Normalize()

	c.Right = c.Front.Cross(c.WorldUp)
	c.Up = c.Right.Normalize().Cross(c.Front)

	fmt.Println(c.Yaw, c.Pitch, c.Front, c.Right, c.Up)

}
