package privilege

// PermissionType the definition for the result of authority
type PermissionType int

const (
	NONE PermissionType = iota
	OK   PermissionType = 1
	DENY PermissionType = 2
)

// ResourceType the definition for various resources
type ResourceType string

// IResource ...
type IResource interface {
	GetID() string
	GetType() ResourceType
	ParseConf(string, interface{}) error
	CheckPriv(string) PermissionType
}

// IRole ...
type IRole interface {
	Get(interface{}) error
}

// PrivilegeItem the item of privilege details
type PrivilegeItem struct {
	Role               IRole
	Resource           IResource
	PermissionExpected PermissionType
}

type IPrivilegeCenter interface {
	// RegistRole regist various type roles
	RegistRole(IRole) error
	// RegistResource regist various type resources which will be requested by roles
	RegistResource(IResource) error
	// CheckPriv using this function, caller can check whether the permission expected is granted or not
	CheckPriv(IRole, IResource) (PermissionType, error)
}
