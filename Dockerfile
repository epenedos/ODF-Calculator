FROM image-registry.openshift-image-registry.svc:5000/openshift/golang:latest AS odf
WORKDIR /opt/app-root/src
RUN git clone https://github.com/epenedos/ODF-Calculator.git 
WORKDIR /opt/app-root/src/ODF-Calculator
RUN rm ./odf
RUN go build -o odf
RUN ls -lisa
EXPOSE 8080
