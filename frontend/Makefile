SHELL := /bin/bash
.PHONY: dev build clean run start

export NVM_DIR := $(HOME)/.nvm
NVM := source "$(NVM_DIR)/nvm.sh" && nvm use 24 &&

clean:
	rm -rf .next

build: clean
	$(NVM) npm run build

run: clean
	$(NVM) npm run dev&&
	osascript -e 'tell application "Google Chrome" to tell active tab of front window to execute javascript "caches.keys().then(ks => Promise.all(ks.map(k => caches.delete(k))))"' 2>/dev/null || true &&
	open -a "Google Chrome" "http://localhost:3000"

start: build
	$(NVM) npm run start

dev: run
