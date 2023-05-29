import { check } from 'k6';
import http from 'k6/http';

export function testLogin(params) {
  http.get('http://localhost:8081/hello');
}

export default function () {
  testLogin();
}
