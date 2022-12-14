import React, {useEffect, useState} from "react";
import axios from "axios";
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import 'react-tippy/dist/tippy.css'
import {
    Tooltip,
} from 'react-tippy';

function PlaceOrder() {
    const [priceType, setPriceType] = useState('')
    const [price, setPrice] = useState('')
    const [quantity, setQuantity] = useState('')
    const placeOrder = async (orderType) => {
        try {
            const res = await axios.post(
                'http://0.0.0.0:8080/api/v1/trade/orders',
                {
                    order_type: orderType,
                    price_type : priceType,
                    price : priceType === 'limit' ? price : '0',
                    quantity : quantity,
                },
            )
            toast.success(
                `order_type: ${orderType} price_type: ${priceType} price: ${price} quantity: ${quantity}`,
                {
                    position: "top-right",
                    autoClose: 2000,
                    hideProgressBar: false,
                    closeOnClick: true,
                    pauseOnHover: true,
                    draggable: true,
                    progress: undefined,
                })
        } catch (e) {
            toast.error(`error: ${e.response.data.message}`)
        }
    }

    const createRandomOrder = async (orderType) => {
        try {
            const res = await axios.post(
                'http://0.0.0.0:8080/api/v1/trade/randomOrders',
                {
                    order_type: orderType,
                },
            )
            toast.success(
                `order_type: ${orderType} ask_len: ${res.data.data.ask_len} bid_len: ${res.data.data.bid_len}`,
                {
                    position: "top-right",
                    autoClose: 2000,
                    hideProgressBar: false,
                    closeOnClick: true,
                    pauseOnHover: true,
                    draggable: true,
                    progress: undefined,
                })
        } catch (e) {
            toast.error(`error: ${e.response.data.message}`)
        }
    }

    useEffect(() => {
        if (priceType !== 'market') return

    }, [priceType])

    return (
        <div>
            <ToastContainer
                position="top-left"
                autoClose={5000}
                hideProgressBar={false}
                newestOnTop={false}
                closeOnClick
                rtl={false}
                pauseOnFocusLoss
                draggable
                pauseOnHover
            />
            <div className="row mb-3 mt-3">
                    <label htmlFor="selection" className="col-sm-2 col-form-label">????????????</label>
                    <div className="col-sm-10">
                        <select
                            onChange={(e)=> setPriceType(e.target.value)}
                            className="form-select col-form-label"
                            aria-label="Default select example"
                            id="selection">
                            <option value="" key={1}>???????????????</option>
                            <option value="market" key={2}>?????????</option>
                            <option value="limit" key={3}>?????????</option>
                        </select>
                    </div>
                </div>
                <div className="row mb-3">
                    <label htmlFor="input3" className="col-sm-2 col-form-label">??????</label>
                    <div className="col-sm-10">
                        <Tooltip
                            disabled={priceType !== 'market'}
                            title="????????????????????????????????????????????????????????????"
                            position="top"
                        >
                            <input
                                disabled={priceType === 'market'}
                                value={priceType === 'market' ? '?????????' : price}
                                onChange={(e)=> setPrice(e.target.value)}
                                type="text"
                                className="form-control"
                                id="input3"
                            />
                        </Tooltip>
                    </div>
                </div>
                <div className="row mb-3">
                    <label htmlFor="input13" className="col-sm-2 col-form-label">??????</label>
                    <div className="col-sm-10">
                        <input
                            onChange={(e)=> setQuantity(e.target.value)}
                            type="text"
                            className="form-control"
                            id="input13"/>
                    </div>
                </div>
                <div>
                    <button className="btn btn-primary m-2" onClick={()=> placeOrder('bid')}>??????</button>
                    <button className="btn btn-danger m-2" onClick={()=> placeOrder('ask')}>??????</button>
                    <button className="btn btn-primary m-2" onClick={()=> createRandomOrder('bid')}>???????????? 10 ???</button>
                    <button className="btn btn-danger m-2" onClick={()=> createRandomOrder('ask')}>???????????? 10 ???</button>
                </div>
            <div>

            </div>
        </div>
    )
}

export default PlaceOrder;