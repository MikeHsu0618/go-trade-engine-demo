import React, {useEffect} from "react";
import axios from "axios";

function OrderHistory() {

    return (
        <table className="table">
            <thead>
            <tr>
                <th scope="col">類型</th>
                <th scope="col">價格</th>
                <th scope="col">數量/已成交</th>
                <th scope="col">金額</th>
                <th scope="col">時間</th>
                <th scope="col">操作</th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <th scope="row">1</th>
                <td>Mark</td>
                <td>Otto</td>
                <td>@mdo</td>
                <td>Mark</td>
                <td>Otto</td>
            </tr>
            <tr>
                <th scope="row">2</th>
                <td>Mark</td>
                <td>Otto</td>
                <td>@mdo</td>
                <td>Mark</td>
                <td>Otto</td>
            </tr>
            <tr>
                <th scope="row">3</th>
                <td>Mark</td>
                <td>Otto</td>
                <td>@mdo</td>
                <td>Mark</td>
                <td>Otto</td>
            </tr>
            </tbody>
        </table>
    )
}

export default OrderHistory;