'use strict';

/**
 * @ngdoc service
 * @name mgmApp.mgm
 * @description
 * # mgm
 * Service in the mgmApp.
 */
angular.module('mgmApp').service('mgm', function ($location, $rootScope) {
  console.log("mgm service instantiated");

  var remoteURL = "ws://" + $location.host() + ":" + $location.port() + "/ws";

  self = this;

  self.regions = {};
  self.estates = {};
  self.activeUsers = {};
  self.suspendedUsers = {};
  self.pendingUsers = {};
  self.groups = {};
  self.hosts = {};
  self.serverConnected = false;

  $rootScope.$on("AuthChange", function (event, value) {
    if (value === false) {
      //we logged out
      self.regions = {};
      self.estates = {};
      self.activeUsers = {};
      self.suspendedUsers = {};
      self.pendingUsers = {};
      self.groups = {};
      self.hosts = {};
      self.serverConnected = false;
    }
  });

  this.connect = function () {
    $rootScope.$broadcast("SyncBegin");
    console.log("Connecting to: " + remoteURL);
    self.ws = new ReconnectingWebSocket(remoteURL);

    self.ws.onopen = function () {
      console.log("Socket has been opened!");
      $rootScope.$broadcast("ServerConnected");
      self.serverConnected = true;
    };

    self.ws.onmessage = function (evt) {
      var message = $.parseJSON(evt.data);
      switch (message.MessageType) {
      case "UserUpdate":
        var user = message.Message;
        if (user.Suspended) {
          self.suspendedUsers[user.UserID] = user;
          if (user.UserID in self.activeUsers) {
            delete self.activeUsers[user.UserID];
          }
        } else {
          self.activeUsers[user.UserID] = user;
          if (user.UserID in self.suspendedUsers) {
            delete self.suspendedUsers[user.UserID];
          }
        }
        $rootScope.$broadcast("UserUpdate", user);
        break
      case "PendingUserUpdate":
        self.pendingUsers[message.Message.UserID] = message.Message;
        $rootScope.$broadcast("PendingUserUpdate", message.Message);
        break;
      case "RegionUpdate":
        self.regions[message.Message.UUID] = message.Message;
        $rootScope.$broadcast("RegionUpdate", message.Message);
        break;
      case "EstateUpdate":
        self.estates[message.Message.ID] = message.Message;
        $rootScope.$broadcast("EstateUpdate", message.Message);
        break;
      case "GroupUpdate":
        self.groups[message.Message.ID] = message.Message;
        $rootScope.$broadcast("GroupUpdate", message.Message);
        break;
      case "ConfigUpdate":
        $rootScope.$broadcast("ConfigUpdate", message.Message);
        break;
      case "HostUpdate":
        message.Message.Status = JSON.parse(message.Message.Status);
        self.hosts[message.Message.Address] = message.Message;
        $rootScope.$broadcast("HostUpdate", message.Message);
        break;
      case "SyncComplete":
        $rootScope.$broadcast("SyncComplete");
        break;
      default:
        console.log("Error parsing message:");
        console.log(message);
      };

    }

    self.ws.onclose = function (message) {
      console.log("Connection closed");
      self.serverConnected = false;
    }
  };

  this.disconnect = function () {
    self.ws.close();
  };

  this.request = function (req) {
    self.ws.send(JSON.stringify(req));
  }

  var locationStack = new Array();
  this.pushLocation = function (url) {
    locationStack.push(url);
  }
  this.popLocation = function () {
    return locationStack.pop();
  }
});
