FROM phusion/baseimage:latest
MAINTAINER Sergey Arkhipov <serge@aerialsounds.org>

ENV DEBIAN_FRONTEND noninteractive

# Add Vagrant insecure public key
# Taken from Vagrant itself: https://github.com/mitchellh/vagrant/tree/master/keys
RUN echo 'ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6NF8iallvQVp22WDkTkyrtvp9eWW6A8YVr+kz4TjGYe7gHzIw+niNltGEFHzD8+v1I2YJ6oXevct1YeS0o9HZyN1Q9qgCgzUFtdOKLv6IedplqoPkcmF0aYet2PkEDo3MlTBckFXPITAMzF8dJSIFo9D8HfdOV0IAdx4O7PtixWKn5y2hMNG0zQPyUecp4pzC6kivAIhyfHilFR61RGL+GPXQ2MWZWFYbAGjyiYJnAmCP3NOTd0jMZEnDkbUvxhMmBYSdETk1rRgm+R4LOzFUGaHqHDLKLX+FIPKcF96hrucXzcWyLbIbEgE98OHlnVYCzRdK8jlqm8tehUc9c9WhQ== vagrant insecure public key' >> ~/.ssh/authorized_keys

RUN apt-get -qq update
RUN apt-get -qq upgrade -y
RUN apt-get install -y -qq git curl
RUN apt-get install -y -qq bash zsh fish

RUN curl -L http://install.ohmyz.sh | sh || true
RUN chsh -s /usr/bin/zsh
