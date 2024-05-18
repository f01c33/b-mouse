package main

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.6-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	return prog
}

func drawLine(x1, y1, x2, y2 int32) {
	gl.Color3f(1.0, 0, 0)
	gl.LineWidth(2.0)
	gl.Begin(gl.LINES)
	gl.Vertex2i(x1, y1)
	gl.Vertex2i(x2, y2)
	gl.End()
}

var RX, UY, LX, DY, w, h int

func reset() {
	RX, UY = 0, 0
	LX, DY = w, h
	robotgo.Move((RX+LX)/2, (UY+DY)/2)
}

func main() {

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	pm := glfw.GetPrimaryMonitor()
	vm := pm.GetVideoMode()
	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)
	glfw.WindowHint(glfw.Floating, glfw.True)
	glfw.WindowHint(glfw.Decorated, glfw.False)

	// if it's the same size of the window it doesn't work (on windows), idk why
	window, err := glfw.CreateWindow(vm.Width-1, vm.Height-1, "much wow", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	prog := gl.CreateProgram()
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	w, h = window.GetSize()

	gl.Ortho(0, // left
		float64(w), // right
		float64(h), // bottom
		0,          // top
		0,          // zNear
		1,          // zFar
	)

	gl.UseProgram(prog)

	mX, mY := w/2, h/2
	reset()
	robotgo.Move(mX, mY)
	hidden := false
	EvChan := hook.Start()
	defer hook.End()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	drawLine(int32(mX), int32(UY), int32(mX), int32(DY))
	drawLine(int32(LX), int32(mY), int32(RX), int32(mY))
	window.SwapBuffers()
	lastX, lastY := mX, mY
	lastClick := 0
	for !window.ShouldClose() {
		glfw.PollEvents()
		for len(EvChan) > 1 {
			<-EvChan
		}
		select {
		case ev := <-EvChan:
			if ev.Kind == hook.MouseMove {
				lastX = int(ev.X)
				lastY = int(ev.Y)
			}
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			drawLine(int32(lastX), int32(UY), int32(lastX), int32(DY))
			drawLine(int32(LX), int32(lastY), int32(RX), int32(lastY))
			drawLine(int32(LX), int32(DY), int32(RX), int32(DY))
			drawLine(int32(LX), int32(UY), int32(RX), int32(UY))
			drawLine(int32(LX), int32(DY), int32(LX), int32(UY))
			drawLine(int32(RX), int32(DY), int32(RX), int32(UY))
			window.SwapBuffers()
			if ev.Rawcode == 65299 { // pause
				if !hidden {
					window.Hide()
				} else {
					window.Show()
				}
				hidden = !hidden
			}
			if ev.Kind == hook.KeyDown && !hidden {
				if ev.Rawcode == 'w' || ev.Rawcode == 'k' || ev.Rawcode == 65362 /*🠕*/ {
					mY = (UY + DY) / 2
					DY = mY
				} else if ev.Rawcode == 's' || ev.Rawcode == 'j' || ev.Rawcode == 65364 /*🠗*/ {
					mY = (UY + DY) / 2
					UY = mY
				} else if ev.Rawcode == 'a' || ev.Rawcode == 'h' || ev.Rawcode == 65361 /*🠔*/ {
					mX = (LX + RX) / 2
					LX = mX
				} else if ev.Rawcode == 'd' || ev.Rawcode == 'l' || ev.Rawcode == 65363 /*🠖*/ {
					mX = (LX + RX) / 2
					RX = mX
				} else if ev.Rawcode == 65307 || ev.Rawcode == 'q' { //␛
					reset()
				} else if ev.Rawcode == 65293 || ev.Rawcode == '1' { //enter
					window.Hide()
					time.Sleep(time.Millisecond * 100)
					robotgo.Click()
					window.Show()
					reset()
				} else if ev.Rawcode == 65506 || ev.Rawcode == '2' { //shift
					window.Hide()
					time.Sleep(time.Millisecond * 100)
					robotgo.Click("center")
					window.Show()
					reset()
				} else if ev.Rawcode == ' ' || ev.Rawcode == '3' { //␠
					window.Hide()
					time.Sleep(time.Millisecond * 100)
					robotgo.Click("right")
					window.Show()
					reset()
				}
				robotgo.Move((RX+LX)/2, (UY+DY)/2)
			} else if ev.Kind == hook.MouseDown && !hidden {
				lastClick = int(ev.Button)
			} else if ev.Kind == hook.MouseUp && !hidden {
				if lastClick != 0 {
					window.Hide()
					time.Sleep(time.Millisecond * 100)
					if ev.Button == 1 {
						robotgo.Click()
					} else if ev.Button == 3 {
						robotgo.Click("right")
					} else if ev.Button == 2 {
						robotgo.Click("center")
					}
					window.Show()
					reset()
					lastClick = 0
				}
			}
		}
	}
}
