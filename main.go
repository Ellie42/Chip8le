package main

import (
	"flag"
	"git.agehadev.com/elliebelly/chip8le/internal"
	"git.agehadev.com/elliebelly/chip8le/internal/engine"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}


	renderer := internal.NewRenderer()

	renderer.ResolutionX = 64
	renderer.ResolutionY = 32

	renderer.Init()
	defer renderer.Stop()

	e := engine.NewEngine(renderer)

	e.Init()

	program := &internal.Program{
		FilePath: "games/Breakout [Carmelo Cortez, 1979].ch8",
	}

	e.LoadProgram(program)

	e.Run()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}