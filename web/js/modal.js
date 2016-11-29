
//Verify Modal
$('#verifyButton').click(function() {
	$.get("http://localhost:12345/verify", function(data, status){
		$(data).replaceAll("#verify-modal")

	});
});

//Refresh Button
$('#refreshButton').click(function() {
	$.get("http://localhost:12345/refresh", function(data, status){
		$(data).replaceAll("#LoginInformation")

	});
});

//Message Modal Button
$('#send-message-button').click(function() {
	$.post("http://localhost:12345/message", $( "#testform" ).serialize() , function(data, status){
		$(data).replaceAll("#LoginInformation")

	});
});