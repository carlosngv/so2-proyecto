import React, { useEffect, useState } from 'react'
import { Pie } from "react-chartjs-2";
import Chart from 'chart.js/auto';
import { DiskUsage } from '../../wailsjs/go/main/App';

export const DiskChart = () => {

  const [userData, setUserData] = useState({
    labels: [ "Free", "Used"],
    datasets: [{
        label: "Disk Usage",
        data: [50, 50],
        backgroundColor: ["cyan", "red"]
    }]
});

  useEffect(() => {

    setInterval(() => {
      DiskUsage().then(res => {
        console.log(res)
        updateData([res.Used, res.All - res.Used])
      })
    }, 2000)

  }, []);

  const updateData = (perc) => {
    setUserData({
      ...userData,
      datasets: [{
        label: "Disk Usage",
        data: perc,
        backgroundColor: ["cyan", "red"]
      }]
    })
  }


  return (
    <Pie data={userData} />
  )
}
