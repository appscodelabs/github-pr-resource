FROM ubuntu

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN apt-get install -y git

ADD check/check /opt/resource/check
ADD in/in /opt/resource/in
ADD in/git_script.sh /git_script.sh
ADD out/out /opt/resource/out
ADD out/find_hash.sh /find_hash.sh
ADD out/fetch_pr.sh /fetch_pr.sh
