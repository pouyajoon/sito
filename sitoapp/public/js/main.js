(function() {
  'use strict';
  /*global WebSocket*/

  var ws = new WebSocket('ws://localhost:8081', "ProtocolOne");
  ws.onopen = function() {
    console.log('ws open');
  };

  console.log('hello');
}());
