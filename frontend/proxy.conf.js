const target = 'http://127.0.0.1:8080';

const PROXY_CONFIG = [
  {
    context: [
      '/api',
    ],
    target: target,
    secure: false,
  },
];

module.exports = PROXY_CONFIG;
