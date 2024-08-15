# SPDX-License-Identifier: Apache-2.0

# set a global Docker argument for the default CLI version
#
# https://github.com/moby/moby/issues/37345
ARG TERRAFORM_VERSION=1.7.4

################################################################################
##     docker build --no-cache --target binary -t vela-terraform:binary .     ##
################################################################################

FROM alpine:latest@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5 as binary

ARG TERRAFORM_VERSION

ENV TERRAFORM_ZIP="https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
ENV CHECKSUM_URL="https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_SHA256SUMS"
ENV CHECKSUM_FILE="SHA256SUMS"

# download and verify the Terraform binary
RUN wget -q "${TERRAFORM_ZIP}" -O terraform.zip && \
    wget -q "${CHECKSUM_URL}" -O "${CHECKSUM_FILE}" && \
    EXPECTED_CHECKSUM=$(grep "terraform_${TERRAFORM_VERSION}_linux_amd64.zip" "${CHECKSUM_FILE}" | awk '{ print $1 }') && \
    ACTUAL_CHECKSUM=$(sha256sum terraform.zip | awk '{ print $1 }') && \
    if [ "$EXPECTED_CHECKSUM" != "$ACTUAL_CHECKSUM" ]; then echo "Checksum verification failed"; exit 1; fi && \
    unzip terraform.zip -d /bin && \
    rm -f terraform.zip "${CHECKSUM_FILE}"

##############################################################################
##     docker build --no-cache --target certs -t vela-terraform:certs .     ##
##############################################################################

FROM alpine:latest@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5 as certs

RUN apk add --update --no-cache ca-certificates curl

###############################################################
##     docker build --no-cache -t vela-terraform:local .     ##
###############################################################

FROM alpine:3.20.2@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5

ARG TERRAFORM_VERSION

ENV PLUGIN_TERRAFORM_VERSION=${TERRAFORM_VERSION}

COPY --from=binary /bin/terraform /bin/terraform

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-terraform /bin/terraform-plugin

ENTRYPOINT ["/bin/terraform-plugin"]
