var term = new Terminal();
term.open(document.getElementById('terminal'));

var uuid = window.location.pathname.split("/session/")[1]


const conn = new WebSocket("ws://" + document.location.host + "/session/" + uuid + "/ws");
conn.onclose = function (evt) {
    var item = document.createElement("div");
    item.innerHTML = "<b>Connection closed.</b>";
    console.log(item);
};
conn.onmessage = function (evt) {
        console.log(`Should write ${evt.data} to the terminal\n`);
        term.write(evt.data, () => {
            console.log("Finished write..\n");
        });
};
conn.onopen = function (evt) {
    console.log("Websocket connection opened\n");
};

term.onData((data) => {
    console.log(`Sending: ${data}\n`)
    conn.send(data)
});


// dragable
dragElement(document.getElementById("draggable"));

function dragElement(elmnt) {
  var pos1 = 0, pos2 = 0, pos3 = 0, pos4 = 0;
    elmnt.onmousedown = dragMouseDown;

  function dragMouseDown(e) {
    e = e || window.event;
    e.preventDefault();
    // get the mouse cursor position at startup:
    pos3 = e.clientX;
    pos4 = e.clientY;
    document.onmouseup = closeDragElement;
    // call a function whenever the cursor moves:
    document.onmousemove = elementDrag;
  }

  function elementDrag(e) {
    e = e || window.event;
    e.preventDefault();
    // calculate the new cursor position:
    pos1 = pos3 - e.clientX;
    pos2 = pos4 - e.clientY;
    pos3 = e.clientX;
    pos4 = e.clientY;
    // set the element's new position:
    elmnt.style.top = (elmnt.offsetTop - pos2) + "px";
    elmnt.style.left = (elmnt.offsetLeft - pos1) + "px";
  }

  function closeDragElement() {
    // stop moving when mouse button is released:
    document.onmouseup = null;
    document.onmousemove = null;
  }
}