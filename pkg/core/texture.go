package core

import (
	"image"
	"image/draw"
	"log"
	"os"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type WrapParams int32

const (
	Repeat       WrapParams = gl.REPEAT
	MirrorRepeat WrapParams = gl.MIRRORED_REPEAT
	ClampEdge    WrapParams = gl.CLAMP_TO_EDGE
	ClampBorder  WrapParams = gl.CLAMP_TO_BORDER
)

type FilterParams int32

const (
	Nearest              FilterParams = gl.NEAREST
	Linear               FilterParams = gl.LINEAR
	NearestMipmapNearest FilterParams = gl.NEAREST_MIPMAP_NEAREST
	LinearMipmapNearest  FilterParams = gl.LINEAR_MIPMAP_NEAREST
	NearestMipmapLinear  FilterParams = gl.NEAREST_MIPMAP_LINEAR
	LinearMipmapLinear   FilterParams = gl.LINEAR_MIPMAP_LINEAR
)

type Texture struct {
	sync.Mutex
	ID             uint32
	textureAddress uint32
}

func NewTexture(texIndex uint32) *Texture {
	var ID uint32
	gl.GenTextures(1, &ID)
	textureAddress := gl.TEXTURE0 + uintptr(texIndex)

	tex := &Texture{
		ID:             ID,
		textureAddress: uint32(textureAddress),
	}

	tex.Bind()

	return tex
}

func (t *Texture) SetWrapX(param WrapParams) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int32(param))
}

func (t *Texture) SetWrapY(param WrapParams) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int32(param))
}

func (t *Texture) SetMagFilter(param FilterParams) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(param))
}

func (t *Texture) SetMinFilter(param FilterParams) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(param))
}

func (t *Texture) Activate() {
	gl.ActiveTexture(t.textureAddress)
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}

func (t *Texture) Use() {
	t.Activate()
	t.Bind()
}

func (t *Texture) LoadImage(file string) {
	t.Lock()
	defer t.Unlock()

	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatalf("texture %q decoding error: %v", file, err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Fatalf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
}
