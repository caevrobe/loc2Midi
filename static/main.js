var readyToSend = false;

function main() {
   let socket = new WebSocket("ws://127.0.0.1:7777/socket");
   socket.onopen = function(e) {
      readyToSend = true;
   };

   /*socket.onmessage = function(event) {
      alert(`[message] Data received from server: ${event.data}`);
   };*/

   var elem = document.getElementById('graph');
   graph = Desmos.GraphingCalculator(elem);
   $.getJSON("https://www.desmos.com/calculator/qzcw2goadz").done(data => cont(graph, data, socket));
}

function cont(graph, data, socket) {
   graph.setState(data.state);

   // todo get columns x, y (array ind 0, 1) and load values from config.json
   // x: loaded_val\cdot c_{0}.x+c.x
   // y: loaded_val\cdot c_{0}.y+c.y
   //absoluteSoundCoords = graph.getExpressions()[16];

   var x_val, y_val;
   var listener_x = graph.HelperExpression({ latex: 'l_0'});
   listener_x.observe('numericValue', function() {
      x_val = listener_x.numericValue;
      triggerSend(socket, x_val, y_val);
   });
   var listener_y = graph.HelperExpression({ latex: 'l_1'});
   listener_y.observe('numericValue', function() {
      y_val = listener_y.numericValue;
      triggerSend(socket, x_val, y_val);
   });
}

function triggerSend(socket, x_val, y_val) {
   console.log(`${x_val}, ${y_val}`)
   if (readyToSend && x_val != null  && y_val != null) {
      socket.send(`${x_val}, ${y_val}`)
   }
}