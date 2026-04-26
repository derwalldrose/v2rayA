import assert from 'node:assert/strict';

function parseLegacyPlugin(pluginString) {
  const obj = {
    obfs: 'http',
    plugin: '',
    impl: '',
    tls: '',
    mode: 'websocket',
    path: '/',
    host: '',
  };
  const arr = pluginString.split(';');
  obj.plugin = arr[0];
  switch (obj.plugin) {
    case 'obfs-local':
    case 'simpleobfs':
      obj.plugin = 'simple-obfs';
      break;
    case 'v2ray-plugin':
      obj.tls = '';
      obj.mode = 'websocket';
      break;
  }
  for (let i = 1; i < arr.length; i++) {
    const a = arr[i].split('=');
    if (a.length > 2) {
      a[1] = a.slice(1).join('=');
      a.splice(2);
    }
    switch (a[0]) {
      case 'obfs':
        obj.obfs = a[1];
        break;
      case 'host':
      case 'obfs-host':
        obj.host = a[1];
        break;
      case 'path':
      case 'obfs-path':
        obj.path = a[1];
        break;
      case 'mode':
        obj.mode = a[1];
        break;
      case 'tls':
        obj.tls = 'tls';
        break;
      case 'impl':
        obj.impl = a[1];
        break;
    }
  }
  return obj;
}

const parsed = parseLegacyPlugin('v2ray-plugin;mode=websocket;host=ss.batch1.workers.dev;path=/?enc=aes-128-gcm&proxyip=kr.270376.xyz:50001;mux=0');
assert.equal(parsed.plugin, 'v2ray-plugin');
assert.equal(parsed.mode, 'websocket');
assert.equal(parsed.host, 'ss.batch1.workers.dev');
assert.equal(parsed.path, '/?enc=aes-128-gcm&proxyip=kr.270376.xyz:50001');

console.log('2.2.7.5 frontend path parse regression test passed');
