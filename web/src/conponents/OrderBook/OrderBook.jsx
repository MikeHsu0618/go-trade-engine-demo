import React, {useCallback, useEffect, useRef, useState} from "react";
// Suppose we are provided with the orderbook data like below...
const CRYPTO = 'BTC';
const FIAT = 'USD';
const spreadColor = {
    up: 'rgba(0, 255, 0, 0.8)',
    down: 'rgba(255, 0, 0, 0.8)',
    fair: 'rgb(138, 147, 159)'
}

function OrderBook(props) {
    const {
        askDepth,
        setAskDepth,
        bidDepth,
        setBidDepth,
        latestPrice,
        setLatestPrice,
        lastMessage,
    } = props
    const [trend, setTrend] = useState(spreadColor['fair'])
    useEffect(() => {
        if (!lastMessage) return
        let data = lastMessage.data.split('\n')
        data.forEach(message => {
            let msg = JSON.parse(message)
            if (msg.tag === "depth") {
                setAskDepth(() => {
                    let ask =msg.data.ask.map(item => ({
                        size: item[1],
                        price: item[0]
                    }))
                    return ask.sort((a, b) => {
                        let aPrice = Number.parseFloat(a.price)
                        let bPrice = Number.parseFloat(b.price)
                        if (aPrice > bPrice) return -1;
                        if (aPrice < bPrice) return 1;
                        return 0;
                    })
                })
                setBidDepth(msg.data.bid.map(item => ({
                    size: item[1],
                    price: item[0]
                })))
            }

            // if (msg.tag === 'trade') setTrend(spreadColor['fair'])

            if (msg.tag === "latest_price") {
                setLatestPrice(msg.data.latest_price)
                if (msg.data.latest_price === latestPrice) setTrend(spreadColor['fair'])
                if (msg.data.latest_price > latestPrice) setTrend(spreadColor['up'])
                if (msg.data.latest_price < latestPrice) setTrend(spreadColor['down'])
            }
        })
    }, [lastMessage]);

  return (
    <div className='orderbook' >
      <OrderBookHeader />
      <div className='orderBookContainer'>
        <OrderBookLabel fiat={FIAT} />
        <OrderBookTable
          sellEntries={askDepth}
          buyEntries={bidDepth}
          latestPrice={latestPrice}
          trend={trend}
          fiat={FIAT}
        />
      </div>
    </div>
  );
}

function OrderBookHeader() {
  return <div className='header'>Order Book</div>;
}

function OrderBookLabel() {
  return (
    <div className='labelContainer'>
      <div className='leftLabelCol'>
        <span>Market Size</span>
      </div>
      <div className='centerLabelCol'>
        <span>Price ({FIAT})</span>
      </div>
    </div>
  );
}

function OrderBookTable({ sellEntries, buyEntries, latestPrice, trend }) {
    return (
    <div className='tableContainer'>
      <OrderBookTableEntries entries={sellEntries} side={'sell'} />
      <OrderBookTableLabel spread={latestPrice} trend={trend}/>
      <OrderBookTableEntries entries={buyEntries} side={'buy'} />
    </div>
  );
}

function OrderBookTableEntries({ side, entries }) {
  return (
    <div className='entriesContainer'>
      {entries.map((order, i) => (
        <OrderBookTableEntry
          side={side}
          size={order.size}
          price={order.price}
          pos={order.pos}
          key={i}
         />
      ))}
    </div>
  );
}

function OrderBookTableEntry({ side, size, price, pos }) {
  return (
    <div className='entryContainer'>
      <div className='leftLabelCol'>
        <span className='entryText'>
          {size ? size : "No Data"}
        </span>
      </div>
      <div className='centerLabelCol'>
        <span className={`${side === 'sell' ? 'entrySellText' : 'entryBuyText'}`}>
          {price ? price : "No Data"}
        </span>
      </div>
    </div>
  );
}

function OrderBookTableLabel({ spread, trend }) {
  return (
    <div className='labelContainer'>
      <div className='leftLabelCol'>
        <span>{FIAT} Spread</span>
      </div>
      <div className='centerLabelCol'>
        <span style={{
            fontSize: 20,
            color: trend
        }}><span style={{
            marginRight: -3,
        }}>{spread}</span>
            {
                trend === spreadColor['up'] ? <i className="bi bi-arrow-up"></i> :
                    trend === spreadColor['down'] ? <i className="bi bi-arrow-down"></i> :
                    ''
            }
        </span>
      </div>
    </div>
  );
}

export default OrderBook;