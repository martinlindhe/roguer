!function e(t,s,a){function i(o,n){if(!s[o]){if(!t[o]){var l="function"==typeof require&&require;if(!n&&l)return l(o,!0);if(r)return r(o,!0);var h=new Error("Cannot find module '"+o+"'");throw h.code="MODULE_NOT_FOUND",h}var u=s[o]={exports:{}};t[o][0].call(u.exports,function(e){var s=t[o][1][e];return i(s?s:e)},u,u.exports,e,t,s,a)}return s[o].exports}for(var r="function"==typeof require&&require,o=0;o<a.length;o++)i(a[o]);return i}({1:[function(e,t,s){"use strict";function a(e,t){if(!(e instanceof t))throw new TypeError("Cannot call a class as a function")}var i=function(){function e(e,t){for(var s=0;s<t.length;s++){var a=t[s];a.enumerable=a.enumerable||!1,a.configurable=!0,"value"in a&&(a.writable=!0),Object.defineProperty(e,a.key,a)}}return function(t,s,a){return s&&e(t.prototype,s),a&&e(t,a),t}}();Object.defineProperty(s,"__esModule",{value:!0});var r=function(){function e(t){a(this,e),this.gameState=t,this.socket=new WebSocket("ws://localhost:3322/ws"),this.sessionToken=window.sessionStorage.getItem("_token");var s=this;this.socket.onmessage=function(e){var t=JSON.parse(e.data);switch(t.Type){case"xy":s.handleXyMessage(t);break;case"move_res":s.handleMoveResMessage(t);break;case"ok":console.log("server OK: "+e.data);break;case"tick":s.gameState.setServerTime(t.FormattedTime);break;case"msg":s.gameState.messageToLog({time:t.Time,text:t.Message});break;default:console.log("<-recv- "+e.data),console.log("unknown command from server: "+t.Type)}},this.socket.onopen=function(){console.log("Websocket connected"),s.sessionToken?(console.log("Resuming session"),this.send("continue "+s.sessionToken)):this.send("new_player "+s.gameState.playerName)}}return i(e,[{key:"sendMove",value:function(){var e=Math.floor(this.gameState.playerGroup.x/this.gameState.tileWidth),t=Math.floor(this.gameState.playerGroup.y/this.gameState.tileHeight);(this.prevX!=e||this.prevY!=t)&&(this.socket.send("move "+e+" "+t+" "+this.sessionToken),this.prevX=e,this.prevY=t)}},{key:"handleXyMessage",value:function(e){this.gameState.spawnPlayer(e),this.sessionToken=e.Token,window.sessionStorage.setItem("_token",e.Token)}},{key:"handleMoveResMessage",value:function(e){this.gameState.renderLocalSpawns(e.LocalSpawns)}}]),e}();s["default"]=r},{}],2:[function(e,t,s){"use strict";function a(e){return e&&e.__esModule?e:{"default":e}}function i(e,t){if(!(e instanceof t))throw new TypeError("Cannot call a class as a function")}function r(e,t){if(!e)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return!t||"object"!=typeof t&&"function"!=typeof t?e:t}function o(e,t){if("function"!=typeof t&&null!==t)throw new TypeError("Super expression must either be null or a function, not "+typeof t);e.prototype=Object.create(t&&t.prototype,{constructor:{value:e,enumerable:!1,writable:!0,configurable:!0}}),t&&(Object.setPrototypeOf?Object.setPrototypeOf(e,t):e.__proto__=t)}var n=function(){function e(e,t){for(var s=0;s<t.length;s++){var a=t[s];a.enumerable=a.enumerable||!1,a.configurable=!0,"value"in a&&(a.writable=!0),Object.defineProperty(e,a.key,a)}}return function(t,s,a){return s&&e(t.prototype,s),a&&e(t,a),t}}();Object.defineProperty(s,"__esModule",{value:!0});var l=e("./Client.js"),h=a(l),u=e("./MessageLog.js"),c=a(u),d=function(e){function t(){return i(this,t),r(this,Object.getPrototypeOf(t).apply(this,arguments))}return o(t,e),n(t,[{key:"preload",value:function(){this.game.time.advancedTiming=!0,this.game.stage.backgroundColor="#262f71",this.game.load.tilemap("islandMap","/island/full",null,Phaser.Tilemap.TILED_JSON),this.game.load.image("ground","img/tileset/oddball/ground.png",4,8),this.game.load.image("minimap","img/islands/current.png"),this.game.load.atlas("characterAtlas","img/tileset/oddball/characters.png","sprite/character"),this.game.load.atlas("itemAtlas","img/tileset/oddball/items.png","sprite/item"),this.game.load.atlas("ground2Atlas","img/tileset/oddball/ground2.png","sprite/ground2"),this.game.load.spritesheet("button","img/tileset/ui/buttons.png",27,24),this.game.load.image("oddballFont","img/tileset/oddball/font.png"),this.game.load.audio("bgSound",["audio/dead_feelings.mp3"])}},{key:"create",value:function(){this.playerName="Jimpson",this.worldScale=1,this.tileWidth=8,this.tileHeight=4,this.maxMessages=15,this.logTextHeight=15,this.serverTime=0,this.game.physics.startSystem(Phaser.Physics.ARCADE),this.game.scale.setGameSize(window.innerWidth,window.innerHeight),this.music=this.game.add.audio("bgSound"),this.music.volume=.2,this.stageGroup=this.game.add.group(),this.groundMap=this.game.add.tilemap("islandMap"),this.groundMap.addTilesetImage("island_tiles","ground"),this.groundMap.setCollisionBetween(0,112),this.groundLayer=this.groundMap.createLayer(0),this.groundLayer.resizeWorld(),this.stageGroup.add(this.groundLayer),this.playerGroup=this.game.add.group(),this.game.camera.follow(this.playerGroup),this.stageGroup.add(this.playerGroup),this.spawnLayer=this.game.add.group(),this.stageGroup.add(this.spawnLayer),this.initUi(),this.client=new h["default"](this)}},{key:"update",value:function(){if(this.playerSprite){this.game.physics.arcade.collide(this.playerGroup,this.groundLayer);var e=1,t=e/2;this.cursors.up.isDown?(this.playerGroup.y-=t,this.client.sendMove()):this.cursors.down.isDown&&(this.playerGroup.y+=t,this.client.sendMove()),this.cursors.left.isDown?(this.playerGroup.x-=e,this.client.sendMove()):this.cursors.right.isDown&&(this.playerGroup.x+=e,this.client.sendMove()),this.game.input.keyboard.isDown(Phaser.Keyboard.Q)?this.worldScale+=.05:this.game.input.keyboard.isDown(Phaser.Keyboard.A)&&(this.worldScale-=.05),this.worldScale=Phaser.Math.clamp(this.worldScale,.5,4),this.stageGroup.scale.set(this.worldScale)}}},{key:"render",value:function(){this.game.debug.text(this.game.time.fps||"--",1,14,"#00ff00")}},{key:"messageToLog",value:function(e){this.logMessageList.text=this.messageLog.add(e).render()}},{key:"initUi",value:function(){this.uiGroup=this.game.add.group(),this.cursors=this.game.input.keyboard.createCursorKeys();var e=3,t=this.game.width-this.game.cache.getImage("minimap").width/e,s=this.game.add.sprite(t,0,"minimap");s.fixedToCamera=!0,s.scale.set(1/e),s.alpha=.8,s.setScaleMinMax(1/e,1/e),this.uiGroup.add(s);var a=this.game.add.button(this.game.width-102,2,"button",function(){this.music.isPlaying?this.music.stop():this.music.play()},this,0,0,0);a.fixedToCamera=!0,this.uiGroup.add(a);var i={font:"10px topaz",fill:"#fff",backgroundColor:"rgba(0,0,0,0.25)",wordWrap:!0,wordWrapWidth:400};this.logMessageList=this.game.add.text(0,0,"",i),this.logMessageList.stroke="#000000",this.logMessageList.strokeThickness=2,this.logMessageList.lineSpacing=-8,this.logMessageList.x=this.game.width-400,this.logMessageList.y=this.game.height-this.maxMessages*this.logTextHeight,this.logMessageList.fixedToCamera=!0,this.uiGroup.add(this.logMessageList),this.messageLog=new c["default"],this.messageToLog({time:0,text:"Welcome to roguer!"}),this.serverTimeText=this.game.add.text(this.game.width-210,0,"",i),this.serverTimeText.fixedToCamera=!0,this.setServerTime(0)}},{key:"setServerTime",value:function(e){this.serverTime=e,this.serverTime&&(this.serverTimeText.text="time: "+this.serverTime),this.messageLog.save()}},{key:"spawnPlayer",value:function(e){this.playerSprite=this.game.add.sprite(0,0,"characterAtlas"),this.playerSprite.frameName="dwarf",this.playerSprite.anchor.set(.5),this.playerGroup.x=e.X*this.tileWidth,this.playerGroup.y=e.Y*this.tileHeight,this.playerGroup.add(this.playerSprite);var t=this.makeText(this.playerName),s=this.game.add.image(0,-10,t);s.anchor.set(.5),this.playerGroup.add(s),console.log("spawned at "+e.X+", "+e.Y),this.game.physics.enable(this.playerGroup),this.renderLocalSpawns(e.LocalSpawns)}},{key:"renderLocalSpawns",value:function(e){this.spawnLayer.removeAll();for(var t="",s=0;s<e.length;s++){var a=e[s];if("player"!=a.Class||a.Name!=this.playerName){var i=a.Sprite.split(":");switch(i[0]){case"c":t="characterAtlas";break;case"i":t="itemAtlas";break;case"g":t="ground2Atlas";break;default:console.log("ERROR unknown sprite: "+a.Sprite),console.log(a);continue}var r=this.game.add.sprite(0,0,t);r.x=a.X*this.tileWidth,r.y=a.Y*this.tileHeight,r.frameName=i[1],r.anchor.set(.5),this.spawnLayer.add(r)}}}},{key:"makeText",value:function(e){var t="                                !\"#$%&'()  ,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",s=this.game.add.retroFont("oddballFont",8,8,t,16);return s.autoUpperCase=!1,s.text=e,s}},{key:"updateShadowTexture",value:function(){this.shadowTexture.context.fillStyle="rgb(100, 100, 100)",this.shadowTexture.context.fillRect(0,0,this.game.width,this.game.height);var e=this.shadowTexture.context.createRadialGradient(this.playerGroup.x,this.playerGroup.y,.75*this.LIGHT_RADIUS,this.playerGroup.x,this.playerGroup.y,this.LIGHT_RADIUS);e.addColorStop(0,"rgba(255, 255, 255, 1.0)"),e.addColorStop(1,"rgba(255, 255, 255, 0.0)"),this.shadowTexture.context.beginPath(),this.shadowTexture.context.fillStyle=e,this.shadowTexture.context.arc(this.playerGroup.x,this.playerGroup.y,this.LIGHT_RADIUS,0,2*Math.PI),this.shadowTexture.context.fill(),this.shadowTexture.dirty=!0}}]),t}(Phaser.State);s["default"]=d},{"./Client.js":1,"./MessageLog.js":3}],3:[function(e,t,s){"use strict";function a(e,t){if(!(e instanceof t))throw new TypeError("Cannot call a class as a function")}var i=function(){function e(e,t){for(var s=0;s<t.length;s++){var a=t[s];a.enumerable=a.enumerable||!1,a.configurable=!0,"value"in a&&(a.writable=!0),Object.defineProperty(e,a.key,a)}}return function(t,s,a){return s&&e(t.prototype,s),a&&e(t,a),t}}();Object.defineProperty(s,"__esModule",{value:!0});var r=function(){function e(){a(this,e),this.logMessages=[];var t=window.sessionStorage.getItem("_messages");t&&(console.log("restoring saved msgs from "+t),this.logMessages=JSON.parse(t))}return i(e,[{key:"add",value:function(e){return this.logMessages?(this.logMessages.push(e),this):(console.log("error: log wnd not yet ready!"),void console.log(e))}},{key:"save",value:function(){window.sessionStorage.setItem("_messages",JSON.stringify(this.logMessages))}},{key:"render",value:function(){if(!this.logMessages)return void console.log("error: log messages not yet ready!");this.logMessages=this.logMessages.slice(-this.maxMessages);var e="",t=!0,s=!1,a=void 0;try{for(var i,r=this.logMessages[Symbol.iterator]();!(t=(i=r.next()).done);t=!0){var o=i.value;e=e+o.time+": "+o.text+"\n"}}catch(n){s=!0,a=n}finally{try{!t&&r["return"]&&r["return"]()}finally{if(s)throw a}}return e.trim()}}]),e}();s["default"]=r},{}],4:[function(e,t,s){"use strict";function a(e){return e&&e.__esModule?e:{"default":e}}function i(e,t){if(!(e instanceof t))throw new TypeError("Cannot call a class as a function")}function r(e,t){if(!e)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return!t||"object"!=typeof t&&"function"!=typeof t?e:t}function o(e,t){if("function"!=typeof t&&null!==t)throw new TypeError("Super expression must either be null or a function, not "+typeof t);e.prototype=Object.create(t&&t.prototype,{constructor:{value:e,enumerable:!1,writable:!0,configurable:!0}}),t&&(Object.setPrototypeOf?Object.setPrototypeOf(e,t):e.__proto__=t)}var n=e("./GameState.js"),l=a(n),h=function(e){function t(){i(this,t);var e=r(this,Object.getPrototypeOf(t).call(this,800,400,Phaser.CANVAS,"content",null,!1,!1));return e.state.add("GameState",l["default"],!1),e.state.start("GameState"),e}return o(t,e),t}(Phaser.Game);new h},{"./GameState.js":2}]},{},[4]);
//# sourceMappingURL=app.js.map
