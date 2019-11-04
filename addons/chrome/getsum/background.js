// Copyright (c) 2013 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

var config = {
	file : "",
	algorithm : "",
	checksum : "",
	timeout : 60,
	supplier : "go"
};

const algs = [ "MD4", "MD5", "SHA1", "SHA224", "SHA256", "SHA384", "SHA512",
		"RMD160", "SHA3-224", "SHA3-256", "SHA3-384", "SHA3-512", "SHA512-224",
		"SHA512-256", "BLAKE2s256", "BLAKE2b256", "BLAKE2b384", "BLAKE2b512" ]

var ids = []

chrome.contextMenus.onClicked.addListener(function(info, tab) {

	if (tab) {
		if (info.menuItemId.startsWith("validate")) {
			config.file = info.linkUrl;
			config.algorithm = info.menuItemId.split(':')[1];
			config.checksum = "Y";
			console.log(config);
			chrome.tabs.sendMessage(tab.id, {
				"name" : "getsum",
				"config" : config
			});
		} else if (info.menuItemId.startsWith("calculate")) {
			config.file = info.linkUrl;
			config.algorithm = info.menuItemId.split(':')[1];
			config.checksum = "";
			console.log(config);
			chrome.tabs.sendMessage(tab.id, {
				"name" : "getsum",
				"config" : config
			});
		}
	}

});
chrome.contextMenus.removeAll(function() {
	if (!ids.includes('calculate')) {
		chrome.contextMenus.create({
			id : 'calculate',
			title : chrome.i18n.getMessage('openContextMenuCalculate'),
			contexts : [ 'link' ],
		});
		ids.push('calculate');
	}
	if (!ids.includes('validate')) {
		chrome.contextMenus.create({
			id : 'validate',
			title : chrome.i18n.getMessage('openContextMenuValidate'),
			contexts : [ 'link' ],
		});
		ids.push('validate');
	}

	for (let i = 0; i < algs.length; i++) {
		valId = 'validate:' + algs[i];
		calId = 'calculate:' + algs[i];
		if (!ids.includes(valId)) {
			chrome.contextMenus.create({
				id : valId,
				title : algs[i],
				contexts : [ 'link' ],
				parentId : 'validate'
			});
			ids.push(valId);
		}
		if (!ids.includes(calId)) {
			chrome.contextMenus.create({
				id : calId,
				title : algs[i],
				contexts : [ 'link' ],
				parentId : 'calculate'
			});
			ids.push(calId);
		}
	}
});
