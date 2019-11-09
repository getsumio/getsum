var $j = jQuery.noConflict();
var pending = [];
const statuses = [ "PREPARED", "ALLOCATED", "DOWNLOAD", "FETCHED", "STARTED",
		"RESUMING", "RUNNING", "COMPLETED", "SUSPENDED", "ERROR", "TIMEDOUT",
		"MISMATCH", "VALIDATED", "TERMINATED",

];
// user selected text
var validationChecksum = "";
// current tab id
var tabId = "";

var DEBUG = (function(){
    var timestamp = function(){};
    timestamp.toString = function(){
        return "[GetSum - " + (new Date).toLocaleTimeString() + "]:";    
    };

    return {
        log: console.log.bind(console, '%s', timestamp)
    }
})();
// listen messages from backend.js
chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
	// message received
	DEBUG.log("A message received", request);
	if (request.name == "getsum") {
		config = request.config;
		config.checksum = getSelectedText();
		validationChecksum = config.checksum;
		DEBUG.log("Validation checksum: "+ validationChecksum);
		var cfg = JSON.stringify(config);
		reset();
		$j('#getsumModalDialog').show(200);
		tabId = request.ids;
		start(config, request.ids);
	} else if (request.name == 'status') {
		status = validate(request);
		$j('#getsumLabel').text(status);
		$j('#getsumValue').text(request.dataStr.value);
		$j('#getsumChecksum').text(request.dataStr.checksum);
	}
});

/*
 * If user selected text checks if its a valid base64 value
 */
function getSelectedText() {
	
	selection = window.getSelection().toString();
	DEBUG.log("Parsing selected text:", selection);
	selection = selection.trim();
	if (selection.includes(" ") || !isBase64(selection)) {
		DEBUG.log("Not a valid base64!:", selection);
		return "";
	}
	return selection;
}

// https://stackoverflow.com/questions/7860392/determine-if-string-is-in-base64-using-javascript
function isBase64(str) {
	if (str === '' || str.trim() === '') {
		return false;
	}
	try {
		return btoa(atob(str)) == str;
	} catch (err) {
		return false;
	}
}

/*
 * if there is validation updates status as MISMATCH/VALIDATED adds also
 * error/checked chars
 */
function validate(request) {
	DEBUG.log("validating request:", request);
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

/*
 * if there is validation updates status as MISMATCH/VALIDATED
 */
function validateStatus(request) {
	if (validationChecksum != "") {
		if (request.dataStr.checksum != validationChecksum) {
			DEBUG.log("Request mismatch:", validationChecksum);
			return '\u24e7  MISMATCH';
		} else {
			DEBUG.log("Request validated:", validationChecksum);
			return '\u2713  VALIDATED';
		}
	} else {
		return '\u2713  ' + statuses[request.dataStr.type];
	}
}

/*
 * Append popup to bady
 */
$j.get(chrome.extension.getURL('/modal.html'), function(data) {
	DEBUG.log("Modal added to content");
	$j($j.parseHTML(data)).appendTo('body');
	// add close listener, make sure process terminated
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

// toggle
function toggleDialog() {
	$j('#getsumModalDialog').toggle("slide:right");
}

// reset all values
function reset() {
	DEBUG.log("Resetting");
	$j('#getsumLabel').text("");
	$j('#getsumValue').text("");
	$j('#getsumChecksum').text("");
}

// start calculation
function start(configData, id) {
	DEBUG.log("Calling backend to start process:", configData);
	data = JSON.stringify(configData);
	chrome.runtime.sendMessage({
		requestType : 'start',
		dataStr : data,
		ids : id
	});
}
