import React, {useState} from "react";
import axios from "axios";
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
function PlaceOrder() {
    const [priceType, setPriceType] = useState('')
    const [price, setPrice] = useState('')
    const [quantity, setQuantity] = useState('')
    const placeOrder = async (orderType) => {
        try {
            const res = await axios.post(
                'http://localhost:8888/api/v1/trade/orders',
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
                    <label htmlFor="selection" className="col-sm-2 col-form-label">訂單類型</label>
                    <div className="col-sm-10">
                        <select
                            onChange={(e)=> setPriceType(e.target.value)}
                            className="form-select col-form-label"
                            aria-label="Default select example"
                            id="selection">
                            <option value="" key={1}>請選擇類型</option>
                            <option value="market" key={2}>市價單</option>
                            <option value="limit" key={3}>限價單</option>
                        </select>
                    </div>
                </div>
                <div className="row mb-3">
                    <label htmlFor="input3" className="col-sm-2 col-form-label">價格</label>
                    <div className="col-sm-10">
                        <input
                            onChange={(e)=> setPrice(e.target.value)}
                            type="text"
                            className="form-control"
                            id="input3"
                        />
                    </div>
                </div>
                <div className="row mb-3">
                    <label htmlFor="input13" className="col-sm-2 col-form-label">數量</label>
                    <div className="col-sm-10">
                        <input
                            onChange={(e)=> setQuantity(e.target.value)}
                            type="text"
                            className="form-control"
                            id="input13"/>
                    </div>
                </div>
                <div>
                    <button type={null} className="btn btn-primary m-2" onClick={()=> placeOrder('bid')}>買入</button>
                    <button type={null} className="btn btn-danger m-2" onClick={()=> placeOrder('ask')}>賣出</button>
                </div>
        </div>
    )
}

export default PlaceOrder;