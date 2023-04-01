///////                    Websocket Online                       //////
const connectionOnline = new WebSocket("ws://localhost:8080/online");
const connectionChat = new WebSocket("ws://localhost:8080/chat");

connectionOnline.onopen = (event) => {
    console.log("WebSocketOnline is open now.", event);
};
connectionOnline.onclose = (event) => {
    console.log("WebSocketOnline is closed now.", event);
    chat.innerHTML = event.data;
};

connectionOnline.onerror = (event) => {
    console.error("WebSocketOnline error observed:", event);
};


connectionOnline.onmessage = (event) => {
    // append received message from the server to the DOM element
    const chat = document.querySelector("#online");
    chat.innerHTML = `Online: ${event.data}`;
};

connectionOnline.addEventListener("message", () => {
});

///////                    Websocket Chat                       //////


const button = document.querySelector("#send");

connectionChat.onopen = (event) => {
    console.log("WebSocketChat is open now.");
};

connectionChat.onclose = (event) => {
    console.log("WebSocketChat is closed now.");
};

connectionChat.onerror = (event) => {
    console.error("WebSocketChat error observed:", event);
};

connectionChat.onmessage = (event) => {
    // append received message from the server to the DOM element
    const chat = document.querySelector("#chat");

    let jsonArray = JSON.parse(event.data);
    let tableHtml = ""
    for (let i = 0; i < jsonArray.length; i++) {
        let obj = jsonArray[i];
        tableHtml += `<p>User: ${obj.user_id} message : ${obj.message}</p>`;
    }
    chat.innerHTML += tableHtml;
    if (typeof jsonArray === "object") {
        chat.innerHTML += `<p>User: ${jsonArray.user_id} message : ${jsonArray.message}</p>`;
    }
};

button.addEventListener("click", () => {
    const name = document.querySelector("#name");
    const message = document.querySelector("#message");

    // Send composed message to the server
    msg = {
        type: "message",
        message: message.value,
        user_id: name.value,
        date: Date.now(),
    };

    connectionChat.send(JSON.stringify(msg));

    // clear input fields
    name.value = "";
    message.value = "";
});