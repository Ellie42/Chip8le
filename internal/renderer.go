package internal

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"go.uber.org/zap"
	"runtime"
	"strings"
)

type Renderer struct {
	WindowWidth  int
	WindowHeight int

	ResolutionX int
	ResolutionY int

	OpenGL
}

type OpenGL struct {
	Window  *glfw.Window
	Program uint32
	VaoList []uint32
}

func NewRenderer() *Renderer {
	return &Renderer{
		WindowWidth:  800,
		WindowHeight: 400,
	}
}

func (r *Renderer) Init() {
	r.VaoList = make([]uint32, r.ResolutionX*r.ResolutionY)

	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		Logger.Panic("failed to initialise glfw", zap.Error(err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(r.WindowWidth, r.WindowHeight, "Chip8le - Spicy emulation", nil, nil)

	if err != nil {
		Logger.Panic("failed to open window", zap.Error(err))
	}

	r.Window = win

	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		Logger.Panic("failed to initialise opengl", zap.Error(err))
	}

	r.Program = gl.CreateProgram()

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)

	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	if err != nil {
		panic(err)
	}

	gl.AttachShader(r.Program, vertexShader)
	gl.AttachShader(r.Program, fragmentShader)
	gl.LinkProgram(r.Program)
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func (r *Renderer) RenderFrame(pixels *[]bool) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(r.Program)

	for y := 0; y < int(r.ResolutionY); y++ {
		for x := 0; x < int(r.ResolutionX); x++ {
			filled := (*pixels)[y*int(r.ResolutionX)+x]

			if filled {
				r.drawSquare(x, y)
			}
		}
	}

	r.Window.SwapBuffers()
	glfw.PollEvents()
}

var square = []float32{
	0, 1, 0,
	0, 0, 0,
	1, 0, 0,

	0, 1, 0,
	1, 0, 0,
	1, 1, 0,
}

func (r *Renderer) drawSquare(x int, y int) {
	index := x + y*r.ResolutionX
	vao := r.VaoList[index]

	y = (r.ResolutionY - 1) - y

	if vao == 0 {
		points := make([]float32, len(square))

		pixelWidth := 1.0 / float32(r.ResolutionX)
		pixelHeight := 1.0 / float32(r.ResolutionY)

		for i := 0; i < len(square); i++ {

			if i%3 == 0 {
				points[i] = (float32(x)*pixelWidth+square[i]*pixelWidth)*2 - 1
			} else if i%3 == 1 {
				points[i] = (float32(y)*pixelHeight+square[i]*pixelHeight)*2 - 1
			}
		}

		vao = makeVao(points)
		r.VaoList[index] = vao
	}

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

const (
	vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 410
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1, 1, 1, 1);
    }
` + "\x00"
)

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

func (r *Renderer) Stop() {
	//r.Window.Destroy()
	//glfw.PollEvents()
	//glfw.Terminate()
}
