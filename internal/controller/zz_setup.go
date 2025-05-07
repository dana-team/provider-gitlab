// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/pkg/controller"

	file "github.com/dana-team/provider-gitlab/internal/controller/gitlab/file"
	group "github.com/dana-team/provider-gitlab/internal/controller/gitlab/group"
	groupmembership "github.com/dana-team/provider-gitlab/internal/controller/gitlab/groupmembership"
	project "github.com/dana-team/provider-gitlab/internal/controller/gitlab/project"
	projectmembership "github.com/dana-team/provider-gitlab/internal/controller/gitlab/projectmembership"
	user "github.com/dana-team/provider-gitlab/internal/controller/gitlab/user"
	providerconfig "github.com/dana-team/provider-gitlab/internal/controller/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		file.Setup,
		group.Setup,
		groupmembership.Setup,
		project.Setup,
		projectmembership.Setup,
		user.Setup,
		providerconfig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
