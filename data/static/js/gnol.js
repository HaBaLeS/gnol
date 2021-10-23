$( "#enableWebAuthn" ).on( "click", function( event ) {
    $( "#enableWebAuthn" ).addClass("is-active");
    $( "#enableClassicAuth" ).removeClass("is-active");
    $("#classicAuthForm").hide();
    $("#classicWebAuthnForm").show();
});

$( "#enableClassicAuth" ).on( "click", function( event ) {
    $( "#enableClassicAuth" ).addClass("is-active");
    $( "#enableWebAuthn" ).removeClass("is-active");
    $("#classicAuthForm").show();
    $("#classicWebAuthnForm").hide();
});

