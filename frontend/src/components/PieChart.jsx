import React, { useEffect, useState } from 'react'
import { Pie } from "react-chartjs-2";
import Chart from 'chart.js/auto';
import { GetCPUPercentage } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime';

export const PieChart = () => {

  const [userData, setUserData] = useState({
    // labels: UserData.map(data => data.year), // returns a new array of "year"
    labels: ["User", "System"], // returns a new array of "year"
    datasets: [{
        label: "CPU Usage",
        data: [50, 50],
        backgroundColor: ["blue", "cyan"]
    }]
});

  useEffect(() => {

    setInterval(() => {
      GetCPUPercentage().then(res => {
        console.log(res)
        updateData([res.UserUsage, 100 - res.UserUsage])
      })
    }, 2000)

  }, []);

  const updateData = (perc) => {
    setUserData({
      ...userData,
      datasets: [{
        label: "CPU Usage",
        data: perc,
        backgroundColor:["blue", "cyan"]
      }]
    })
  }


  return (
    <Pie data={userData} />
  )
}
