const LOCAL_SOCKET_URL = "http://localhost:7777";

chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
	if (changeInfo.status === "complete" && tab.url) {
		const url = new URL(tab.url);
		if (url.hostname.includes("geoguessr.com")) {
			sendURLToServer({ action: "open", url: tab.url });
		}
	}
});

chrome.tabs.onRemoved.addListener(() => {
	checkForGeoGuessrTabs();
});

function checkForGeoGuessrTabs() {
	chrome.tabs.query({}, (tabs) => {
		const geoGuessrTabs = tabs.filter((tab) => {
			try {
				const url = new URL(tab.url);
				return url.hostname.includes("geoguessr.com");
			} catch {
				return false;
			}
		});

		if (geoGuessrTabs.length === 0) {
			sendURLToServer({
				action: "close",
				url: "https://www.geoguessr.com",
			});
		}
	});
}

function sendURLToServer(data) {
	fetch(`${LOCAL_SOCKET_URL}`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(data),
	})
		.then((response) => {
			if (!response.ok) {
				console.error("Failed to send data to server");
			} else {
				console.log("Data sent to server successfully", data);
			}
		})
		.catch((error) => {
			console.error("Error sending data to server:", error);
		});
}
