import React, {useCallback, useEffect, useState} from "react";

// Suppose we are provided with the orderbook data like below...
const CRYPTO = 'BTC';
const FIAT = 'USD';


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
            if (msg.tag === "latest_price") {
                setLatestPrice(msg.data.latest_price)
            }
        })
    }, [lastMessage]);

  return (
    <div className='orderbook' style={{width: "300px"}}>
      <OrderBookHeader />
      <div className='orderBookContainer'>
        <OrderBookLabel fiat={FIAT} />
        <OrderBookTable
          sellEntries={askDepth}
          buyEntries={bidDepth}
          latestPrice={latestPrice}
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

function OrderBookTable({ sellEntries, buyEntries, latestPrice }) {
    return (
    <div className='tableContainer'>
      <OrderBookTableEntries entries={sellEntries} side={'sell'} />
      <OrderBookTableLabel spread={latestPrice} />
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

function OrderBookTableLabel({ spread }) {
  return (
    <div className='labelContainer'>
      <div className='leftLabelCol'>
        <span>{FIAT} Spread</span>
      </div>
      <div className='centerLabelCol'>
        <span>{spread}</span>
      </div>
    </div>
  );
}

export default OrderBook;