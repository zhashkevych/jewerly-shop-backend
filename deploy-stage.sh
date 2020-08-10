if [ ! "$(git status | grep 'develop')" ]; then
    git checkout develop
fi

git pull origin develop

export HOST=stage
docker image build -t jewerly-api:0.1 -f ./deploy/Dockerfile .

if [ "$(docker ps -q -f name=jewerly-api-stage)" ]; then
    if [ ! "$(docker ps -aq -f status=exited -f name=jewerly-api-stage)" ]; then
        docker stop jewerly-api-stage
    fi
fi

docker run -e HOST -e ACCESS_KEY -e SECRET_KEY -e POSTGRES_PASSWORD --rm -d --publish 8001:8000 --name jewerly-api-stage --link=jewerly-db:db jewerly-api:0.1