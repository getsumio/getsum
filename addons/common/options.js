// Copyright (c) 2013 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

var Config = {
	hostname : "",
	proxy : "",
	lib : "",
	timeout: 180
}

function loadRules() {
	var config = localStorage.config;
	try {
		config = JSON.parse(config);
	} catch (e) {
		config = Config;
	}
	Config = config;
	
	document.getElementById('hostname').value = Config.hostname;
	document.getElementById('proxy').value = Config.proxy;
	document.getElementById('timeout').value = Config.timeout;
	lib = document.getElementById('lib');
	setSelectedValue(lib, Config.lib);
	
	
}

function storeRules() {
	Config.hostname = document.getElementById('hostname').value;
	Config.proxy = document.getElementById('proxy').value;
	Config.timeout = document.getElementById('timeout').value;
	var lib = document.getElementById('lib');
	Config.lib = lib.options[lib.selectedIndex].text
	if(!Config.timeout || Config.timeout == null){
		Config.timeout = 180;
	}
	localStorage.config = JSON.stringify(Config);
	alert("Changes stored!");
}

window.onload = function() {
	loadRules();
	document.getElementById('save').onclick = function() {
		storeRules();
	};
}

function setSelectedValue(selectObj, valueToSet) {
    for (var i = 0; i < selectObj.options.length; i++) {
        if (selectObj.options[i].text== valueToSet) {
            selectObj.options[i].selected = true;
            return;
        }
    }
}
