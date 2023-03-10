
RELEASE_MANIFEST_FILE ?=
RELEASE_CHGLOG_DESC_FILE ?=

REPOTOOLS_VERSION ?= latest
REPOTOOLS_MODULE = github.com/awslabs/aws-go-multi-module-repository-tools
REPOTOOLS_CMD_CALCULATE_RELEASE = ${REPOTOOLS_MODULE}/cmd/calculaterelease@${REPOTOOLS_VERSION}
REPOTOOLS_CMD_UPDATE_REQUIRES = ${REPOTOOLS_MODULE}/cmd/updaterequires@${REPOTOOLS_VERSION}
REPOTOOLS_CMD_UPDATE_MODULE_METADATA = ${REPOTOOLS_MODULE}/cmd/updatemodulemeta@${REPOTOOLS_VERSION}
REPOTOOLS_CMD_GENERATE_CHANGELOG = ${REPOTOOLS_MODULE}/cmd/generatechangelog@${REPOTOOLS_VERSION}
REPOTOOLS_CMD_CHANGELOG = ${REPOTOOLS_MODULE}/cmd/changelog@${REPOTOOLS_VERSION}
REPOTOOLS_CMD_TAG_RELEASE = ${REPOTOOLS_MODULE}/cmd/tagrelease@${REPOTOOLS_VERSION}

smithy-publish-local:
	cd codegen && ./gradlew publishToMavenLocal

smithy-build:
	cd codegen && ./gradlew build

smithy-clean:
	cd codegen && ./gradlew clean

#####################
#  Release Process  #
#####################
.PHONY: preview-release pre-release-validation release

preview-release:
	go run ${REPOTOOLS_CMD_CALCULATE_RELEASE}

pre-release-validation:
	@if [[ -z "${RELEASE_MANIFEST_FILE}" ]]; then \
		echo "RELEASE_MANIFEST_FILE is required to specify the file to write the release manifest" && false; \
	fi
	@if [[ -z "${RELEASE_CHGLOG_DESC_FILE}" ]]; then \
		echo "RELEASE_CHGLOG_DESC_FILE is required to specify the file to write the release notes" && false; \
	fi

release: pre-release-validation
	go run ${REPOTOOLS_CMD_CALCULATE_RELEASE} -o ${RELEASE_MANIFEST_FILE}
	go run ${REPOTOOLS_CMD_UPDATE_REQUIRES} -release ${RELEASE_MANIFEST_FILE}
	go run ${REPOTOOLS_CMD_UPDATE_MODULE_METADATA} -release ${RELEASE_MANIFEST_FILE}
	go run ${REPOTOOLS_CMD_GENERATE_CHANGELOG} -release ${RELEASE_MANIFEST_FILE} -o ${RELEASE_CHGLOG_DESC_FILE}
	go run ${REPOTOOLS_CMD_CHANGELOG} rm -all
	go run ${REPOTOOLS_CMD_TAG_RELEASE} -release ${RELEASE_MANIFEST_FILE}

##############
# Repo Tools #
##############
.PHONY: install-changelog

install-changelog:
	go install ${REPOTOOLS_MODULE}/cmd/changelog@${REPOTOOLS_VERSION}
