FROM ubuntu

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN apt-get install -y git

ADD check/check /opt/resource/check
ADD in/in /opt/resource/in
ADD in/git_script.sh /opt/resource/git_script.sh
ADD out/out /opt/resource/out
