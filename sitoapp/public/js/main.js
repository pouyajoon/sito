(function() {
  'use strict';
  /*global WebSocket*/

  var ws = new WebSocket('ws://' + window.location.host + '/ws');

  function sendMousePosition(e) {
    if (e === undefined) {
      e = {
        pageX: -1,
        pageY: -1
      };
    }
    var d, ds, p = document.getElementById('player');
    d = {
      'player': p !== null ? p.value : 'P.',
      'x': e.pageX,
      'y': e.pageY
    };
    ds = JSON.stringify(d);
    console.log('send d', ds);
    ws.send(ds);
  }


  ws.onopen = function() {

    console.log('ws open');
    document.getElementById('player').value = 'Guest ' + Math.random();
    sendMousePosition();

    // window.onmousemove = sendMousePosition;
  };

  ws.onclose = function(e) {
    console.log('ws close', e.data);
  };

  ws.onmessage = function(e) {
    console.log('ws message', JSON.parse(e.data));
    return e;
  };



  console.log('hello', window.location.host);
}());
