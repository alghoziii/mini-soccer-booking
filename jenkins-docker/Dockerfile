FROM jenkins/jenkins:lts

USER root

# Install Docker CLI, Docker Buildx, Git, dan Curl di dalam container Jenkins
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        docker-cli \
        docker-buildx \
        git \
        curl && \
    rm -rf /var/lib/apt/lists/*

USER root
