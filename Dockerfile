FROM golang:latest AS odf
WORKDIR /opt/app-root/src
RUN git clone https://github.com/epenedos/ODF-Calculator.git 
WORKDIR /opt/app-root/src/ODF-Calculator
RUN go build -o odf
RUN ls -lisa
CMD /opt/app-root/src/ODF-Calculator/odf
EXPOSE 8080
