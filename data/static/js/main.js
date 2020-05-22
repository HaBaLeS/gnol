//Color wheel is: https://iro.js.org/


$('#testbutton').click(function() {
    $('#testbutton').toggleClass('is-link');
    if($('#testbutton').hasClass('is-link')){
        $('#testbutton').html("ON")
    } else {
        $('#testbutton').html("OFF")
    }
   /* $.get( "sendtest", function( data ) {
        $( ".result" ).html( data );

    });*/

    $.ajax({
        url: "sendtest",
        type: "get", //send it through get method
        data: {
            ajaxid: 4,
            UserID: "uis",
            EmailAddress: "mailla"
        },
        success: function(response) {
            //Do Something
            //alert(response.room)
        },
        error: function(xhr) {
            //Do Something to handle error
            alert(xhr)
        }
    });
});




var colorWheel = new iro.ColorPicker("#colorWheelDemo", {
    width: 280,
    color: "rgb(255, 0, 0)",
    borderWidth: 1,
    borderColor: "#fff",
});


var colorWheel2 = new iro.ColorPicker("#colorWheelDemo2", {
    width: 280,
    color: "rgb(255, 0, 0)",
    borderWidth: 1,
    borderColor: "#fff",
});
//BG #01004a


$(".card-header").click(function () {

    $header = $(this);
    //getting the next element
    $content = $header.next();
    //open up the content needed - toggle the slide- if visible, slide up, if not slidedown.
    $content.slideToggle(500, function () {
        //execute this after slideToggle is done
        //change text of header based on visibility of content div
        $('#toggle_x').toggleClass('fa-angle-down');
        $('#toggle_x').toggleClass('fa-angle-up');
       /* $header.text(function () {
            //change text based on condition

            return $content.is(":visible") ? "Collapse" : "Expand";
        });*/
    });

});


var lastSend = 0;
// listen to a color picker's color:change event
// color:change callbacks receive the current color
colorWheel.on('color:change', function(color) {
    // log the current color as a HEX string
    //console.log(color.hexString);

    //FIXME debounce
    now = new Date().getTime();
    dt = now - lastSend;
    if(dt >200) {
        lastSend = now;
        $.get("sendtest?" + color.hexString, function (data) {
            //$( ".result" ).html( data );
        });
    }


});