import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 50 },
    { duration: '10s', target: 100 },
    { duration: '10s', target: 200 },
    { duration: '20s', target: 400 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(95)<2000'], // %95 cevap süresi 2sn altında olmalı
  },
};

export default function () {
  http.get('http://localhost:3000/api/hello');
  sleep(0.05);
}
