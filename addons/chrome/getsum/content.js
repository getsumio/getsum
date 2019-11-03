var $j = jQuery.noConflict();
chrome.runtime.onMessage
		.addListener(function(request, sender, sendResponse) {
			if (request.name == "getsum") {
				config = request.config;
				if (config.checksum == "Y") {
					config.checksum = window.getSelection().toString();
				}else {
					config.checksum = "";
				}
				var cfg = JSON.stringify(config);
				
				$j('#getsumModalDialog').toggle("slide:right");
				console.log('done');
			}
		});

$j.get(chrome.extension.getURL('/modal.html'), function(data) {
    // Or if you're using jQuery 1.8+:
	$j($j.parseHTML(data)).appendTo('body');
	/*$j('#getsumModalDialog').dialog({
		'autoOpen' : false
	});*/
});