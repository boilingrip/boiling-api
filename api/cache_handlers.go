package api

type FormatsResponse struct {
	Formats []string `json:"formats"`
}

func (a *API) getFormats(ctx *context) {
	ctx.Success(FormatsResponse{Formats: a.c.formats.Keys()})
}

type LeechTypesResponse struct {
	LeechTypes []string `json:"leech_types"`
}

func (a *API) getLeechTypes(ctx *context) {
	ctx.Success(LeechTypesResponse{LeechTypes: a.c.leechTypes.Keys()})
}

type MediaResponse struct {
	Media []string `json:"media"`
}

func (a *API) getMedia(ctx *context) {
	ctx.Success(MediaResponse{Media: a.c.media.Keys()})
}

type ReleaseGroupTypesResponse struct {
	ReleaseGroupTypes []string `json:"release_group_types"`
}

func (a *API) getReleaseGroupTypes(ctx *context) {
	ctx.Success(ReleaseGroupTypesResponse{ReleaseGroupTypes: a.c.releaseGroupTypes.Keys()})
}

type ReleasePropertiesResponse struct {
	ReleaseProperties []string `json:"release_properties"`
}

func (a *API) getReleaseProperties(ctx *context) {
	ctx.Success(ReleasePropertiesResponse{ReleaseProperties: a.c.releaseProperties.Keys()})
}

type ReleaseRolesResponse struct {
	ReleaseRoles []string `json:"release_roles"`
}

func (a *API) getReleaseRoles(ctx *context) {
	ctx.Success(ReleaseRolesResponse{ReleaseRoles: a.c.releaseRoles.Keys()})
}

type PrivilegesResponse struct {
	Privileges []string `json:"privileges"`
}

func (a *API) getPrivileges(ctx *context) {
	ctx.Success(PrivilegesResponse{Privileges: a.c.privileges.Keys()})
}
