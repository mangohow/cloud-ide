FROM ubuntu:20.04

WORKDIR /.workspace

COPY code-server-4.9.0-linux-amd64.tar.gz .
COPY go1.21.3.linux-amd64.tar.gz .
COPY SourceCodeVariable* ./
COPY .init.sh /root/
COPY code-server-bak /root/.local/share/code-server-bak
COPY install.sh .

RUN apt-get -y update                                && \
        apt-get -y install fontconfig                    && \
        mkdir -p ~/.fonts/source-code-pro             && \
        cp SourceCodeVariable* ~/.fonts/source-code-pro  && \
        cd ~/.fonts/source-code-pro                   && \
        fc-cache -f -v                                && \
        cd /.workspace                                && \
        rm -f SourceCodeVariable*                     && \
        tar zxvf code-server-4.9.0-linux-amd64.tar.gz && \
        mv code-server-4.9.0-linux-amd64 code-server  &&\
        rm -f code-server-4.9.0-linux-amd64.tar.gz    && \
        tar zxvf  go1.21.3.linux-amd64.tar.gz -C /usr/local     && \
        rm -f  go1.21.3.linux-amd64.tar.gz             && \
        mkdir -p /go/src /go/pkg /go/bin              && \
        apt-get -qq update                            && \
        apt-get -qq install -y --no-install-recommends ca-certificates curl && \
		apt install -y git make zsh   && \
		chsh -s /bin/zsh              && \
		sh -c ./install.sh

COPY go_tools/* /go/bin/

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct
ENV GOROOT /usr/local/go
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
ENV USER_WORKSPACE /user_data/workspace
ENV CODE_SERVER_DIR /.workspace/code-server
ENV PATH /.workspace/code-server/bin:$PATH

WORKDIR /root

EXPOSE 9999

# 执行脚本文件
CMD ["/bin/bash", ".init.sh"]

