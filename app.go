package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/pbnjay/memory"
	"github.com/shirou/gopsutil/cpu"
	"github.com/wailsapp/wails"
)

// App struct
var usbFileNames string
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
	a.deleteFile()
	a.deleteFolder()
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
	cmd.Stdin = os.Stdin

	fmt.Printf("OS: %v", runtime.GOOS)

	if runtime.GOOS == "darwin" {
		cmd = exec.Command("sudo", "pmset", "-a", "hibernatemode", "0")
	} else if runtime.GOOS == "linux" {
		//cmd = exec.Command("sudo", "modprobe", "-r", "usb_storage")
		cmd = exec.Command("sh", "-c", "echo 'Xunix..unix' | sudo -S chmod 0000 /media")

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
        cmd = exec.Command("sudo", "pmset", "-a", "hibernatemode", "0")
    } else if runtime.GOOS == "linux" {
        //cmd = exec.Command("sudo", "sh", "-c", "echo '0' > /sys/bus/usb/drivers/usb/unbind")
		cmd = exec.Command("sh", "-c", "echo 'Xunix..unix' | sudo -S chmod 0000 /media")
    } else {
        fmt.Errorf("Sistema operativo no soportado")
    }

    // Ejecutar el comando y comprobar si se ha producido algún error.
	a.WriteLog("Puertos USB bloqueados.")
    err := cmd.Run()
    if err != nil {
        return
    }
}

// Desbloquear todos los dispositivos USB.
func (a *App) UnblockAllDevices() {
    // Ejecutar el comando adecuado para desbloquear todos los dispositivos USB en función del sistema operativo.
    var cmd *exec.Cmd

	fmt.Printf("Sistema Operativo: %v\n", runtime.GOOS)

    if runtime.GOOS == "darwin" {
        cmd = exec.Command("sudo", "kextload", "-b", "com.apple.driver.AppleUSBFTDI")
    } else if runtime.GOOS == "linux" {
		cmd = exec.Command("sh", "-c", "echo 'Xunix..unix' | sudo -S chmod 0777 /media")

    } else {
 		fmt.Errorf("Sistema operativo no soportado")
    }

	a.WriteLog("Puertos USB desbloqueados.")
    // Ejecutar el comando y comprobar si se ha producido algún error
    err := cmd.Run()
    if err != nil {

        return
    }

}

func (a *App) readFilesInFolder(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(usbFileNames, f.Name()) == false {
			usbFileNames = fmt.Sprintf("%v|%v", f.Name(), usbFileNames)
			if(strings.Contains(path, "Volumes")) {
				a.WriteLog(fmt.Sprintf("Archivo %v ya se encuentra en el USB", f.Name()))
			} else {
				a.WriteLog(fmt.Sprintf("Archivo %v ya se encuentra en la carpeta local", f.Name()))
			}
		}
	}
	fmt.Printf("Archivos en USB y local: %v\n", usbFileNames)
}

func (a *App) ManageLogs() {
	mediaPath :="/Volumes/Seagate Expansion Drive/"
	localPath := "ArchivosUSB"
	pathExists, err := validateUSBPath(mediaPath)

	if err != nil {
		fmt.Printf("La ruta no existe, se ha generado el siguiente error: %v\n", err)
		time.Sleep(10 * time.Second)
		return
	}

	fmt.Printf("PathExists: %v\nPath: %v", pathExists, mediaPath)


	if pathExists {
		a.readFilesInFolder(mediaPath)
		a.readFilesInFolder(localPath)
	}
}



// Comprueba si la ruta especificada existe
func validateUSBPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func (a *App) deleteFile() {
	if _, err := os.Stat("bitacora.log"); err == nil {
		e := os.Remove("bitacora.log")
		if e != nil {
			log.Fatal(e)
		}
	 } else {
		fmt.Printf("File does not exist\n");
	 }
}

func (a *App) deleteFolder() {
	dir := "ArchivosUSB"
	if _, err := os.Stat(dir); os.IsNotExist(err) {

		err := os.Remove(dir)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Directory", dir, "removed successfully")
		}

		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	} else {
	   fmt.Println("The provided directory named", dir, "exists")
	}
}

func(a *App) WriteLog(text string) {
	file, e := os.OpenFile("bitacora.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if e != nil {
		log.Fatalln("Failed")
		return
	}

	log.SetOutput(file)
	dt := time.Now()
	fmt.Println(dt.Format(time.Kitchen))
	log.Println(dt.String() + " || " +text)


}
