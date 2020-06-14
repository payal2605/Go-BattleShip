var table1 = document.getElementById("table1");
for (var i = 0; i < 10; i++) {
    var row = document.createElement("tr");
    for (var j = 0; j < 10; j++) {
        var col = document.createElement("td");
        col.setAttribute("id", "" + i + j + "me");
        col.style.border = "1px solid green"
        col.style.height = "36.1px"
        col.style.width = "36.1px"
        row.appendChild(col);
    }
    table1.appendChild(row);
}

var table = document.getElementById("table2");
for (var i = 0; i < 10; i++) {
    var row = document.createElement("tr");
    for (var j = 0; j < 10; j++) {
        var col = document.createElement("td");
        col.setAttribute("id", "" + i + j);
        col.style.border = "1px solid green"
        col.style.height = "36.1px"
        col.style.width = "36.1px"
        col.setAttribute("onclick", 'sendId(this.id)')
        row.appendChild(col);
    }
    table.appendChild(row);
}





