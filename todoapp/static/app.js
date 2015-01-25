/*global $ */
/*jshint unused:false */
var app = app || {};
var ENTER_KEY = 13;
var ESC_KEY = 27;

$(function () {
  'use strict';

  // kick things off by creating the `App`
  new app.AppView();

  // if a browser supports websocket we're going to use it :)
  if (window["WebSocket"]) {
    var conn = new WebSocket("ws://" + url + "/ws/todo");

    conn.onclose = function(evt) {
    }

    conn.onmessage = function(evt) {
      var data = {};

      try {
        data = JSON.parse(evt.data)
        if(data.__action) {
          data.__doc.id = data.__doc._id;
          if(data.__action === "create") {
            app.todos.add(data.__doc)
          }
          else if(data.__action === "update") {
            var model = app.todos.findWhere({id: data.__doc._id})
            if(model) {
              model.set('title', data.__doc.title)
              model.set('completed', data.__doc.completed)
              model.set('order', data.__doc.order)
            }
          }
          else if(data.__action === "delete") {
            var model = app.todos.findWhere({id: data.__doc._id})
            if(model) {
              model.url = function() {}
              model.destroy()
            }
          }
        }
      }
      catch (e) {
      }
    }
    conn.onerror = function(err) {
    }
    conn.onopen = function() {
    }
  }
});
