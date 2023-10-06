package main

import (
	"fmt"
	"os/exec"
)

const (
	BlenderPath = "E:/blender/blender.exe"
	ProjectPath = "cube_diorama.blend"
	ScriptPath  = "render.py"
)

func main() {
	render()
}

func render() {
	var c *exec.Cmd
	c = exec.Command("cmd", "/C", BlenderPath, ProjectPath, "-P", ScriptPath, "-f", "0")

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}
