// Copyright (c) 2013 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//logger
//https://stackoverflow.com/questions/18410119/is-it-possible-to-bind-a-date-time-to-a-console-log
var DEBUG = (function(){
    var timestamp = function(){};
    timestamp.toString = function(){
        return "[GetSum - " + (new Date).toLocaleTimeString() + "]:";    
    };

    return {
        log: console.log.bind(console, '%s', timestamp)
    }
})();

//base configration will be passed arround
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
// default ops
var ops = {
	hostname : "",
	proxy : "",
	lib : "",
	timeout: 180
}
// current process id received from server on post
var processId = "";
// algos
const algs = [ "MD4", "MD5", "SHA1", "SHA224", "SHA256", "SHA384", "SHA512",
		"RMD160", "SHA3-224", "SHA3-256", "SHA3-384", "SHA3-512", "SHA512-224",
		"SHA512-256", "BLAKE2S256", "BLAKE2B256", "BLAKE2B384", "BLAKE2B512" ]
// tab id used by extension
var ids = [];
var options = ops;
setDefaults();
DEBUG.log("Options parsed: ", options);
/*
 * Reads options Tries to parse if failure then uses default one
 */
function setDefaults(){
	if(localStorage.config && localStorage.config != null && localStorage.config != ""){
		try{
			options = JSON.parse(localStorage.config);
		}catch (e) {
			DEBUG.log("Can not parse options!",e);
		}
	}

	config.supplier = options.lib;
	config.timeout = options.timeout;
	config.proxy = options.proxy
}

/*
 * Add menu items
 */
chrome.contextMenus.removeAll(function() {
	for (let i = 0; i < algs.length; i++) {
		// make sure no duplicate
		valId = 'getsum:' + algs[i];
		DEBUG.log('Generating menu item: ' + valId);
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
/*
 * Add Context menu listener
 */
chrome.contextMenus.onClicked.addListener(function(info, tab) {
	if (tab) {
		if (info.menuItemId.startsWith("getsum")) {
			
			
			config.file = info.linkUrl;
			config.algorithm = [];
			config.algorithm.push(info.menuItemId.split(':')[1]);
			config.checksum = "";
			DEBUG.log('Menu click: ' + info.menuItemId,config);
			// send message on click to start process
			// content.js will read selected text if exist
			// and it will callback
			chrome.tabs.sendMessage(tab.id, {
				"name" : "getsum",
				"config" : config,
				"ids" : tab.id
			});
		} 
	}

});



/*
 * Checks response is not ok and throws response text
 */
function handleErrors(response) {
    if (!response.ok) {
    	response.text().then(text => {throw Error(text);});
    }
    return response;
}

/*
 * Validates if hostname is set otherwise callbacks with error code
 */
function checkOptions(id){
	if(!options || options.hostname == ""){
		DEBUG.log('invalid option settings')
		callBack( {
    		type: 9, // ERROR
    		value: 'Please first visit extension options and set hostname',
    		checksum: ""
    	},id); // tab id
		return false;
	}
	return true;
}
/*
 * Starts the process by posting config to server
 */
function postToServer(dataStr, id){
	processId = "";
	
	// validate options
	if(!checkOptions(id)){
		return;
	}
	
	// post
	fetch(this.options.hostname, {
	    method: 'post',
	    cache: 'no-cache',
	    headers: {
	        'Accept': 'application/json',
	        'Content-Type': 'application/json'
	    },
	    body: dataStr,
	}) 
	// check response code
	.then(handleErrors)
	// parse response and call back
	.then(result => {
        result.json().then(function (json) {
        	// 4 == STARTED
        	DEBUG.log('Post response received', json);
        	if(json.type == 4){
        		// process id needed to watch status
        		processId = json.value;
	        	callBack(json, id);
	        	listen(id);
        	}else{
        		// couldnt start return response
        		callBack(json, id);
        	}
        });
    })
    // catch error
	.catch(error => {
		errorStr = getErrorString(error);
       	callBackError(errorStr, id);
	});
}

/*
 * If processId present retrieves status from server
 */
function getFromServer(id){
	// check options
	if(!checkOptions(id) || processId == null || processId == ""){
		return;
	}
	// get
	return fetch(options.hostname + "/" + processId, {
	    method: 'get',
	    cache: 'no-cache',
	}) 
	// check response code
	.then(handleErrors)
	// return result
	.then(result => {
		DEBUG.log('Get response received', result);
        return result.json();
    }).catch(error => {
    		errorStr = getErrorString(error);
	       	callBackError(errorStr, id);
    });
}

function terminate(id){
	// check options
	if(!checkOptions(id) || processId == null || processId == ""){
		return;
	}
	// call delete
	 return fetch(options.hostname + "/"+ processId, {
	    method: 'delete',
	    cache: 'no-cache',
	}) 
	// check errors
	.then(handleErrors)
	// call back if ok
    .then(result => {
    	DEBUG.log('Terminate response received', result);
        callBackFromParameters(id,13,"Process Terminated","");
    }).catch(error => {
    	errorStr = getErrorString(error);
       	callBackError(errorStr, id);
    });
	// reset response id
	processId = "";
}

/*
 * Watch status until we got error
 */
function listen(id){
	DEBUG.log("Starting listening");
	setTimeout(function () {  
		// get
		getFromServer(id)
		// parse return
		.then(result => {
			// we got result
			if(result){
				// notify content js
				callBack(result, id);
				// still running keep watching
				if(result.type < 7){
					listen(id);
				}
			}
		});
	   }, 1000);
	
}

/*
 * error to string
 */
function getErrorString(error){
	if(!error){
		return "";
	}
	
	errorStr = String(error);
   	if(errorStr.includes("Failed to fetch")){
   		errorStr = "Can not reach server @ " + options.hostname ;
   	}
   	DEBUG.log(errorStr);
   	return errorStr;
}

/*
 * Notifies content.js
 */
function callBack(result, id){
	//4 == STARTED, add meaningful text
	if(result.type == 4){
		result.value = "Waiting server response";
	}
	
	//callback
	DEBUG.log("informing back content.js", result);
	chrome.tabs.sendMessage(id, {
		name : 'status',
		dataStr : result,
		id : id
	});
}

/*
 * utility to callback with json
 */
function callBackFromParameters(id,type,value,checksum){
	callBack( {
		type: type,
		value: value,
		checksum: checksum
	},id);
}

/*
 * utility to callback with error
 */
function callBackError(errorStr, id){
	callBack( {
		type: 9,
		value: errorStr,
		checksum: ""
	},id);
}


/*
 * When content.js sends message do related http request
 */
chrome.runtime.onMessage.addListener(function (message, sender, sendResponse) {
	DEBUG.log('message received', message);
	if(message.requestType == 'start'){
		if(processId != ""){
			DEBUG.log('There is an existing process cleaning|terminating ' + processId);
			terminate(message.ids);
		}
		postToServer(message.dataStr,message.ids);
		
	}else if(message.requestType == 'terminate'){
		terminate(message.ids);
	}
});



