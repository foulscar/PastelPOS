const host = window.location.host;
const wrapperMain = document.getElementById("wrapper-main");

var socket;

var menu;
fetch("/menu")
  .then(
    response => {
      response.json().then(function(data) {
        menu = data;
        console.log("Fetched Menu:");
        console.log(menu)
      });
    });

function handleOrderButtonClicked(orderID) {
  console.log("Sending removeOrder to wsServer");
  socket.emit("removeOrder fohOrderBagger", orderID);
}

function addOrder(orderID, order) {
  const orderDivMain = document.createElement("div");
  orderDivMain.className = "order";
  orderDivMain.id = "order-" + orderID;

  const orderTitle = document.createElement("h2");
  orderTitle.innerHTML = orderID;
  orderDivMain.append(orderTitle);
  
  const orderDivContent = document.createElement("div");
  orderDivContent.className = "order-content";
  orderDivMain.append(orderDivContent);

  const orderButton = document.createElement("button");
  orderButton.innerHTML = "✔️";
  orderButton.addEventListener("click", function() {
    handleOrderButtonClicked(orderID)
  }, {once: true});
  orderDivMain.append(orderButton);

  for (var index in order.meals) {
    const orderEntryGroup = document.createElement("div");
    orderEntryGroup.className = "order-entry group";
    orderDivContent.append(orderEntryGroup);
    
    const orderEntryGroupCount = document.createElement("h3");
    orderEntryGroupCount.innerHTML = "x" + order.meals[index].count;
    orderEntryGroup.append(orderEntryGroupCount);

    for (var itemType in order.meals[index].mealItems) {
      const itemID = order.meals[index].mealItems[itemType];
      if (itemType == "drink") {
        itemType = "drinks";
      }

      const orderEntry = document.createElement("div");
      orderEntry.className = "order-entry";
      orderDivContent.append(orderEntry);

      const orderEntryBracket = document.createElement("h3");
      orderEntryBracket.innerHTML = "└";
      orderEntry.append(orderEntryBracket);

      const orderEntryName = document.createElement("h3");
      orderEntryName.innerHTML = menu.itemsAvailable[itemType][itemID].icon + " ";
      orderEntryName.innerHTML += menu.itemsAvailable[itemType][itemID].name;
      orderEntry.append(orderEntryName)
    }
  }

  for (var singleItem in order.singleItems) {
    const itemID = order.singleItems[singleItem].item;
    var itemType = order.singleItems[singleItem].type;
    if (itemType == "drink") {
      itemType = "drinks";
    }

    const orderEntryGroup = document.createElement("div");
    orderEntryGroup.className = "order-entry group";
    orderDivContent.append(orderEntryGroup);

    const orderEntryCount = document.createElement("h3");
    orderEntryCount.innerHTML = "x" + order.singleItems[singleItem].count;
    orderEntryGroup.append(orderEntryCount);

    const orderEntryName = document.createElement("h3");
    orderEntryName.innerHTML = menu.itemsAvailable[itemType][itemID].icon + " ";
    orderEntryName.innerHTML += menu.itemsAvailable[itemType][itemID].name;
    orderEntryGroup.append(orderEntryName);
  }

  wrapperMain.append(orderDivMain);
}

function parseInitPayload(payloadJSON) {
  var orderArray = [];

  for (var orderID in payloadJSON.orders) {
    var time = payloadJSON.orders[orderID].time;
    orderArray.push({orderID: orderID, time: time});
  }

  orderArray.sort(function(a, b) {
    return a.time - b.time;
  });

  var sortedOrderIDs = orderArray.map(function(obj) {
    return obj.orderID;
  });

  for (var index in sortedOrderIDs) {
    addOrder(sortedOrderIDs[index], payloadJSON.orders[sortedOrderIDs[index]]);
  }
}

function ready() {
  console.log("Starting wsClient")
  socket = io("http://" + host + "/socket.io/");
  
  socket.on("init", (payload) => {
    console.log("Received 'init': " + payload)

    const payloadJSON = JSON.parse(payload)
    wrapperMain.innerHTML = "";

    parseInitPayload(payloadJSON);
  })

  socket.on("addOrder", (payload) => {
    console.log("Received 'addOrder': " + payload)

    const payloadJSON = JSON.parse(payload);
    addOrder(payloadJSON.orderID, payloadJSON.order) 
  });

  socket.on("removeOrder", (orderID) => {
    console.log("Received 'removeOrder': " + orderID)
    const order = document.getElementById("order-" + orderID);
    order.remove();
  });

  console.log("Joining Room: fohOrderBagger")
  socket.emit("join", "fohOrderBagger");
  
  setTimeout(function(){
    console.log("Requesting Initialization from wsServer")
    socket.emit("requestInit fohOrderBagger")
  }, 1000)
}

window.onload = ready()
