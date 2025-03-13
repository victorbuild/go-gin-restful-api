import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 50, // 50 個並發用戶
    duration: '30s', // 測試 30 秒
};

export default function () {
    let payload = JSON.stringify({
        email: "***REMOVED***",
        password: "***REMOVED***123"
    });

    let params = {
        headers: {
            "Content-Type": "application/json"
        }
    };

    let res = http.post("http://localhost:8000/v1/auth/login", payload, params);

    check(res, {
        "status is 200": (r) => r.status === 200,
        "status is 401": (r) => r.status === 401,
        "status is 500": (r) => r.status === 500,
    });

    console.log(`Response: ${res.status}, Time: ${res.timings.duration}ms`);

    sleep(1);
}
