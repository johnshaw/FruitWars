minispade.register('application/Application.js', function() {
window.App = Ember.Application.create({
  customEvents: {
    'touchenter': "touchEnter",
    'touchleave': "touchLeave"
  }
});minispade.require('controllers/ApplicationController.js');minispade.require('controllers/BaseController.js');minispade.require('models/Tower.js');minispade.require('models/Player.js');minispade.require('models/Dude.js');minispade.require('views/TowerView.js');minispade.require('views/BaseView.js');minispade.require('views/ResourcesView.js');minispade.require('views/BuildunitView.js');minispade.require('store/Router.js');minispade.require('store/adapter.js');
});

minispade.register('controllers/ApplicationController.js', function() {
App.ApplicationController = Ember.Controller.extend({
  sampleText: "hey this application controller is sure neat",
  Player: null,
  dudeCount: 1,
  newTower: null
});
});

minispade.register('controllers/BaseController.js', function() {
App.BaseController = Em.ArrayController.extend();
});

minispade.register('models/Dude.js', function() {
App.Dude = DS.Model.extend({
  primaryKey: "PlayerId",
  PlayerId: DS.attr('string'),
  Type: DS.attr('string')
});
});

minispade.register('models/Player.js', function() {
App.Player = DS.Model.extend({
  primaryKey: "PlayerId",
  PlayerId: DS.attr('string'),
  name: DS.attr('string'),
  team: DS.attr('string'),
  water: DS.attr('number')
});
});

minispade.register('models/Tower.js', function() {
App.Tower = Ember.Object.extend({
  top: 0,
  left: 0,
  position: 0
});
});

minispade.register('store/Router.js', function() {
App.IndexRoute = Ember.Route.extend({
  enter: function() {
    window.ws1 = App.__container__.lookup('controller:base').get('store.adapter.socket');
    return Ember.run.later(this, function() {
      $('#chooseTeam').modal();
      $('#banana').on('click touch', function() {
        var charSelected, msg;

        charSelected = JSON.stringify({
          PlayerId: "banana"
        });
        msg = JSON.stringify({
          MsgType: "SelectPlayer",
          data: "" + charSelected
        });
        console.log(msg);
        return ws1.send(msg);
      });
      $('#grape').on('click touch', function() {
        var charSelected, msg;

        charSelected = JSON.stringify({
          PlayerId: "grape"
        });
        msg = JSON.stringify({
          MsgType: "SelectPlayer",
          data: "" + charSelected
        });
        console.log(msg);
        return ws1.send(msg);
      });
      $('#apple').on('click touch', function() {
        var charSelected, msg;

        charSelected = JSON.stringify({
          PlayerId: "apple"
        });
        msg = JSON.stringify({
          MsgType: "SelectPlayer",
          data: "" + charSelected
        });
        console.log(msg);
        return ws1.send(msg);
      });
      return $('#watermelon').on('click touch', function() {
        var charSelected, msg;

        charSelected = JSON.stringify({
          PlayerId: "watermelon"
        });
        msg = JSON.stringify({
          MsgType: "SelectPlayer",
          data: "" + charSelected
        });
        console.log(msg);
        return ws1.send(msg);
      });
    }, 250);
  },
  setupController: function(controller, model) {
    var appCon;

    appCon = this.controllerFor('application');
    return appCon.set("Player", DS.defaultStore.all(App.Player));
  }
});
});

minispade.register('store/adapter.js', function() {
App.SocketAdapter = DS.RESTAdapter.extend({
  socket: void 0,
  find: function(type) {
    return DS.defaultStore.all(type);
  },
  init: function() {
    var context;

    this._super();
    context = this;
    window.ws = new WebSocket("ws://" + document.location.host + "/control");
    ws.onopen = function(err) {
      return console.log("connected");
    };
    ws.onerror = function(err) {
      var ws;

      ws = new WebSocket("ws://" + document.location.host + "/control");
      return this.set("socket", ws);
    };
    ws.onmessage = function(evt) {
      var actionType, appCon, data, newTower, newplayer;

      window.received_msg = evt.data;
      window.store = DS.defaultStore;
      window.parsed = $.parseJSON(received_msg);
      actionType = parsed.MsgType;
      data = JSON.parse(parsed.Data);
      switch (actionType) {
        case "ConfirmPlayer":
          newplayer = data;
          console.log(data.PlayerId, "Playerrr");
          newplayer.player_id = data.PlayerId;
          delete newplayer.PlayerId;
          store.load(App.Player, newplayer);
          return context.loadPlayer(data.PlayerId);
        case "RejectPlayer":
          return console.log("reject");
        case "BuyTowerConfirm":
          newTower = data;
          console.log(newTower.Pos);
          appCon = App.__container__.lookup('controller:application');
          return appCon.set('newTower', newTower.Pos);
      }
    };
    ws.onclose = function(err) {
      return console.log("closed");
    };
    return this.set("socket", ws);
  },
  loadPlayer: function(newplayer) {
    return $.modal.close();
  }
});

App.Store = DS.Store.extend({
  revision: 12,
  adapter: App.SocketAdapter.create()
});
});

minispade.register('views/BaseView.js', function() {
App.BaseView = Ember.ContainerView.extend({
  childViews: ['topLTower', 'topMLTower', 'topMRTower', 'topRTower', 'upLTower', 'upRTower', 'lowLTower', 'lowRTower', 'BotLTower', 'BotMLTower', 'BotMRTower', 'BotRTower'],
  classNames: "towerGrid",
  tagName: 'section',
  newTower: (function() {
    var newtower, player;

    newtower = this.get('controller.newTower');
    player = this.get('controller.Player').toArray()[0];
    console.log("PLAYERTOWER", player.get('PlayerId'));
    if (newtower !== null) {
      switch (player.get('PlayerId')) {
        case "banana":
          $('.pos' + newtower).removeClass("sprite-bananabuildtower");
          return $('.pos' + newtower).addClass("sprite-Tower-Banana");
        case "grape":
          $('.pos' + newtower).removeClass("sprite-grapebuildtower");
          return $('.pos' + newtower).addClass("sprite-Tower-Grape");
        case "watermelon":
          $('.pos' + newtower).removeClass("sprite-watermelonbuildtower");
          return $('.pos' + newtower).addClass("sprite-Tower-Watermelon");
        case "apple":
          $('.pos' + newtower).removeClass("sprite-applebuildtower");
          return $('.pos' + newtower).addClass("sprite-Tower-Apple");
      }
    }
  }).observes('controller.newTower'),
  topLTower: App.TowerView.create({
    content: App.Tower.create({
      position: 0,
      top: 0,
      left: 0
    })
  }),
  topMLTower: App.TowerView.create({
    content: App.Tower.create({
      position: 1,
      top: 0,
      left: "25%"
    })
  }),
  topMRTower: App.TowerView.create({
    content: App.Tower.create({
      position: 2,
      top: 0,
      left: "50%"
    })
  }),
  topRTower: App.TowerView.create({
    content: App.Tower.create({
      position: 3,
      top: 0,
      left: "75%"
    })
  }),
  upLTower: App.TowerView.create({
    content: App.Tower.create({
      position: 11,
      top: "25%",
      left: 0
    })
  }),
  upRTower: App.TowerView.create({
    content: App.Tower.create({
      position: 4,
      top: "25%",
      left: "75%"
    })
  }),
  lowLTower: App.TowerView.create({
    content: App.Tower.create({
      position: 10,
      top: "50%",
      left: 0
    })
  }),
  lowRTower: App.TowerView.create({
    content: App.Tower.create({
      position: 5,
      top: "50%",
      left: "75%"
    })
  }),
  BotLTower: App.TowerView.create({
    content: App.Tower.create({
      position: 9,
      top: "75%",
      left: 0
    })
  }),
  BotMLTower: App.TowerView.create({
    content: App.Tower.create({
      position: 8,
      top: "75%",
      left: "25%"
    })
  }),
  BotMRTower: App.TowerView.create({
    content: App.Tower.create({
      position: 7,
      top: "75%",
      left: "50%"
    })
  }),
  BotRTower: App.TowerView.create({
    content: App.Tower.create({
      position: 6,
      top: "75%",
      left: "75%"
    })
  })
});
});

minispade.register('views/BuildunitView.js', function() {
App.UnitView = Em.View.extend({
  templateName: 'buildunit',
  tagName: 'section',
  type: null,
  dudeCount: 1,
  startX: 0,
  startY: 0,
  deltaX: 0,
  deltaY: 0,
  moveX: 0,
  moveY: 0,
  touchStart: function(e) {
    var dudeCount, _this;

    e.preventDefault();
    this.set('startX', e.originalEvent.touches[0].pageX);
    this.set('startY', e.originalEvent.touches[0].pageY);
    dudeCount = 1;
    _this = this;
    return this.dude = setInterval(function() {
      dudeCount += 1;
      _this.set("dudeCount", dudeCount);
      return _this.set("controller.dudeCount", dudeCount);
    }, 100);
  },
  touchMove: function(e) {
    return e.preventDefault();
  },
  mouseDown: function(e) {
    var dudeCount, _this;

    dudeCount = 1;
    _this = this;
    return this.dude = setInterval(function() {
      dudeCount += 1;
      _this.set("dudeCount", dudeCount);
      return _this.set("controller.dudeCount", dudeCount);
    }, 100);
  },
  mouseUp: function(e) {
    var a, b, dudeCount, i;

    window.ws = this.get('controller.store.adapter.socket');
    i = 0;
    dudeCount = this.get('dudeCount');
    clearInterval(this.dude);
    while (dudeCount > 0) {
      b = JSON.stringify({
        Type: "seedling"
      });
      a = JSON.stringify({
        MsgType: "BuyDude",
        data: b
      });
      ws.send(a);
      dudeCount = dudeCount - 1;
    }
    return this.set("dudeCount", 1);
  },
  touchEnd: function(e) {
    var data, deltaX, deltaY, direction, dudeCount, msg, startX, startY, ws;

    dudeCount = this.get('dudeCount');
    clearInterval(this.dude);
    e.preventDefault();
    startX = this.get('startX');
    startY = this.get('startY');
    deltaX = e.originalEvent.changedTouches[0].pageX - startX;
    deltaY = e.originalEvent.changedTouches[0].pageY - startY;
    if (deltaX !== 0 && deltaY !== 0) {
      this.set('deltaX', deltaX);
      this.set('deltaY', deltaY);
      ws = this.get('controller.store.adapter.socket');
      direction = {
        X: deltaX,
        Y: deltaY
      };
      data = JSON.stringify({
        Type: this.get('type'),
        Count: dudeCount,
        Dir: direction
      });
      msg = JSON.stringify({
        MsgType: "BuyDude",
        data: data
      });
      ws.send(msg);
    }
    this.set("dudeCount", 1);
    return this.set("controller.dudeCount", 1);
  }
});

App.BuildunitView = Em.ContainerView.extend({
  childViews: ['leftUnit', 'middleUnit', 'rightUnit'],
  tagName: 'section',
  classNames: 'buildunitcontainer',
  leftUnit: App.UnitView.create({
    classNames: 'antidude',
    type: 'antidude'
  }),
  middleUnit: App.UnitView.create({
    classNames: 'antitower',
    type: 'antitower'
  }),
  rightUnit: App.UnitView.create({
    classNames: 'antibase',
    type: 'antibase'
  })
});
});

minispade.register('views/ResourcesView.js', function() {
App.ResourcesView = Em.View.extend({
  templateName: "resources",
  tagName: 'section',
  classNames: 'resources'
});
});

minispade.register('views/TowerView.js', function() {
App.TowerView = Em.View.extend({
  classNames: ['baseTower', 'contsprite'],
  templateName: 'basetower',
  classNameBindings: ['player', 'pos'],
  tagName: 'section',
  attributeBindings: ['style'],
  pos: (function() {
    return "pos" + this.get('content.position');
  }).property('content.position'),
  player: (function() {
    var player, type;

    player = this.get('controller.Player').toArray()[0];
    if (player) {
      type = player.get('PlayerId');
      switch (type) {
        case "banana":
          $("body").css("background", "#aaa671");
          $("#innerBaseMid").addClass("contsprite sprite-Nexus-Banana");
          $(".antidude").addClass("contsprite sprite-bananadude-dude-icon");
          $(".antitower").addClass("contsprite sprite-bananadude-tower-icon");
          $(".antibase").addClass("contsprite sprite-bananadude-base-icon");
          return "sprite-bananabuildtower";
        case "grape":
          $("body").css("background", "#9d71a8");
          $("#innerBaseMid").addClass("contsprite sprite-Nexus-Grape");
          $(".antidude").addClass("contsprite sprite-grapedude-dude-icon");
          $(".antitower").addClass("contsprite sprite-grapedude-tower-icon");
          $(".antibase").addClass("contsprite sprite-grapedude-base-icon");
          return "sprite-grapebuildtower";
        case "apple":
          $("body").css("background", "#a87171");
          $("#innerBaseMid").addClass("contsprite sprite-Nexus-Apple");
          $(".antidude").addClass("contsprite sprite-appledude-dude-icon");
          $(".antitower").addClass("contsprite sprite-appledude-tower-icon");
          $(".antibase").addClass("contsprite sprite-appledude-base-icon");
          return "sprite-applebuildtower";
        case "watermelon":
          $("body").css("background", "#81936a");
          $("#innerBaseMid").addClass("contsprite sprite-Nexus-Watermelon");
          $(".antidude").addClass("contsprite sprite-watermelondude-dude-icon");
          $(".antitower").addClass("contsprite sprite-watermelondude-tower-icon");
          $(".antibase").addClass("contsprite sprite-watermelondude-base-icon");
          return "sprite-watermelonbuildtower";
      }
    }
    return false;
  }).property('controller.Player.@each'),
  built: false,
  touchStart: function(e) {
    return this.$().addClass("scaleXYUP");
  },
  touchEnd: function(e) {
    var data, msg, ws;

    this.$().removeClass("scaleXYUP");
    ws = this.get('controller.store.adapter.socket');
    data = JSON.stringify({
      Pos: this.get('content.position')
    });
    msg = JSON.stringify({
      MsgType: "BuyTower",
      Data: data
    });
    console.log(msg);
    return ws.send(msg);
  },
  touchCancel: function(e) {
    return this.$().removeClass("scaleXYUP");
  },
  style: (function() {
    var left, top;

    top = this.get('content.top');
    left = this.get('content.left');
    return "top:" + top + ";left:" + left + ";";
  }).property('content.top', 'content.left')
});
});
