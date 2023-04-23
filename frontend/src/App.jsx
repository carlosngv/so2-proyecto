
import { useState } from 'react';
import './App.css';
import { BarChart } from './components/BarChart';
import { DiskChart } from './components/DiskChart';
import { MemoryChart } from './components/MemoryChart';
import { PieChart } from './components/PieChart';
import { BlockAllDevices, UnblockAllDevices } from '../wailsjs/go/main/App';



function App() {

    const handleClick = () => {
        BlockAllDevices();
      };

    const handleClickUnblock = () => {
        UnblockAllDevices();
      };


    return (
        <div className="container" id="App">

            <h1>Administrador de Tareas</h1>

            <div className="row">
                <div className="col-6">
                    <h2>Bloquear puertos</h2>
                        <button className="btn btn-primary" onClick={handleClick}>Block USB Ports</button>
                </div>

                <div className="col-6">
                    <h2>Desbloquear puertos</h2>
                        <button className="btn btn-primary" onClick={handleClickUnblock}>Unblock USB Ports</button>
                </div>
            </div>
            <div className="row">

                <div className="col-6">
                    <div className="cpu-container">
                        <span>CPU Usage</span>
                        <div style={{maxWidth: 400}} >
                            <PieChart />
                        </div>
                    </div>
                </div>

                <div className="col-6">
                <div className="cpu-container">
                        <span>Disk Usage</span>
                        <div style={{maxWidth: 400}} >
                            <DiskChart />
                        </div>
                    </div>
                </div>

            </div>

            <div className="row">
                <div className="col-6">
                    <div className="cpu-container">
                        <span>Memory Usage</span>
                        <div style={{maxWidth: 400}} >
                            <MemoryChart />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default App
