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

ENV TERRAFORM_RELEASE_URL="https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}"
ENV TERRAFORM_ZIP_FILENAME="terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
ENV TERRAFORM_CHECKSUMS_FILENAME="terraform_${TERRAFORM_VERSION}_SHA256SUMS"

# download and verify the Terraform binary
RUN wget -q "${TERRAFORM_RELEASE_URL}/${TERRAFORM_ZIP_FILENAME}" -O "${TERRAFORM_ZIP_FILENAME}" && \
    wget -q "${TERRAFORM_RELEASE_URL}/${TERRAFORM_CHECKSUMS_FILENAME}" -O "${TERRAFORM_CHECKSUMS_FILENAME}" && \
    cat "${TERRAFORM_CHECKSUMS_FILENAME}" | grep "${TERRAFORM_ZIP_FILENAME}" | sha256sum -c && \
    unzip "${TERRAFORM_ZIP_FILENAME}" -d /bin && \
    rm -f "${TERRAFORM_ZIP_FILENAME}" "${TERRAFORM_CHECKSUMS_FILENAME}"

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
