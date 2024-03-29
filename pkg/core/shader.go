package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	ID uint32
}

func NewShader(vertexShaderPath, fragmentShaderPath string) *Shader {
	shaderProgram := createShaderProgram(vertexShaderPath, fragmentShaderPath)

	return &Shader{
		ID: shaderProgram,
	}
}

func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

func (s *Shader) SetBool(name string, value bool) {
	var v int32
	if value {
		v = 1
	} else {
		v = 0
	}
	gl.Uniform1i(int32(s.ID), v)
}

func (s *Shader) SetInt(name string, value int32) {
	loc := gl.GetUniformLocation(s.ID, gl.Str(fmt.Sprintf("%s\x00", name)))
	gl.Uniform1i(loc, value)
}

func (s *Shader) SetFloat(name string, value float32) {
	loc := gl.GetUniformLocation(s.ID, gl.Str(fmt.Sprintf("%s\x00", name)))
	gl.Uniform1f(loc, value)
}

func (s *Shader) SetMatrix(name string, value mgl32.Mat4) {
	loc := gl.GetUniformLocation(s.ID, gl.Str(fmt.Sprintf("%s\x00", name)))
	gl.UniformMatrix4fv(loc, 1, false, &value[0])
}

func getShaderFromFile(path string) string {
	shader, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error opening the shader file %s: %e\n", path, err)
	}

	return string(shader)
}

func compileShader(path string, shaderType uint32) uint32 {
	source := getShaderFromFile(path)
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

		logStr := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logStr))

		log.Fatalf("failed to compile %v: %v", source, logStr)
		return 0
	}

	return shader

}

func createShaderProgram(vertexShaderPath, fragmentShaderPath string) uint32 {
	vertexShader := compileShader(vertexShaderPath, gl.VERTEX_SHADER)

	fragmentShader := compileShader(fragmentShaderPath, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		logStr := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logStr))

		log.Fatalf("failed to link program: %v", logStr)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}
