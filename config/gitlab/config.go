package gitlab

import (
	"github.com/crossplane/upjet/pkg/config"
)

const (
	apiVersion = "v1alpha1"
	shortGroup = "gitlab"
)

func Configure(p *config.Provider) {
	p.AddResourceConfigurator("gitlab_group", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "Group"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("gitlab_project", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "Project"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("gitlab_project", func(r *config.Resource) {
		r.References["namespace_id"] = config.Reference{
			TerraformName: "gitlab_group",
			Extractor:     "github.com/crossplane/upjet/pkg/resource.ExtractResourceID()",
		}
	})

	p.AddResourceConfigurator("gitlab_repository_file", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "File"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("gitlab_repository_file", func(r *config.Resource) {
		r.References["project"] = config.Reference{
			TerraformName: "gitlab_project",
			Extractor:     `github.com/crossplane/upjet/pkg/resource.ExtractResourceID()`,
		}
	})

	p.AddResourceConfigurator("gitlab_project_membership", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ProjectMembership"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("gitlab_group_membership", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "GroupMembership"
		r.Version = apiVersion
	})

	p.AddResourceConfigurator("gitlab_user", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "User"
		r.Version = apiVersion
	})
}
