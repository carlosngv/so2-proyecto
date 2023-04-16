package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/pbnjay/memory"
	"github.com/shirou/gopsutil/cpu"
	"github.com/wailsapp/wails"
)

// App struct
type App struct {
	ctx context.Context
}

type CPUData struct {
	UserUsage float64
	SystemUsage float64
	IdleUsage float64
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, you're welcome!", name)
}

func (a *App) LogInfo(text string) string {
	return fmt.Sprintf("New log: %s", text)
}

func (a * App) WailsInit(runtime *wails.Runtime) error {
	go func() {
		runtime.Events.Emit("cpu_usage", a.GetCPUPercentage())
		time.Sleep(1 * time.Second)
	}()

	return nil
}



// func (a *App) GetCPUPercentage() *CPUData{
// 	perc, err := cpu.Percent(time.Second, true)
// 	if err != nil {

// 		return nil
// 	}

// 	data := CPUData{perc[cpu.CPUser], perc[cpu.CPSys], perc[cpu.CPIdle]}

// 	return &data
// }
func (a *App) GetCPUPercentage() CPUData{
	perc, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		return CPUData{}
	}

	return CPUData{perc[cpu.CPUser], 0, 0}
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// Creating structure for DiskStatus
type DiskStatus struct {
	All  float64
	Used float64
	Free float64
}

func (a *App) DiskUsage() DiskStatus {
	disk := DiskStatus{}
	fs := syscall.Statfs_t{}
	err := syscall.Statfs("/", &fs)
	if err != nil {
		return disk
	}
	all := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	used := disk.All - disk.Free
	disk.All = float64(all)/float64(GB)
	disk.Free = float64(used)/float64(GB)
	disk.Used = float64(free)/float64(GB)
	return disk
}

type Memory struct {
    MemTotal     uint64
    MemFree      uint64
    MemUsed uint64
}

func (a *App) ReadMemoryStats() Memory {
	memTotal := memory.TotalMemory()/1000000
	memFree := memory.FreeMemory()/100000
	memUsed := memTotal - memFree
    fmt.Printf("Total system memory: %d\n", memTotal)
    fmt.Printf("Free system memory: %d\n", memFree)
    fmt.Printf("Used system memory: %d\n", memUsed)
	mem := Memory{
		MemTotal: memTotal,
		MemFree: memFree,
		MemUsed: memUsed,
	}
	return mem
}
