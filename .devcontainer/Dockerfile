FROM mcr.microsoft.com/devcontainers/go:1-bookworm

# Install Python for pre-commit (pre-commit requires Python)
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
  && apt-get -y install --no-install-recommends python3 python3-pip python3-venv

# Clean up
RUN apt-get autoremove -y && apt-get clean -y
