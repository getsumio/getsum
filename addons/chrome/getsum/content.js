var $j = jQuery.noConflict();
var pending = [];
const statuses = [ "PREPARED", "ALLOCATED", "DOWNLOAD", "FETCHED", "STARTED",
		"RESUMING", "RUNNING", "COMPLETED", "SUSPENDED", "ERROR", "TIMEDOUT",
		"MISMATCH", "VALIDATED", "TERMINATED",

];
var validationChecksum = "";
var tabId = "";
chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
	console.log('message received!');
	console.log(request);
	if (request.name == "getsum") {
		config = request.config;
		if (config.checksum == "Y") {
			config.checksum = getSelection();
		} else {
			config.checksum = "";
		}
		validationChecksum = config.checksum;
		var cfg = JSON.stringify(config);
		reset();
		$j('#getsumModalDialog').show(200);
		tabId = request.ids;
		start(config, request.ids);
	} else if (request.name == 'status') {
		console.log(request.dataStr);
		status = validate(request);
		$j('#getsumLabel').text(status);
		$j('#getsumValue').text(request.dataStr.value);
		$j('#getsumChecksum').text(request.dataStr.checksum);
	}
});

function getSelection(){
	selection = window.getSelection().toString();
	selection = selection.trim();
	if(selection.includes(" ") || !isBase64(selection)){
		return "";
	}
	return selection;
}

//https://stackoverflow.com/questions/7860392/determine-if-string-is-in-base64-using-javascript
function isBase64(str) {
    if (str ==='' || str.trim() ===''){ return false; }
    try {
        return btoa(atob(str)) == str;
    } catch (err) {
        return false;
    }
}

function validate(request) {
	if (request.dataStr.type >= 7) {
		switch (request.dataStr.type) {
		case 7:
			return validateStatus(request);
		case 12:
			return validateStatus(request);
		default:
			return "\u24e7  " + statuses[request.dataStr.type];
		}
	}
	return statuses[request.dataStr.type];
}

function validateStatus(request) {
	if (validationChecksum != "") {
		if (request.dataStr.checksum != validationChecksum) {
			return '\u24e7  MISMATCH';
		} else {
			return '\u2713  VALIDATED';
		}
	} else {
		return '\u2713  ' + statuses[request.dataStr.type];
	}
}

$j.get(chrome.extension.getURL('/modal.html'), function(data) {
	$j($j.parseHTML(data)).appendTo('body');
	$j('#closeGetsum').click(function() {
		reset();
		chrome.runtime.sendMessage({
			requestType : 'terminate',
			ids : tabId
		});
		toggleDialog();
	});
	$j('#getsumModalDialog').hide();
});
function toggleDialog() {
	$j('#getsumModalDialog').toggle("slide:right");
}
function reset() {
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
