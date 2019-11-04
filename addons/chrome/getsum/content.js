var $j = jQuery.noConflict();
var processId;
chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
	if (request.name == "getsum") {
		config = request.config;
		if (config.checksum == "Y") {
			config.checksum = window.getSelection().toString();
		} else {
			config.checksum = "";
		}
		var cfg = JSON.stringify(config);

		toggleDialog();
		start(config);
		console.log('done');
	}
});

$j.get(chrome.extension.getURL('/modal.html'), function(data) {
	$j($j.parseHTML(data)).appendTo('body');
	$j('#closeGetsum').click(function() {
		toggleDialog()
	});
	$j('#getsumModalDialog').hide();
});
function toggleDialog() {
	$j('#getsumModalDialog').toggle("slide:right");
}
function start(configData){
	data = JSON.stringify(configData);
	console.log(data);
	$j.post('https://example.com:8088',   // url
		       { config: data }, // data to be submit
		       function(data, status, jqXHR) {// success callback
		                $j('#getsumModalDialog').append('status: ' + status + ', data: ' + data);
		        })
}
