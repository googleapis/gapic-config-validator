
build:
	gcloud builds submit --tag gcr.io/config-checker/config-checker

deploy:
	gcloud beta run deploy --image gcr.io/config-checker/config-checker --platform managed
