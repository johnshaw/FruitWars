$(document).ready( function() {
  var board = new Board(0, 0, 960, 540, "#123456");
  var socket = new SocketAdapter(board, "ws://192.168.2.13:8080/screen");
  board.runLoop();
});

var Tile = {};

Tile = function(x, y, width, height, color) {
  this.x = x;
  this.y = y;
  this.width = width;
  this.height = height;
  this.color = color; 
};

Tile.prototype = {

  constructor: Tile

};

var Board = {};

Board = function(x, y, width, height, backgroundColor) {
  this.canvas = $("<canvas id='board' height=" + height +" width=" + width + "></canvas>");
  this.canvas.appendTo('body');
  this.ctx = document.getElementById('board').getContext('2d'); 
  this.x = x;
  this.y = y;
  this.width = width;
  this.height = height;
  this.tilesize = 30;
  this.backgroundColor = backgroundColor;

  this.tiles = [];
};

Board.prototype = {
  
  constructor: Board,
 
  //refresh, probably will remove
  refresh: function () {
    this.drawSelf(this.ctx, this.backgroundColor);
    this.drawTiles(this.ctx, this.tiles, this.tilesize);
  },

  runLoop: function() {
    this.drawSelf(this.ctx, this.backgroundColor); 
    this.drawTiles(this.ctx, this.tiles, this.tilesize);
    window.requestAnimationFrame(this.runLoop.bind(this));
  },

  drawTiles: function(ctx, tiles, tilesize) {
    tiles.forEach( function (tile) {
      this.drawTile(ctx, tile, tilesize);
    }, this);
  },

  drawTile: function(ctx, tile, tilesize) {
    ctx.fillStyle = tile.color;
    ctx.fillRect( tilesize * tile.x, tilesize * tile.y,
                  tilesize * tile.width, tilesize * tile.height ); 
  },

  drawSelf: function(ctx) {
    ctx.fillStyle = this.backgroundColor; 
    ctx.fillRect(this.x, this.y, this.width, this.height);
  },
 
};


//SOCKET ADAPTER OBJECT 
var SocketAdapter = {};

SocketAdapter = function(board, url) {

  this.board = board;
  
  this.ws = new WebSocket(url);

  this.ws.onopen = function(event) {
    console.log("connection established");
  }.bind(this);

  this.ws.onmessage = function(event) {
    this.parse(JSON.parse(event.data)); 
  }.bind(this);

};

SocketAdapter.prototype = {

  constructor: SocketAdapter,
  
  parse: function (data) {
    var baseData = data.Bases
      , dudeData = data.Dudes;    

    //here we wipe the tiles
    this.board.tiles = [];
    for (var i in baseData) {
      var base = baseData[i];
      this.board.tiles.push(new Tile(base.Pos.X, base.Pos.Y, 4, 4, "#00ff00"));
    }
    for (var i in dudeData) {
      var dude = dudeData[i];
      this.board.tiles.push(new Tile(dude.Pos.X, dude.Pos.Y, 1, 1, "#0000ff"));
    }
    this.board.refresh();
  },
};

//Helper methods
calculatePixelSize = function(height, pdenY) {
  return Math.floor(height/pdenY);
};

updateCanvasSize = function (pdenX, pdenY, pixelSize) {
   var canvasHeight = pdenY * pixelSize
     , canvasWidth = pdenX * pixelSize;

   return {height: canvasHeight, width: canvasWidth};  
};
