import http from "k6/http";
import { check } from "k6";

export const options = {
  vus: 20,
  duration: "30s",
  thresholds: {
    http_req_failed: ["rate<0.001"],
    http_req_duration: ["p(95)<200"],
  },
};

const baseURL = __ENV.BASE_URL || "http://localhost:8080";
const roomID = __ENV.ROOM_ID || "replace-room-id";
const token = __ENV.TOKEN || "replace-token";
const date = new Date().toISOString().slice(0, 10);

export default function () {
  const res = http.get(`${baseURL}/rooms/${roomID}/slots/list?date=${date}`, {
    headers: { Authorization: `Bearer ${token}` },
  });

  check(res, {
    "status is 200": (r) => r.status === 200,
  });
}
