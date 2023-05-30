import { check } from 'k6';
import http from 'k6/http';

export function test(params) {
  // http.get('http://localhost:8081/hello');
  const responses = http.batch([
    {
      method: 'POST',
      url: 'http://localhost:8081/sum',
      body: JSON.stringify([1, 2, 3]),
    },
    {
      method: 'POST',
      url: 'http://localhost:8081/sum',
      body: JSON.stringify([1, 2, 3, 5]),
    },
    {
      method: 'POST',
      url: 'http://localhost:8081/sum',
      body: JSON.stringify([0]),
    },
  ])
  check(responses[0], {
    'sum data OK': (res) => JSON.parse(res.body) == 6,
  });
  check(responses[1], {
    'sum data OK': (res) => JSON.parse(res.body) == 11,
  });
  check(responses[2], {
    'sum data OK': (res) => JSON.parse(res.body) == 0,
  });
}

export default function () {
  test();
}
