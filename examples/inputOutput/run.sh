#!/usr/bin/env bash
TIMEOUT=1m
PATHS=$1
IMAGE=$2
FILES=$3
CONTAINER_ID=

term_handler() {
    # This catch both temporary and persistent containers
    if docker inspect -f {{.State.Running}} $CONTAINER_ID > /dev/null; then
        docker rm -f $CONTAINER_ID > /dev/null
    fi

    exit 0
}

trap 'term_handler' SIGINT SIGTERM

# Warning: only --rm can save you from losing control over the spawned container if the host OS
# will terminated suddenly
# Time for pulling the image will not included to timeout
CONTAINER_ID=$(docker run -d -it --rm $PATHS $IMAGE $FILES)

echo "Container $CONTAINER_ID started in background for $TIMEOUT"
# Sleep in the foreground will block traps thus doing this in the background.
sleep $TIMEOUT &
SLEEP_PID=$!
# Wait until sleep is finished
wait $SLEEP_PID

docker rm -f $CONTAINER_ID > /dev/null
