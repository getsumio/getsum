var $j = jQuery.noConflict();
var pending = [];
const statuses = [ "PREPARED", "ALLOCATED", "DOWNLOAD", "FETCHED", "STARTED",
		"RESUMING", "RUNNING", "COMPLETED", "SUSPENDED", "ERROR", "TIMEDOUT",
		"MISMATCH", "VALIDATED", "TERMINATED",

];
var validationChecksum = "";
chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
	console.log('message received!');
	console.log(request);
	if (request.name == "getsum") {
		config = request.config;
		if (config.checksum == "Y") {
			config.checksum = window.getSelection().toString();
		} else {
			config.checksum = "";
		}
		validationChecksum = config.checksum;
		var cfg = JSON.stringify(config);
		reset();
		$j('#getsumModalDialog').show("slide:left");
		start(config, request.ids);
	} else if (request.name == 'status') {
		console.log(request.dataStr);
		status = statuses[request.dataStr.type];
		if (request.dataStr.type >= 7) {
			char = "";
			switch (request.dataStr.type) {
			case 7:
				char = "\u2713";
				if(validationChecksum != "" ){
					if(request.dataStr.checksum != validationChecksum){
						char = "\u24e7";
						status = 'MISMATCH';
					}else{
						status = 'VALIDATED';
					}
					break;
				}
				break;
			case 12:
				char = "\u2713";
				if(validationChecksum != "" ){
					if(request.dataStr.checksum != validationChecksum){
						char = "\u24e7";
						status = 'MISMATCH';
					}else{
						status = 'VALIDATED';
					}
					break;
				}
				break;
			default:
				char = "\u24e7";
			}
			status = char + " " + status;
		}
		$j('#getsumLabel').text(status);
		$j('#getsumValue').text(request.dataStr.value);
		$j('#getsumChecksum').text(request.dataStr.checksum);
	}
});

$j.get(chrome.extension.getURL('/modal.html'), function(data) {
	$j($j.parseHTML(data)).appendTo('body');
	$j('#closeGetsum').click(function() {
		reset();
		toggleDialog()
	});
	$j('#getsumModalDialog').hide();
});
function toggleDialog() {
	$j('#getsumModalDialog').toggle("slide:right");
}
function reset(){
	$j('#getsumLabel').text("");
	$j('#getsumValue').text("");
	$j('#getsumChecksum').text("");
}
function start(configData, id) {
	data = JSON.stringify(configData);
	console.log('starting calculation')
	console.log(data);

	chrome.runtime.sendMessage({
		requestType : 'start',
		dataStr : data,
		ids : id
	});
}
