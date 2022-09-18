import React, {useState, useEffect} from "react";
import axios from "axios";
import {toast} from "react-toastify";

function formatTime(t) {
    let d = new Date(t);
    let minutes = d.getMinutes();
    let month = d.getMonth() + 1;
    let date = d.getDate() + 1;
    let second = d.getSeconds();
    let hour = d.getHours();
    if (minutes < 10) minutes = `0${minutes}`
    if (month < 10) month = `0${month}`
    if (date < 10) date = `0${date}`
    if (second < 10) second = `0${second}`
    if (hour < 10) hour = `0${hour}`
    return d.getFullYear() + '-' + month + '-' + date + ' ' + hour + ':' + minutes + ':' + second;
}

function TradeLog(props) {
    const {
        tradeLog,
        setTradeLog,
        lastMessage,
        setLatestPrice
    } = props

    useEffect(() => {
        if (!lastMessage) return
        let data = lastMessage.data.split('\n')
        data.forEach(message => {
            let msg = JSON.parse(message)
            if (msg.tag !== "trade") return
            setTradeLog((prev) => {
                return prev.length >= 10
                    ? [msg.data, ...prev.slice(0, 9)]
                    : [msg.data, ...prev]
            })
        })
    }, [lastMessage]);

    const fetch = async () => {
        try {
        const res = await axios.get('http://localhost:8080/api/v1/trade/log')
            setTradeLog(() => {
                let tradeLog = res.data.data.trade_log.length >= 10 ? res.data.data.trade_log.slice(0,10) : res.data.data.trade_log
                return tradeLog.sort((a, b) => {
                    if (a.trade_time > b.trade_time) return -1;
                    if (a.trade_time < b.trade_time) return 1;
                    return 0;
                })} )
            setLatestPrice(res.data.data.latest_price)
        } catch (e) {
            toast.error(e.response.data.message)
        }
    }
    useEffect(() => {
        fetch()
    }, [])
    return (
        <div>
            <hr/>
            <table className="table">
                <thead>
                <tr>
                    <th scope="col">價格</th>
                    <th scope="col">數量</th>
                    <th scope="col">金額</th>
                    <th scope="col">時間</th>
                </tr>
                </thead>
                <tbody>
                {
                    tradeLog.map((log, index) => (
                        <tr key={index}>
                            <td>{log.trade_price}</td>
                            <td>{log.trade_quantity}</td>
                            <td>{log.trade_amount}</td>
                            <td>{formatTime(log.trade_time/1e6)}</td>
                        </tr>
                    ))
                }

                </tbody>
            </table>
        </div>
    )
}

export default TradeLog;