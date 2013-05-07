Ember.TEMPLATES["application"] = Ember.Handlebars.template(function anonymous(Handlebars,depth0,helpers,partials,data) {
this.compilerInfo = [2,'>= 1.0.0-rc.3'];
helpers = helpers || Ember.Handlebars.helpers; data = data || {};
  var buffer = '', hashTypes, escapeExpression=this.escapeExpression;


  data.buffer.push("<div id = \"leftHalf\">\r\n    ");
  hashTypes = {};
  data.buffer.push(escapeExpression(helpers.view.call(depth0, "App.ResourcesView", {hash:{},contexts:[depth0],types:["STRING"],hashTypes:hashTypes,data:data})));
  data.buffer.push("\r\n    \r\n    <section id = \"baseGrid\">\r\n      ");
  hashTypes = {};
  data.buffer.push(escapeExpression(helpers.view.call(depth0, "App.BaseView", {hash:{},contexts:[depth0],types:["STRING"],hashTypes:hashTypes,data:data})));
  data.buffer.push("\r\n      <section id = \"innerBaseMid\"></section>\r\n    </section>\r\n</div>\r\n    \r\n<div id = \"rightHalf\" >\r\n  <span id = \"dudecount\"> Dudes: ");
  hashTypes = {};
  data.buffer.push(escapeExpression(helpers._triageMustache.call(depth0, "dudeCount", {hash:{},contexts:[depth0],types:["ID"],hashTypes:hashTypes,data:data})));
  data.buffer.push("</span>\r\n    ");
  hashTypes = {};
  data.buffer.push(escapeExpression(helpers.view.call(depth0, "App.BuildunitView", {hash:{},contexts:[depth0],types:["STRING"],hashTypes:hashTypes,data:data})));
  data.buffer.push("\r\n</div>\r\n\r\n<div class=\"modal\" id=\"chooseTeam\">\r\n      <p> FRUIT WARS </p>\r\n      <p id = \"tagline\"> Pick your Fruit: </p>\r\n      <div class = \"teamchoose\" id = \"banana\"></div>\r\n      <div class = \"teamchoose\" id = \"grape\"></div>\r\n      <div class = \"teamchoose\" id = \"apple\"></div>\r\n      <div class = \"teamchoose\" id = \"watermelon\"></div>\r\n</div>\r\n\r\n");
  return buffer;
  
});

Ember.TEMPLATES["basetower"] = Ember.Handlebars.template(function anonymous(Handlebars,depth0,helpers,partials,data) {
this.compilerInfo = [2,'>= 1.0.0-rc.3'];
helpers = helpers || Ember.Handlebars.helpers; data = data || {};
  var buffer = '';


  return buffer;
  
});

Ember.TEMPLATES["buildunit"] = Ember.Handlebars.template(function anonymous(Handlebars,depth0,helpers,partials,data) {
this.compilerInfo = [2,'>= 1.0.0-rc.3'];
helpers = helpers || Ember.Handlebars.helpers; data = data || {};
  


  data.buffer.push("\r\n\r\n");
  
});

Ember.TEMPLATES["resources"] = Ember.Handlebars.template(function anonymous(Handlebars,depth0,helpers,partials,data) {
this.compilerInfo = [2,'>= 1.0.0-rc.3'];
helpers = helpers || Ember.Handlebars.helpers; data = data || {};
  var buffer = '', hashTypes, escapeExpression=this.escapeExpression;


  data.buffer.push("<p id = \"resourceTitle\"> Resources </p>\r\n<section class = \"row-fluid\">\r\n<section class = \"span6\" id = \"wateramount\"> Water: ");
  hashTypes = {};
  data.buffer.push(escapeExpression(helpers._triageMustache.call(depth0, "waterAmount", {hash:{},contexts:[depth0],types:["ID"],hashTypes:hashTypes,data:data})));
  data.buffer.push(" </section>\r\n<section class = \"span6\" id = \"time\"> Time: ");
  hashTypes = {};
  data.buffer.push(escapeExpression(helpers._triageMustache.call(depth0, "time", {hash:{},contexts:[depth0],types:["ID"],hashTypes:hashTypes,data:data})));
  data.buffer.push(" </section>\r\n</section>");
  return buffer;
  
});