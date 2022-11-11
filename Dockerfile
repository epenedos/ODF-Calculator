FROM golang AS odf
WORKDIR /opt/app-root/src
RUN git clone https://github.com/epenedos/ODF-Calculator.git 
WORKDIR /opt/app-root/src/ODF-Calculator
RUN PWD
RUN go build
RUN ls -lisa
CMD  /opt/app-root/src/ODF-Calculator/main
EXPOSE 8080