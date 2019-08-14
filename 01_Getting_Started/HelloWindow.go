package main
import (
	"log"
	"runtime"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
    runtime.LockOSThread()
	
	if err := glfw.Init(); err != nil {
        panic(err)
    }
	// glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 3) //OpenGL最大版本
    glfw.WindowHint(glfw.ContextVersionMinor, 3) //OpenGl最小版本
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) //明确核心模式
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True) //Mac使用
	window, err := glfw.CreateWindow(800, 600, "LearnOpenGL", nil, nil)
	log.Println("created window")
	
	if  window == nil || err!= nil {
		panic(err)
	}
    defer glfw.Terminate()
	window.MakeContextCurrent()//通知glfw将当前窗口上下文设置为线程主上下文

	
	if err := gl.Init(); err != nil {
        panic(err)
	}
	
	gl.Viewport(0, 0, 800, 600)
	window.SetFramebufferSizeCallback(framebuffer_size_callback)

	for !window.ShouldClose(){
		window.SwapBuffers()
		glfw.PollEvents()    
	}
	
}

func framebuffer_size_callback(window *glfw.Window, width int , height int){
	log.Printf("resize width:%d,height:%d",width,height)
	gl.Viewport(0, 0, int32(width), int32(height))
}