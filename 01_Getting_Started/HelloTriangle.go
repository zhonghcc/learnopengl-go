package main
import (
	"log"
	// "fmt"
	"strings"
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

	// Hello Triangle

	var vao uint32
	gl.GenVertexArrays(1,&vao)
	gl.BindVertexArray(vao)

	vertices:=[]float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0,  0.5, 0.0 }
		
	var vbo uint32
	gl.GenBuffers(1,&vbo)//创建顶点缓冲对象，绑定id
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)//把新创建的缓冲绑定到GL_ARRAY_BUFFER目标
	gl.BufferData(gl.ARRAY_BUFFER,len(vertices)*4,gl.Ptr(vertices),gl.STATIC_DRAW)//把用户定义的数据复制到当前绑定缓冲

	//连接顶点属性
	gl.VertexAttribPointer(0,3,gl.FLOAT,false,3*4,gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)


	//顶点着色器
	vertexShaderSource:=`
		#version 330 core
		layout (location = 0) in vec3 aPos;

		void main()
		{
			gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
		}
	`

	var vertexShader uint32
	vertexShader = gl.CreateShader(gl.VERTEX_SHADER)//创建一个Shader对象
	vertexSourcePointer, freeFunc := gl.Strs(nullTerminatedString(vertexShaderSource))
	gl.ShaderSource(vertexShader,1,vertexSourcePointer,nil)
	gl.CompileShader(vertexShader)
	freeFunc()

	var result int32 = gl.FALSE
	var infoLogLength int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &result)
	gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		vertexShaderErrorMessage := strings.Repeat("\x00", int(infoLogLength))
		
		var messageLength int32
		gl.GetShaderInfoLog(vertexShader, infoLogLength, &messageLength, gl.Str(vertexShaderErrorMessage))
		log.Printf("Shader Compile Error: %s (compile status: %d)\n", vertexShaderErrorMessage, result)
	}


	//片段着色器
	fragmentShaderSource:=`
		#version 330 core
		out vec4 FragColor;

		void main()
		{
			FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
		} 
	`

	var fragmentShader uint32
	fragmentShader = gl.CreateShader(gl.FRAGMENT_SHADER)//创建一个Shader对象
	fragmentSourcePointer, freeFunc := gl.Strs(nullTerminatedString(fragmentShaderSource))
	gl.ShaderSource(fragmentShader,1,fragmentSourcePointer,nil)
	gl.CompileShader(fragmentShader)
	freeFunc()

	result = gl.FALSE
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &result)
	gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		fragmentShaderErrorMessage := strings.Repeat("\x00", int(infoLogLength))
		
		var messageLength int32
		gl.GetShaderInfoLog(fragmentShader, infoLogLength, &messageLength, gl.Str(fragmentShaderErrorMessage))
		log.Printf("Shader Compile Error: %s (compile status: %d)\n", fragmentShaderErrorMessage, result)
	}

	//着色器程序
	var shaderProgram uint32
	shaderProgram = gl.CreateProgram()
	gl.AttachShader(shaderProgram,vertexShader)
	gl.AttachShader(shaderProgram,fragmentShader)
	gl.LinkProgram(shaderProgram)

	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &result)
	gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		programErrorMessage := strings.Repeat("\x00", int(infoLogLength))
		
		var messageLength int32
		gl.GetProgramInfoLog(shaderProgram, infoLogLength, &messageLength, gl.Str(programErrorMessage))
		log.Printf("Program Link Error: %s (compile status: %d)\n", programErrorMessage, result)
	}

	gl.UseProgram(shaderProgram)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)



	//渲染循环
	for !window.ShouldClose(){
		//用户输入
		processInput(window)

		//渲染
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)//状态设置
		gl.Clear(gl.COLOR_BUFFER_BIT)//状态使用

		gl.UseProgram(shaderProgram);
		gl.BindVertexArray(vao);
		gl.DrawArrays(gl.TRIANGLES, 0, 3);

		//检查调用事件，交换缓冲
		glfw.PollEvents()    
		window.SwapBuffers()
	}
	
}

func framebuffer_size_callback(window *glfw.Window, width int , height int){
	log.Printf("resize width:%d,height:%d",width,height)
	gl.Viewport(0, 0, int32(width), int32(height))
}

func processInput(window *glfw.Window){
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		log.Println("escape pressed")
		window.SetShouldClose(true)
	}
}
func nullTerminatedString(str string)string {
	return str+ "\x00"
}