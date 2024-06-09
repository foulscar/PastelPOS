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

function ready() {
  console.log("Starting wsClient")
  const socket = io("http://" + host + "/socket.io/");
  
  socket.on("init", (payload) => {
    console.log("Received 'init': " + payload)

    const payloadJSON = JSON.parse(payload)
    wrapperTrackerGroupInProgress.innerHTML = "";
    wrapperTrackerGroupReady.innerHTML = "";

    for (orderID of payloadJSON.inProgress) {
      addOrder(orderID, false)
    }
    for (orderID of payloadJSON.ready) {
      addOrder(orderID, true)
    }
  })

  socket.on("addOrder", (orderID) => {
    console.log("Received 'addOrder': " + orderID)
    addOrder(orderID, false) 
  });

  socket.on("progressOrder", (orderID) => {
    console.log("Received 'progressOrder': " + orderID)
    const order = document.getElementById("order-" + orderID);
    wrapperTrackerGroupReady.prepend(order);
  });

  socket.on("removeOrder", (orderID) => {
    console.log("Received 'removeOrder': " + orderID)
    const order = document.getElementById("order-" + orderID);
    order.remove();
  });

  console.log("Joining Room: fohOrderTracker")
  socket.emit("join", "fohOrderTracker");
  
  setTimeout(function(){
    console.log("Requesting Initialization from wsServer")
    socket.emit("requestInit fohOrderTracker")
  }, 1000)
}

window.onload = ready()
