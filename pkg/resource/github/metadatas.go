package github

import "github.com/cloudskiff/driftctl/pkg/resource"

func InitMetadatas(resourceSchemaRepository *resource.SchemaRepository) {
	initGithubBranchProtectionMetadata(resourceSchemaRepository)
	initGithubTeamMembershipMetadata(resourceSchemaRepository)
	initGithubMembershipMetadata(resourceSchemaRepository)
	initGithubTeamMetadata(resourceSchemaRepository)
	initGithubRepositoryMetadata(resourceSchemaRepository)
}
