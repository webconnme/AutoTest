
//setInterval(function() {
//  $("#log-box").append("[" + Date().toLocaleString() + "] abc\n");
//  $("#log-box")[0].scrollTop = $("#log-box")[0].scrollHeight;
//}, 5000);
var socket = io();


socket.on('message', function(mesg) {
//    console.log(mesg);

    $("#log-box").append(mesg);
    $("#log-box")[0].scrollTop = $("#log-box")[0].scrollHeight;

    // 메세지 확인
    // start
    if (mesg.match("start")) {
        console.log("start");
//        $('.progress-bar').css('width: 100%');
        $('.progress-bar').show();
    }
    // total
    if (mesg.match("total")) {
        console.log("total");
    }
    // end
    if (mesg.match("end")) {
        console.log("end");

        $('.progress-bar').removeClass('active');
        $('.progress-bar').hide();
        alert('종료 하였습니다.!!!');
    }

});

$('#btn_start').click(function () {
    $("#log-box").empty();
//    console.log('start>>>>>>>>>>>');
//    $('#btn_start').addClass("disabled");
//    $('#btn_stop').removeClass("disabled");

    socket.emit('start');

});

$('#btn_stop').click(function () {
//    console.log('stop>>>>>>>>>>>');
//    $('#btn_stop').addClass("disabled");
//    $('#btn_start').removeClass("disabled");
    socket.emit('stop');

});