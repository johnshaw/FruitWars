var express = require('express')
  , socket = require('engine.io');

var server = express();

server.configure(function() {
  server.use(express.static(__dirname));
  server.use(express.bodyParser());
});

server.get('/', function(req, res) {
  return res.sendfile(__dirname + "/index.html");
});

server.listen(8000);
