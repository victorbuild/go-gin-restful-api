import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '10s', target: 50 },   // 10秒內逐步增加到50個用戶
        { duration: '30s', target: 100 },  // 維持在100個用戶
        { duration: '10s', target: 200 },  // 逐步增加到200個用戶
        { duration: '10s', target: 0 },    // 逐步降回0
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95%的請求在500ms內完成
        http_req_failed: ['rate<0.01'],   // 錯誤率低於1%
    },
};

export default function () {
    let payload = JSON.stringify({
        email: "example@example.com",
        password: "example@example.com123"
    });

    let params = {
        headers: {
            "Content-Type": "application/json"
        }
    };

    let res = http.post("http://localhost:8000/v1/auth/login", payload, params);

    check(res, {
        "status is 200": (r) => r.status === 200,
        "response time < 500ms": (r) => r.timings.duration < 500,
        "response has body": (r) => r.body !== null && r.body.length > 0,
    });

    if (res.status === 404) {
        console.log(`404 Not Found: ${res.url} - 請確認應用程式是否正在運行`);
    } else {
        console.log(`Response: ${res.status}, Time: ${res.timings.duration}ms, URL: ${res.url}, Method: ${res.request.method}`);
    }

    sleep(1);
}
