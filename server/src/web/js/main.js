
//setInterval(function() {
//  $("#log-box").append("[" + Date().toLocaleString() + "] abc\n");
//  $("#log-box")[0].scrollTop = $("#log-box")[0].scrollHeight;
//}, 5000);
var socket = io();


socket.on('message', function(mesg) {
//    console.log(mesg);
    $("#log-box").append("[" +mesg + "]\n");
    $("#log-box")[0].scrollTop = $("#log-box")[0].scrollHeight;
});

$('#btn_start').click(function () {

//    console.log('start>>>>>>>>>>>');
    $('#btn_start').addClass("disabled");
    $('#btn_stop').removeClass("disabled");

    socket.emit('start');

});

$('#btn_stop').click(function () {
//    console.log('stop>>>>>>>>>>>');
    $('#btn_stop').addClass("disabled");
    $('#btn_start').removeClass("disabled");
    socket.emit('stop');

});