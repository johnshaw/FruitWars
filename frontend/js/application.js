oCanvas.domReady( function() {
  window.game = new Game();
  window.adapter = new SocketAdapter(game, "ws://192.168.2.13:8080/screen");

  game.init();

  game.canvas = oCanvas.create({
    canvas: "#game",
    background: "#123456",
    clearEachFrame: true,
    drawEachFrame: true,
    fps: .1,
    disableScrolling: true 
  });

  game.canvas.setLoop(function () {
    game.entityManager.updateEntities();
    game.render();
  }).start();
});

//GAME CLASS
var Game = {};

Game = function () {
  this.entityManager = new EntityManager(this);
};

Game.prototype = {
  
  constructor: Game,
  
  init: function() {
    this.initCanvas();
  },

  //initialize the drawing context
  initCanvas: function () {
    var height = document.height
      , width = document.width
      , $canvas;
    this.pDensityX = 32;
    this.pDensityY = 18;
    this.canvasId = "game";
   
    this.pixelSize = calculatePixelSize(height, 
                                        this.pDensityY);
    canvasSize = updateCanvasSize(this.pDensityX, 
                                  this.pDensityY, 
                                  this.pixelSize);

    $canvas = $("<canvas id='"+this.canvasId+"' height="+canvasSize.height+" width="+canvasSize.width+"></canvas>");
    $canvas.appendTo('body');
  },

  render: function() {
    console.log('rendering');
  },
};

var EntityManager = {};

EntityManager = function (game) {
  this.game = game;
  var baseCount = 4
    , dudeCount = 200
    , playerCount = 4
  this.playerStore = [];
  this.baseStore = [];
  this.dudeStore = [];

  this.activePlayers = [];
  this.activeBases = [];
  this.activeDudes = [];

  //initialize an array of bases
  for (var i=0; i<playerCount; i++) {
    this.playerStore.push(new Player(this.game));
  }
  for (var i=0; i<baseCount; i++) {
    this.baseStore.push(new Base(this.game));
  }
  //initialize an array of dudes 
  for (var i=0; i<dudeCount; i++) {
    this.dudeStore.push(new Dude(this.game));
  }
};

EntityManager.prototype = {

  parseData: function (data) {
    if (data.Bases) {
      this.parseBaseData(data.Bases);
    }
  },

  parseBaseData: function (data) {

    //reset bases
    this.activeBases.forEach( function(base) {
      base.accountedFor = false;
    }); 
    
    //loop over data looking for activeBases
    for (var baseIndex in data) {

      var newBase
        , baseData = data[baseIndex];

      var existingBase = this.activeBases.filter( function(base) {
        return (base.Id === baseData.Id) ? true : false;
      });
      
      if (existingBase.length > 0) {
        console.log('existingBase ', existingBase.length);
        existingBase.updateParams();
        existingBase.accountedFor = true;
      } else {
        console.log("baseStore ", this.baseStore);
        newBase = this.baseStore.pop();
        newBase.updateParams(baseData);
        this.activeBases.push(newBase);
      }
    }
  },

  updateEntities: function() {
    //find Objects that were not accounted for during parsing
    console.log(this.activeBases);
    var removedObjects = this.activeBases.filter( function(base) {
      return !base.accountedFor;
    });

    //reset all params on these objects
    removedObjects.forEach( function(base) {
      base.reset();
    });
  },
};

var Base = {};
var Player = {};
var Dude = {};

Player = function(game) {
  this.game = game; 
}

Base = function(game) {
  this.game = game;
  this.size = game.pixelSize*4;
}

Dude = function(game) {
  this.game = game;
  this.size = game.pixelSize;
}

Player.prototype = {
  constructor: Player,  

  reset: function () {
    var keys = Object.keys(this);
    for (var i=0; i<keys; i++) {
      keys[i].delete;
    }
  },

  updateParams: function (data) {
    console.log('base update params called');
  },
};

Base.prototype = {
  constructor: Base,

  reset: function () {
    var keys = Object.keys(this);
    for (var i=0; i<keys; i++) {
      keys[i].delete;
    }
  },

  updateParams: function (data) {
    this.Id = data.Id;
    console.log(this.Id);
  },
};

Dude.prototype = {
  constructor: Dude,

  reset: function () {
    var keys = Object.keys(this);
    for (var i=0; i<keys; i++) {
      keys[i].delete;
    }
  },

  updateParams: function (data) {
    console.log('base update params called');
  },
};
//SOCKET ADAPTER OBJECT 
var SocketAdapter = {};

SocketAdapter = function(game, url) {

  this.game = game;
  
  this.ws = new WebSocket(url);

  this.ws.onopen = function(event) {
    console.log("connecion established");
  }.bind(this);

  this.ws.onmessage = function(event) {
    this.game.entityManager.parseData(JSON.parse(event.data));
  }.bind(this);

};

SocketAdapter.prototype = {
  constructor: SocketAdapter
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
