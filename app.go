package main

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/pbnjay/memory"
	"github.com/shirou/gopsutil/cpu"
	"github.com/wailsapp/wails"
)

// App struct

type CPUData struct {
	UserUsage float64
	SystemUsage float64
	IdleUsage float64
}
// type App struct {
// 	ctx context.Context
// }

// // NewApp creates a new App application struct

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

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
	mem := Memory{
		MemTotal: memTotal,
		MemFree: memFree,
		MemUsed: memUsed,
	}
	return mem
}

func (a *App) BlockUSBPorts() error {
	var cmd *exec.Cmd

	fmt.Printf("OS: %v", runtime.GOOS)

	if runtime.GOOS == "darwin" {
		cmd = exec.Command("sudo", "pmset", "-a", "hibernatemode", "0")
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("sudo", "modprobe", "-r", "usb_storage")
		// cmd := exec.Command("sh", "-c", "echo 'contraseña1234' | sudo -S chmod 0000 /media")

	} else {
		return fmt.Errorf("unsupported operating system")
	}

	err := cmd.Run()
	fmt.Println("Puertos USB Bloqueados!")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) UnblockUSBPorts() error {
	var cmd *exec.Cmd
	fmt.Printf("OS: %v", runtime.GOOS)
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("sudo", "kextload", "-b", "com.apple.driver.AppleUSBFTDI")
	case "linux":
		cmd = exec.Command("sudo", "modprobe", "-r", "usb-storage")
		// cmd := exec.Command("sh", "-c", "echo 'contraseña1234' | sudo -S chmod 0777 /media")
	default:
		fmt.Println("Platform not supported")
		return nil
	}

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	fmt.Println("Puertos USB desbloqueados!")
	return nil
}


// Bloquear todos los dispositivos USB.
func (a *App) BlockAllDevices()  {
    // Ejecutar el comando adecuado para bloquear todos los dispositivos USB en función del sistema operativo.
    var cmd *exec.Cmd

	fmt.Printf("Sistema Operativo: %v\n", runtime.GOOS)
    if runtime.GOOS == "darwin" {
        cmd = exec.Command("sudo", "sh", "-c", "echo 'disable' > /sys/bus/usb/drivers/usb/unbind")
    } else if runtime.GOOS == "linux" {
        cmd = exec.Command("sudo", "sh", "-c", "echo '0' > /sys/bus/usb/drivers/usb/unbind")
    } else {
        fmt.Errorf("Sistema operativo no soportado")
    }

    // Ejecutar el comando y comprobar si se ha producido algún error.
    err := cmd.Run()
    if err != nil {
		fmt.Println("ERROR")
        return
    }
	fmt.Println("Esperando 10 segundos...")
    time.Sleep(10 * time.Second)
}

// Desbloquear todos los dispositivos USB.
func (a *App) UnblockAllDevices() {
    // Ejecutar el comando adecuado para desbloquear todos los dispositivos USB en función del sistema operativo.
    var cmd *exec.Cmd

	fmt.Printf("Sistema Operativo: %v\n", runtime.GOOS)

    if runtime.GOOS == "darwin" {
        cmd = exec.Command("sudo", "sh", "-c", "echo 'enable' > /sys/bus/usb/drivers/usb/bind")
    } else if runtime.GOOS == "linux" {
        cmd = exec.Command("sudo", "sh", "-c", "echo '1' > /sys/bus/usb/drivers/usb/bind")
    } else {
 		fmt.Errorf("Sistema operativo no soportado")
    }

    // Ejecutar el comando y comprobar si se ha producido algún error.
    err := cmd.Run()
    if err != nil {
		fmt.Println("ERROR")
        return
    }

}
