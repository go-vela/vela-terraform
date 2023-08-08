# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

# set a global Docker argument for the default CLI version
#
# https://github.com/moby/moby/issues/37345
ARG TERRAFORM_VERSION=1.2.7

################################################################################
##     docker build --no-cache --target binary -t vela-terraform:binary .     ##
################################################################################

FROM alpine:latest@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a as binary

ARG TERRAFORM_VERSION

RUN wget -q https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip -O terraform.zip && \
  unzip terraform.zip -d /bin && \
  rm -f terraform.zip

##############################################################################
##     docker build --no-cache --target certs -t vela-terraform:certs .     ##
##############################################################################

FROM alpine:latest@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a as certs

RUN apk add --update --no-cache ca-certificates

###############################################################
##     docker build --no-cache -t vela-terraform:local .     ##
###############################################################

FROM alpine:3.18.2@sha256:82d1e9d7ed48a7523bdebc18cf6290bdb97b82302a8a9c27d4fe885949ea94d1

ARG TERRAFORM_VERSION

ENV PLUGIN_TERRAFORM_VERSION=${TERRAFORM_VERSION}

COPY --from=binary /bin/terraform /bin/terraform

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-terraform /bin/terraform-plugin

ENTRYPOINT ["/bin/terraform-plugin"]
