import React, {useEffect, useRef, useState} from "react";
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import axios from "axios";
import {toast} from "react-toastify";

ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
);

function getTime() {
    let d = new Date();
    let minutes = d.getMinutes();
    let second = d.getSeconds();
    let hour = d.getHours();
    if (minutes < 10) minutes = `0${minutes}`
    if (second < 10) second = `0${second}`
    if (hour < 10) hour = `0${hour}`
    return hour + ':' + minutes + ':' + second;
}

function getLastItems(arr, amount) {
    if (arr.length >= amount) {
        return arr.slice(-amount)
    }
    return arr
}



function PriceChart(props) {
    const timer = useRef(null);
    const [data, setData] = useState([])
    const [label,setLabel] = useState([])
    const {latestPrice} = props
    const chartData = {
        labels: label,
        datasets: [
            {
                label: 'USD Spread',
                fill: false,
                lineTension: 0.1,
                backgroundColor: 'rgba(75,192,192,0.4)',
                borderColor: 'rgba(75,192,192,1)',
                borderCapStyle: 'butt',
                borderDash: [],
                borderDashOffset: 0.0,
                borderJoinStyle: 'miter',
                pointBorderColor: 'rgba(75,192,192,1)',
                pointBackgroundColor: '#fff',
                pointBorderWidth: 5,
                pointHoverRadius: 11,
                pointHoverBackgroundColor: 'rgba(75,192,192,1)',
                pointHoverBorderColor: 'rgba(220,220,220,1)',
                pointHoverBorderWidth: 2,
                pointRadius: 1,
                pointHitRadius: 10,
                data: data
            }
        ]
    };

     const pushChartData = async () => {
         try {
             const res = await axios.get('http://localhost:8080/api/v1/trade/log')
             setData((prev)=> getLastItems([...prev, res.data.data.latest_price], 10))
             setLabel((prev) => getLastItems([...prev, getTime()], 10))
         } catch (e) {
             toast.error(e.response.data.message)
         }
    }

    useEffect(() => {
        timer.current = setInterval(pushChartData, 2000)
        return () => clearInterval(timer.current);
    }, [])

    return (
        <div><Line data={chartData} /></div>
);
}

export default PriceChart;