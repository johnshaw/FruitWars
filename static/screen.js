minispade.register('application.js', function() {$(document).ready( function() {
  var board = new Board(0, 0, 1920, 1080, "#123456");
  var socket = new SocketAdapter(board, "ws://192.168.2.13:8080/screen");
  board.runLoop();
});

var Tile = {};

Tile = function(x, y, width, height, color, imagemap) {
  this.x = x;
  this.y = y;
  this.width = width;
  this.height = height;
  this.color = color; 
  this.imagemap = imagemap;
};

Tile.prototype = {
  constructor: Tile
};

//Player Model
var Player = {};

Player = function () {};

Player.prototype = {
  constructor: Player,  
};

//ImageMap model
var ImageMap = {};

ImageMap = function (spritesheet, x, y, width, height) {
  this.spritesheet = spritesheet;
  this.x = x;
  this.y = y;
  this.width = width;
  this.height = height;
};

ImageMap.prototype = {
  constructor: ImageMap,
};

var ImageController = {};

ImageController = function(board) {
  this.board = board;
};

ImageController.prototype = {
  init: function() {
    this.spriteStore = [];
    this.spritesAreLoaded = 0;
    this.spriteSheetCount = 2;

    this.spritesheet = document.getElementById('spritesheet');
    this.maptiles = document.getElementById('maptiles');
    this.retrieveSpriteJSON("screen/json/spritesheet.json", this.spritesheet);
    this.retrieveSpriteJSON("screen/json/maptiles.json", this.maptiles);
  },

  getImageMap: function(imageName) {
    if (this.spriteStore[imageName]) {
      return this.spriteStore[imageName];
    } else {
      return null;
    }    
  },

  retrieveSpriteJSON: function(url, spritesheet) {
    $.ajax({
      url: url,
      dataType: "json",
      contentType: 'application/json; charset=utf-8',
      context: this,
      success: function (data) {
        for (var key in data.frames) {
          var spriteData = data.frames[key].frame
            , newImageMap = new ImageMap( spritesheet,
                                          spriteData.x, spriteData.y,
                                          spriteData.w, spriteData.h );  
          this.spriteStore[key] = newImageMap;
        }
        this.spritesAreLoaded++;
      }
    });
  },  
};

//Board Model
var Board = {};

Board = function(x, y, width, height, backgroundColor) {
  //drawing is done to each layer based on "type"
  this.background = $("<canvas id='bgcan' height=" + height +" width=" + width + "></canvas>");
  this.maptiles = $("<canvas id='mapcan' height=" + height +" width=" + width + "></canvas>");
  this.entitylayer = $("<canvas id='entcan' height=" + height +" width=" + width + "></canvas>");
  this.foreground = $("<canvas id='fgcan' height=" + height +" width=" + width + "></canvas>");

  this.background.appendTo('body');
  this.maptiles.appendTo('body');
  this.entitylayer.appendTo('body');
  this.foreground.appendTo('body');

  this.bgctx = document.getElementById('bgcan').getContext('2d'); 
  this.mapctx = document.getElementById('mapcan').getContext('2d'); 
  this.entityctx = document.getElementById('entcan').getContext('2d'); 
  this.fgctx = document.getElementById('fgcan').getContext('2d'); 

  this.x = x;
  this.y = y;
  this.width = width;
  this.height = height;
  this.tilesize = this.height / 18;

  this.backgroundImage = document.getElementById('map');

  this.backgroundColor = backgroundColor;
  this.foregroundColor = backgroundColor;
  
  this.imageController = new ImageController(this);
  this.imageController.init(); 
 
  this.mapTiles = [];
  this.entities = [];
};

Board.prototype = {
  
  constructor: Board,
 
  runLoop: function() {
    this.drawImageLayer(this.bgctx, this.backgroundImage);
    if (this.imageController.spritesAreLoaded === 
        this.imageController.spriteSheetCount) {
      //this.drawTiles(this.mapctx, this.mapTiles, this.tilesize);
      this.clearLayer(this.entityctx);
      this.drawTiles(this.entityctx, this.entities, this.tilesize);
      //this.drawTiles(this.fgctx, this.foregroundColor);
    }
    window.requestAnimationFrame(this.runLoop.bind(this));
  },

  drawImageLayer: function(ctx, image) {
    console.log(image);
    ctx.drawImage(image, 0, 0);
  },

  drawLayer: function(ctx, color) {
    ctx.fillStyle = color;
    ctx.fillRect(this.x, this.y, this.width, this.height); 
  },

  clearLayer: function(ctx) {
    ctx.clearRect(this.x, this.y, this.width, this.height);
  },

  drawTiles: function(ctx, tiles, tilesize) {
    tiles.forEach( function (tile) {
      this.drawTile(ctx, tile, tilesize);
    }, this);
  },

  drawTile: function(ctx, tile, tilesize) {
    ctx.drawImage(tile.imagemap.spritesheet,
                  tile.imagemap.x, tile.imagemap.y, 
                  tile.imagemap.width, tile.imagemap.height,
                  tilesize * tile.x, tilesize * tile.y,
                  tilesize * tile.width, tilesize * tile.height )
  },

  drawColorTile: function(ctx, tile, tilesize) {
    ctx.fillStyle = tile.color;
    ctx.fillRect( tilesize * tile.x, tilesize * tile.y,
                  tilesize * tile.width, tilesize * tile.height ); 
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
      , dudeData = data.Dudes
      , mapData = data.Map;

    //here we wipe the tiles
    this.board.entities = [];
    this.board.mapTiles = [];

    var dudeMapping = {};
    var baseMapping = {};

    var fruits = ['apple', 'banana', 'watermelon', 'grape'];
    var types = ['antidude', 'antitower', 'antibase'];

    var generateDudeMapping = function(fruit, type) {
      dudeMapping[fruit+type] = fruit + "-" + type + ".png";
      return;
    };

    var generateBaseMapping = function(fruit) {
      baseMapping[fruit] = fruit + "base.png";
      return;
    };
  
    for (var i in fruits) {
      for (var j in types) {
        generateDudeMapping(fruits[i], types[j]);
        generateBaseMapping(fruits[i]);
      }
    }
    
    for (var i in baseData) {
      var base = baseData[i];

      var imageName = baseMapping[base.Id];
      imageMap = this.board.imageController.getImageMap(imageName); 

      this.board.entities.push(new Tile(base.Pos.X, 
                                        base.Pos.Y, 
                                        4, 4, "#00ff00", imageMap));
    }

    for (var i in dudeData) {
      var dude = dudeData[i]; 

      var imageName = dudeMapping[dude.PlayerId+dude.Type];
      imageMap = this.board.imageController.getImageMap(imageName); 

      var dude = dudeData[i];
      this.board.entities.push(new Tile(dude.Pos.X, 
                                        dude.Pos.Y, 
                                        1, 1, "#0000ff", imageMap));
    }


    mapTileMapping = {
      s: 'Sand.png',
      g: 'concrete.png',
      m: 'ruin4.png',
      r: 'sludge.png',
      f: 'Rubble.png',
    }

    for (var i in mapData) {
      for (var j in mapData[i]) {
        var imageName = mapTileMapping[mapData[i][j]];
        imageMap = this.board.imageController.getImageMap(imageName); 
        if (imageMap) {
          this.board.mapTiles.push(new Tile(j, i, 1, 1, "#ff0000", imageMap));
        }
      }
    } 
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

});
