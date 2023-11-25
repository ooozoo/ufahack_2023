const target = 'http://192.168.102.97:8080';

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
