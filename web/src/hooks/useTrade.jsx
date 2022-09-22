import React, {useState} from "react";
import useWebSocket from "react-use-websocket";

function useTrade() {
    const [socketUrl, setSocketUrl] = useState('ws://0.0.0.0:8080/ws');
    const [askDepth, setAskDepth] = useState([]);
    const [bidDepth, setBidDepth] = useState([]);
    const [latestPrice, setLatestPrice] = useState('')
    const [tradeLog, setTradeLog] = useState([])
    const [myTrade, setMyTrade] = useState([])
    const { lastMessage, readyState } = useWebSocket(socketUrl);

    return {
        socketUrl,
        askDepth,
        setAskDepth,
        bidDepth,
        setBidDepth,
        latestPrice,
        setLatestPrice,
        lastMessage,
        readyState,
        tradeLog,
        setTradeLog,
        myTrade,
        setMyTrade,
    }
}

export default useTrade;