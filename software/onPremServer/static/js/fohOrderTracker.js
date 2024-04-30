const host = window.location.host;
const wrapperTrackerGroupInProgress = document.getElementById("wrapper-tracker-group-inprogress");
const wrapperTrackerGroupReady = document.getElementById("wrapper-tracker-group-ready");

const socket = io("http://" + host + "/socket.io/");

socket.emit("join", "fohOrderTracker");

socket.on("addOrder", (orderNumber) => {
 const order = document.createElement("h3");
 order.id = "order-" + orderNumber;
 order.innerHTML = orderNumber;

 console.log(orderNumber)
 wrapperTrackerGroupInProgress.append(order);
});

socket.on("progressOrder", (orderNumber) => {
  const order = document.getElementById("order-" + orderNumber);
  wrapperTrackerGroupReady.prepend(order);
});

socket.on("removeOrder", (orderNumber) => {
  const order = document.getElementById("order-" + orderNumber);
  order.remove();
});
