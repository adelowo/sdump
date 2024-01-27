if [ -z "$1" ]; then
	echo "Please provide the image version"
	exit 10
fi

echo "Updating deployment to image with tag, $1"

kubectl set image deployment/sdump-api server=ghcr.io/adelowo/sdump:$1 --namespace sdump

echo "Checking rollout status... Hang on for around 5 seconds \n"
sleep .5

kubectl rollout status deployment.v1.apps/sdump-api --namespace sdump
