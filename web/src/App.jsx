import React from 'react'
import OrderBook from './components/OrderBook/OrderBook.jsx'
import "./index.css"
import TradeLog from "./components/TradeLog/TradeLog.jsx";
import PlaceOrder from "./components/PlaceOrder/PlaceOrder.jsx";
import useTrade from "./hooks/useTrade.jsx";
import PriceChart from "./PriceChart/PriceChart.jsx";
function App() {
  const {
    askDepth,
    setAskDepth,
    bidDepth,
    setBidDepth,
    latestPrice,
    setLatestPrice,
    lastMessage,
    setTradeLog,
    tradeLog,
    setMyTrade,
    myTrade,
  } = useTrade()

  const props = {
    askDepth,
    setAskDepth,
    bidDepth,
    setBidDepth,
    latestPrice,
    setLatestPrice,
    lastMessage,
    tradeLog,
    setTradeLog,
    myTrade,
    setMyTrade,
  }

  return <div className="container text-center" style={{maxWidth: '1440px'}}>
    <div className="row align-items-start">
      <div className="col">
        <OrderBook {...props}/>
      </div>
      <div className=" col">
        <PriceChart {...props}/>
        <PlaceOrder />
      </div>
      <div className=" col">
        <TradeLog {...props}/>
      </div>
    </div>
    {/*<div className="row align-items-center">*/}
    {/*  <OrderHistory {...props} className="col-6"/>*/}
    {/*</div>*/}
  </div>

}

export default App;
