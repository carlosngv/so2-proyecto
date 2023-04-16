
import { useState } from 'react';
import './App.css';
import { BarChart } from './components/BarChart';
import { DiskChart } from './components/DiskChart';
import { MemoryChart } from './components/MemoryChart';
import { PieChart } from './components/PieChart';


function App() {


    return (
        <div className="container" id="App">

            <h1>Administrador de Tareas</h1>
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
