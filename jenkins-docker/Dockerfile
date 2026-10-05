FROM jenkins/jenkins:lts

USER root

# Install Docker CLI dan Git di dalam container Jenkins
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        docker.io \
        git \
        curl && \
    rm -rf /var/lib/apt/lists/*

USER root
