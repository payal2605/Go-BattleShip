var socketProtocol = "ws:"

if (document.location.protocol == "https:") {
  socketProtocol = "wss:"
}

var socket = new WebSocket(socketProtocol + "//" + document.location.host + "/ws")
var count = 0
socket.onmessage = function (e) {
  var data = JSON.parse(e.data)
  var task = data.task
  var message = data.message
  var turn = data.turn
  console.log(count);

  if (task == "sendcoordinates") {
    message = message.replace("[", "")
    message = message.replace("]", "")
    message = message.replace(/"/g, "")
    message = message.split(",")

    for (var i = 0; i < message.length; i++) {
      document.getElementById(message[i] + "me").style.backgroundColor = "#008000"
    }
    console.log(turn)
    if (turn == "yourturn") {
      document.getElementById("result").innerHTML = "It's your turn"
    }
    else {
      document.getElementById("result").innerHTML = "It's opponent's turn"
    }

  }
  if (task == "hit") {
    if (turn == "me") {
      document.getElementById(message).innerHTML = "<span style='font-size:30px;'>&#8226</span>"
    }
    else {
      document.getElementById(message + "me").innerHTML = "<span style='font-size:30px;'>&#8226</span>"
    }
  }
  if (task == "miss") {

    if (turn == "me") {
      document.getElementById(message).innerHTML = "<td height='36.1px' width='36.1px'>&#10005</td>"
    } else {
      document.getElementById(message + "me").innerHTML = "<td height='36.1px' width='36.1px'>&#10005</td>"
    }
    var upturn = document.getElementById("result").innerHTML

    if (upturn === "It's your turn") {

      upturn = "It's opponent's turn"

    }
    else {
      upturn = "It's your turn"
    }
    document.getElementById("result").innerHTML = upturn
  }

  if (task == "win") {
    document.getElementById("result").innerHTML = "You win";
    document.getElementById("result").style.backgroundColor = "#007bff"
  }
  if (task == "hitShip") {
    message = message.replace("[", "")
    message = message.replace("]", "")
    message = message.replace(/"/g, "")
    message = message.split(",")

    for (var i = 0; i < message.length; i++) {
      document.getElementById(message[i]).style.backgroundColor = "#008000"
    }

    count++;
    if (count == 5) {
      sendWin()
    }
  }
  if (task == "lose") {
    document.getElementById("result").innerHTML = "You lose. Better Luck next time"
    document.getElementById("result").style.backgroundColor = "#800000"
  }
  if (task == "disconnect") {
    document.getElementById("result").innerHTML = "Your Opponent has left"
    document.getElementById("result").style.backgroundColor = "	#8B4513"
  }
  sendId = function (id) {
    var upturn = document.getElementById("result").innerHTML
    if (upturn == "It's your turn") {
      this.socket.send(JSON.stringify({ task: 'shot', message: id, }))
      document.getElementById(id).removeAttribute("onclick")
    }
  }
  sendWin = function () {
    this.socket.send(JSON.stringify({ task: 'win' }))
  }
}



