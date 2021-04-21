package github

import "github.com/cloudskiff/driftctl/pkg/resource"

func InitMetadatas(resourceSchemaRepository *resource.SchemaRepository) {
	initGithubBranchProtectionMetadata()
	initGithubTeamMembershipMetadata()
	initGithubMembershipMetadata()
	initGithubTeamMetadata()
	initGithubRepositoryMetadata()
}
