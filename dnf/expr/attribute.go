package expr

import "errors"

// AttributeMetadataStorer ...
type AttributeMetadataStorer interface {

	// NameMapping provides a mapping between attribute string name to int
	// e.g., attribute Name: browser -> 1, os -> 2
	// this mapping is used to convert string to int, which can save more memory
	AddNameIDMapping(name string, id uint32)
	// NameValueMapping provides a mapping between attribute string value to int
	// e.g., attribute Name: browser -> 1, then Chrome -> 1, Safari -> 2
	// ["browser"]["Chrome"] = 1
	AddValueIDMapping(name string, value string, id uint32) error
	// Get id of give name
	GetNameID(name string) (uint32, bool)
	// Get id of give name/value
	GetValueID(name string, value string) (uint32, bool)
	// Creates a new Attribute, based on stored mapping
	NewAttribute(name string, values []string, contains bool) (Attribute, error)
}

// NewAttributeMetadataStorer ...
func NewAttributeMetadataStorer() AttributeMetadataStorer {
	return &attributeMetadataStore{
		nameToID:  make(map[string]uint32),
		valueToID: make(map[string]map[string]uint32),
	}
}

type attributeMetadataStore struct {
	nameToID  map[string]uint32
	valueToID map[string]map[string]uint32
}

func (a *attributeMetadataStore) NewAttribute(name string, values []string, contains bool) (Attribute, error) {
	nid, ok := a.GetNameID(name)
	if !ok {
		return Attribute{}, errors.New("name not found")
	}

	vids := make([]uint32, len(values))
	for i, v := range values {
		vid, ok := a.GetValueID(name, v)
		if !ok {
			return Attribute{}, errors.New("value not found")
		}
		vids[i] = vid
	}

	return Attribute{
		Name:     nid,
		Values:   vids,
		Contains: contains,
	}, nil
}

func (a *attributeMetadataStore) AddNameIDMapping(name string, id uint32) {
	a.nameToID[name] = id
}

func (a *attributeMetadataStore) AddValueIDMapping(name string, value string, id uint32) error {
	_, ok := a.nameToID[name]
	if !ok {
		return errors.New("name not found")
	}

	vmap, ok := a.valueToID[name]
	if !ok {
		vmap = make(map[string]uint32)
	}
	vmap[value] = id
	a.valueToID[name] = vmap

	return nil
}

func (a *attributeMetadataStore) GetNameID(name string) (uint32, bool) {
	id, exist := a.nameToID[name]
	return id, exist
}

func (a *attributeMetadataStore) GetValueID(name string, value string) (uint32, bool) {
	_, exist := a.valueToID[name]
	if !exist {
		return 0, false
	}

	id, exist := a.valueToID[name][value]
	if !exist {
		return 0, false
	}
	return id, true
}

// Attribute is the pair of key-value, e.g., age:10, representing 'belongs to'
// The value here is discrete, but for range values, like age < 40, we can convert it into multiple pairs:
// age < 40 ---> age = 10, age = 20, age = 30, the granularity is 10.
//
// Name and Values are internally represented as int value, to save memory space, there should be a mapping
// between the integer and human readable string, like:
// Name: browser -> 1, os -> 2
// Value: Chrome -> 1, Safari -> 2
//
// TODO: for longer range values such as dates, using a hierarchy structure to convert
type Attribute struct {
	Name     uint32
	Values   []uint32
	Contains bool
}
