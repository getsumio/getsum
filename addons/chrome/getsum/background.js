// Copyright (c) 2013 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

var config = {
	file : "",
	algorithm : [],
	checksum : "",
	timeout : 60,
	supplier : "go",
	all: false,
	remoteonly:false,
	localonly:false,
	timeout: 60,
	key: "",
	proxy: "",
	insecureskipverify: true
};
var ops = {
		hostname : "",
		proxy : "",
		lib : "",
		timeout: 180
	}
var processId = "";
const algs = [ "MD4", "MD5", "SHA1", "SHA224", "SHA256", "SHA384", "SHA512",
		"RMD160", "SHA3-224", "SHA3-256", "SHA3-384", "SHA3-512", "SHA512-224",
		"SHA512-256", "BLAKE2S256", "BLAKE2B256", "BLAKE2B384", "BLAKE2B512" ]

var ids = []
var options = ops;
if(localStorage.config && localStorage.config != null && localStorage.config != ""){
	options = JSON.parse(localStorage.config);
}

config.supplier = options.lib;
config.timeout = options.timeout;
config.proxy = options.proxy
console.log(options);
chrome.contextMenus.onClicked.addListener(function(info, tab) {

	if (tab) {
		if (info.menuItemId.startsWith("validate")) {
			config.file = info.linkUrl;
			config.algorithm = [];
			config.algorithm.push(info.menuItemId.split(':')[1]);
			config.checksum = "";
			console.log(config);
			chrome.tabs.sendMessage(tab.id, {
				"name" : "getsum",
				"config" : config,
				"ids" : tab.id
			});
		} else if (info.menuItemId.startsWith("calculate")) {
			config.file = info.linkUrl;
			config.algorithm = [];
			config.algorithm.push(info.menuItemId.split(':')[1]);
			config.checksum = "";
			console.log(config);
			chrome.tabs.sendMessage(tab.id, {
				"name" : "getsum",
				"config" : config,
				"ids" : tab.id
			});
		}
	}

});
chrome.contextMenus.removeAll(function() {
	for (let i = 0; i < algs.length; i++) {
		valId = 'validate:' + algs[i];
		if (!ids.includes(valId)) {
			chrome.contextMenus.create({
				id : valId,
				title : algs[i],
				contexts : [ 'link' ]
			});
			ids.push(valId);
		}
	}
});

function handleErrors(response) {
    if (!response.ok) {
        throw Error(response.statusText);
    }
    return response;
}


function postToServer(dataStr, id){
	processId = "";
	console.log(options);
	fetch(this.options.hostname, {
	    method: 'post',
	    cache: 'no-cache',
	    headers: {
	        'Accept': 'application/json',
	        'Content-Type': 'application/json'
	    },
	    body: dataStr,
	}) .then(handleErrors)
	    .then(result => {
	        // Here body is not ready yet, throw promise
	        if (!result.ok) {
	        	res.text().then(text => {throw Error(text);})
	        }
	        
	        result.json().then(function (json) {
	        	console.log("Post Response:");
	        	console.log(json);
	        	if(json.type == 4){
	        		processId = json.value;
		        	callBack(json, id);
		        	listen(id);
	        	}else{
	        		callBack(json, id);
	        	}
	        });
	    }).catch(error => {
	    	 console.log("Get error:");
			    errorStr = String(error);
		       	console.log(errorStr);
		       	if(errorStr.includes("Failed to fetch")){
		       		errorStr = "Can not reach server @ " + options.hostname ;
		       	}
		       	callBack( {
			    		type: 9,
			    		value: errorStr,
			    		checksum: ""
			    	},id);
	    });
}

function getFromServer(id){
	if(processId == null || processId == ""){
		return;
	}
	 return fetch(options.hostname + "/" + processId, {
	    method: 'get',
	    cache: 'no-cache',
	    headers: {
	        'Accept': 'application/json'
	    }
	}) .then(handleErrors)
	    .then(result => {
	        // Here body is not ready yet, throw promise
	    	if (!result.ok) {
	        	res.text().then(text => {throw Error(text);})
	        }
	        return result.json();
	    }).catch(error => {
	    	 console.log("Get error:");
			    errorStr = String(error);
		       	console.log(errorStr);
		       	if(errorStr.includes("Failed to fetch")){
		       		errorStr = "Can not reach server @ " + options.hostname ;
		       	}
		       	callBack( {
			    		type: 9,
			    		value: errorStr,
			    		checksum: ""
			    	},id);
	    });
}

function terminate(id){
	if(processId == null || processId == ""){
		return;
	}
	 return fetch(options.hostname + "/"+ processId, {
	    method: 'delete',
	    cache: 'no-cache',
	    headers: {
	        'Accept': 'application/json'
	    }
	}) .then(handleErrors)
	    .then(result => {
	        // Here body is not ready yet, throw promise
	    	if (!result.ok) {
	        	res.text().then(text => {throw Error(text);})
	        }
	        processId = "";
	    	callBack(  {
	        	type: 13,
	    		value: "Process terminated",
	    		checksum: ""
	        },id);
	    }).catch(error => {
		    console.log("Get error:");
		    errorStr = String(error);
		    processId = "";
	       	console.log(errorStr);
	       	if(errorStr.includes("Failed to fetch")){
	       		errorStr = "Can not reach server @ " + options.hostname ;
	       	}
	       	callBack( {
		    		type: 9,
		    		value: errorStr,
		    		checksum: ""
		    	},id);
	    });
}


function listen(id){
	setTimeout(function () {  
		getFromServer(id).then(result => {
			console.log("Get feedback");
			console.log(result);
			if(result){
				callBack(result, id);
				if(result.type < 7){
					listen(id);
				}
			}
		});
	   }, 1000)
	
}

function callBack(result, id){
	console.log('sending response');
	console.log(result);
	console.log(id);
	if(result.type == 4){
		result.value = "Waiting server response";
	}
	chrome.tabs.sendMessage(id, {
		name : 'status',
		dataStr : result,
		id : id
	});
}



chrome.runtime.onMessage.addListener(function (message, sender, sendResponse) {
	if(options.hostname == ""){
		callBack( {
    		type: 9,
    		value: "Please setup hostname from extension options",
    		checksum: ""
    	},id);
		return;
	}
	if(message.requestType == 'start'){
		console.log("Request received with id: " + message.ids);
		if(processId != ""){
			terminate(message.ids);
			processId = "";
		}
		postToServer(message.dataStr,message.ids);
		
	}else if(message.requestType == 'terminate'){
		console.log("Request received with id: " + message.ids);
		terminate(message.ids);
		
	}
});



