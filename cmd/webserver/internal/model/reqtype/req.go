package reqtype

type SpaceCreateOption struct {
	Name          string `json:"name"`
	TmplId        uint32 `json:"tmpl_id"`
	SpaceSpecId   uint32 `json:"space_spec_id"`
	UserId        uint32 `json:"user_id"`
	GitRepository string `json:"git_repository"`
}
