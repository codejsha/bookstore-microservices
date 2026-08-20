SERVICE_NAME := mobile

.PHONY: none
none:
	@echo "No target specified."

.PHONY: dev
dev:
	npx expo start

.PHONY: ios
ios:
	npx expo start --ios

.PHONY: android
android:
	npx expo start --android

.PHONY: lint
lint:
	npx biome lint ./src ./app

.PHONY: check
check:
	npx biome check --write ./src ./app
