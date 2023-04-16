import React, { useEffect, useState } from 'react'
import { Pie } from "react-chartjs-2";
import { ReadMemoryStats } from '../../wailsjs/go/main/App';

export const MemoryChart = () => {

  const [userData, setUserData] = useState({
    labels: [ "Free", "Used"],
    datasets: [{
        label: "<e,pru Usage",
        data: [50, 50],
        backgroundColor: ["cyan", "red"]
    }]
});

  useEffect(() => {

    setInterval(() => {
        ReadMemoryStats().then(res => {
        console.log(res)
        updateData([res.MemFree, res.MemUsed])
      })
    }, 2000)

  }, []);

  const updateData = (perc) => {
    setUserData({
      ...userData,
      datasets: [{
        label: "Memory Usage",
        data: perc,
        backgroundColor: ["cyan", "red"]
      }]
    })
  }


  return (
    <Pie data={userData} />
  )
}
