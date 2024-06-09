const host = window.location.host;
const wrapperTrackerGroupInProgress = document.getElementById("wrapper-tracker-group-inprogress");
const wrapperTrackerGroupReady = document.getElementById("wrapper-tracker-group-ready");

function addOrder(orderID, isReady) {
  const order = document.createElement("h3");
  order.id = "order-" + orderID;
  order.innerHTML = orderID;

  if (isReady) {
    wrapperTrackerGroupReady.append(order);   
  } else {
    wrapperTrackerGroupInProgress.append(order)
  }

}

function handleWSMessage(event) {
  let message = JSON.parse(event.data);
  switch(message.event) {

  }
}

function ready() {
  console.log("Starting wsClient")
  var socket = new WebSocket("ws://" + host + "/ws/fohOrderTracker");
  
  socket.onopen = function() {
    console.log("Connected to WebSocket Server");
  }
  
  socket.onmessage = function(event) {
    let message = JSON.parse(event.data);
    switch(message.event) {
    case "addOrder":
      console.log("Received 'addOrder': " + message.payload.orderID)
      addOrder(orderID, false)
    default:
      console.log("Received uknown event: " + message)
    }
  }
}

window.onload = ready()
